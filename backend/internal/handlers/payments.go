package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/kibsoft/luna-collections/backend/internal/models"
	"github.com/kibsoft/luna-collections/backend/internal/payments"
	"github.com/kibsoft/luna-collections/backend/internal/repository"
)

type PaymentsAPI struct {
	repo    *repository.ContentRepo
	manager *payments.PaymentManager
}

func NewPaymentsAPI(repo *repository.ContentRepo, manager *payments.PaymentManager) *PaymentsAPI {
	return &PaymentsAPI{
		repo:    repo,
		manager: manager,
	}
}

// GetActiveMethods returns all payment methods that are active/enabled by the admin.
func (p *PaymentsAPI) GetActiveMethods(w http.ResponseWriter, r *http.Request) {
	site, err := p.repo.Site(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load payment settings")
		return
	}
	methods := p.manager.ListActiveMethods(site.Payment)
	writeJSON(w, http.StatusOK, map[string]any{
		"methods": methods,
	})
}

// Initiate executes the strategy for the chosen payment method.
func (p *PaymentsAPI) Initiate(w http.ResponseWriter, r *http.Request) {
	var req models.PaymentInitiateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Method == "" {
		req.Method = "mpesa_express"
	}
	if req.AmountKES <= 0 {
		writeError(w, http.StatusBadRequest, "amount must be greater than zero")
		return
	}

	site, err := p.repo.Site(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load payment settings")
		return
	}

	res, err := p.manager.Process(r.Context(), req, site.Payment)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Link checkout request ID / transaction ID to order if an orderId was supplied
	if req.OrderID != "" && res.TransactionID != "" {
		_ = p.repo.UpdateOrderPaymentRef(r.Context(), req.OrderID, res.TransactionID)
	}

	writeJSON(w, http.StatusOK, res)
}

// Verify verifies a manual transaction reference or checks STK status.
func (p *PaymentsAPI) Verify(w http.ResponseWriter, r *http.Request) {
	var req models.PaymentVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Reference == "" {
		writeError(w, http.StatusBadRequest, "reference code is required")
		return
	}

	site, err := p.repo.Site(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load payment settings")
		return
	}

	res, err := p.manager.Verify(r.Context(), req.Method, req.Reference, site.Payment)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// If verified with a code, update the order status
	if req.OrderID != "" && res.Success {
		_ = p.repo.UpdateOrderStatus(r.Context(), req.OrderID, "paid")
		_ = p.repo.UpdateOrderPaymentRef(r.Context(), req.OrderID, req.Reference)
	}

	writeJSON(w, http.StatusOK, res)
}

// GetOrderStatus checks the live status of an order (used by checkout polling).
func (p *PaymentsAPI) GetOrderStatus(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("orderId")
	if orderID == "" {
		writeError(w, http.StatusBadRequest, "orderId is required")
		return
	}

	order, err := p.repo.GetOrderByID(r.Context(), orderID)
	if err != nil {
		writeError(w, http.StatusNotFound, "order not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":               order.ID,
		"status":           order.Status,
		"paymentMethod":    order.PaymentMethod,
		"paymentReference": order.PaymentReference,
		"amountKes":        order.AmountKES,
		"notes":            order.Notes,
	})
}

// DarajaCallbackPayload represents the authentic Safaricom Daraja STK Push callback JSON schema.
type DarajaCallbackPayload struct {
	Body struct {
		StkCallback struct {
			MerchantRequestID string `json:"MerchantRequestID"`
			CheckoutRequestID string `json:"CheckoutRequestID"`
			ResultCode        int    `json:"ResultCode"`
			ResultDesc        string `json:"ResultDesc"`
			CallbackMetadata  *struct {
				Item []struct {
					Name  string `json:"Name"`
					Value any    `json:"Value"`
				} `json:"Item"`
			} `json:"CallbackMetadata"`
		} `json:"stkCallback"`
	} `json:"Body"`
}

// MpesaCallback receives Safaricom Daraja STK Push callbacks and reconciles orders.
func (p *PaymentsAPI) MpesaCallback(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "cannot read callback body")
		return
	}
	log.Printf("[Daraja Callback] Received payload: %s", string(body))

	var payload DarajaCallbackPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("[Daraja Callback] JSON decode error: %v", err)
		writeJSON(w, http.StatusOK, map[string]any{"ResultCode": 0, "ResultDesc": "Accepted with warning"})
		return
	}

	cb := payload.Body.StkCallback
	checkoutID := cb.CheckoutRequestID
	resultCode := cb.ResultCode
	resultDesc := cb.ResultDesc

	var receipt string
	var amount float64
	var phone string

	if cb.CallbackMetadata != nil {
		for _, item := range cb.CallbackMetadata.Item {
			switch item.Name {
			case "MpesaReceiptNumber":
				receipt = fmt.Sprintf("%v", item.Value)
			case "Amount":
				if v, ok := item.Value.(float64); ok {
					amount = v
				}
			case "PhoneNumber":
				switch pv := item.Value.(type) {
				case float64:
					phone = fmt.Sprintf("%.0f", pv)
				default:
					phone = fmt.Sprintf("%v", item.Value)
				}
			}
		}
	}

	// Reconcile order
	if checkoutID != "" {
		if resultCode == 0 {
			if receipt == "" {
				receipt = checkoutID
			}
			note := fmt.Sprintf("[M-Pesa Callback Verified] Receipt: %s, Amount: KES %.0f, Phone: %s", receipt, amount, phone)
			if o, err := p.repo.ReconcileOrderPayment(r.Context(), checkoutID, receipt, "paid", note); err == nil {
				log.Printf("[Daraja Reconciliation] Order %s marked PAID with receipt %s", o.ID, receipt)
			} else {
				log.Printf("[Daraja Reconciliation] Note: %v", err)
			}
		} else {
			note := fmt.Sprintf("[M-Pesa Callback Failed] ResultCode: %d, Desc: %s", resultCode, resultDesc)
			if o, err := p.repo.ReconcileOrderPayment(r.Context(), checkoutID, "", "payment_failed", note); err == nil {
				log.Printf("[Daraja Reconciliation] Order %s marked PAYMENT_FAILED (%s)", o.ID, resultDesc)
			}
		}
	}

	// Standard Safaricom acknowledge payload
	writeJSON(w, http.StatusOK, map[string]any{
		"ResultCode": 0,
		"ResultDesc": "Callback accepted successfully",
	})
}

// SimulateCallback allows testing end-to-end reconciliation by creating a simulated Safaricom callback.
func (p *PaymentsAPI) SimulateCallback(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrderID           string  `json:"orderId"`
		CheckoutRequestID string  `json:"checkoutRequestId"`
		ResultCode        int     `json:"resultCode"`
		ResultDesc        string  `json:"resultDesc"`
		ReceiptNumber     string  `json:"receiptNumber"`
		Amount            float64 `json:"amount"`
		Phone             string  `json:"phone"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.OrderID == "" && req.CheckoutRequestID == "" {
		writeError(w, http.StatusBadRequest, "orderId or checkoutRequestId is required")
		return
	}

	lookupKey := req.CheckoutRequestID
	if lookupKey == "" {
		lookupKey = req.OrderID
	}

	if req.ReceiptNumber == "" {
		req.ReceiptNumber = fmt.Sprintf("QA%dKL%d", time.Now().Unix()%100000, 7)
	}
	if req.ResultDesc == "" {
		req.ResultDesc = "The service request is processed successfully."
	}
	if req.Phone == "" {
		req.Phone = "254722000000"
	}

	status := "paid"
	note := fmt.Sprintf("[M-Pesa Simulated Callback Verified] Receipt: %s, Amount: KES %.0f, Phone: %s", req.ReceiptNumber, req.Amount, req.Phone)
	if req.ResultCode != 0 {
		status = "payment_failed"
		note = fmt.Sprintf("[M-Pesa Simulated Callback Failed] Code: %d, Desc: %s", req.ResultCode, req.ResultDesc)
	}

	updated, err := p.repo.ReconcileOrderPayment(r.Context(), lookupKey, req.ReceiptNumber, status, note)
	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("reconciliation failed: %v", err))
		return
	}

	log.Printf("[Daraja Simulation] Successfully reconciled order %s to status '%s' (Receipt: %s)", updated.ID, status, req.ReceiptNumber)

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Order payment reconciled successfully via simulated callback",
		"order":   updated,
	})
}

// --- Admin Payment Handlers ---

// GetAdminPayments returns full payment configuration including keys for admin.
func (p *PaymentsAPI) GetAdminPayments(w http.ResponseWriter, r *http.Request) {
	site, err := p.repo.Site(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load payment settings")
		return
	}
	methods := p.manager.ListActiveMethods(site.Payment)
	writeJSON(w, http.StatusOK, map[string]any{
		"config":  site.Payment,
		"methods": methods,
	})
}

// UpdateAdminPayments updates payment configuration and enabled/disabled toggles.
func (p *PaymentsAPI) UpdateAdminPayments(w http.ResponseWriter, r *http.Request) {
	var newCfg models.PaymentConfig
	if err := json.NewDecoder(r.Body).Decode(&newCfg); err != nil {
		writeError(w, http.StatusBadRequest, "invalid payment settings body")
		return
	}

	site, err := p.repo.Site(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load site")
		return
	}

	// Preserve existing credentials and keys if the update left them blank
	if newCfg.ConsumerKey == "" {
		newCfg.ConsumerKey = site.Payment.ConsumerKey
	}
	if newCfg.ConsumerSecret == "" {
		newCfg.ConsumerSecret = site.Payment.ConsumerSecret
	}
	if newCfg.PassKey == "" {
		newCfg.PassKey = site.Payment.PassKey
	}
	if newCfg.ShortCode == "" {
		newCfg.ShortCode = site.Payment.ShortCode
	}
	if newCfg.Environment == "" {
		newCfg.Environment = site.Payment.Environment
	}
	if newCfg.CallbackURL == "" {
		newCfg.CallbackURL = site.Payment.CallbackURL
	}
	if newCfg.CardSecretKey == "" {
		newCfg.CardSecretKey = site.Payment.CardSecretKey
	}

	site.Payment = newCfg
	if err := p.repo.UpdateSite(r.Context(), site); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save payment settings")
		return
	}

	methods := p.manager.ListActiveMethods(site.Payment)
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Payment settings updated successfully",
		"config":  site.Payment,
		"methods": methods,
	})
}
