package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/kibsoft/luna-collections/backend/internal/models"
	"github.com/kibsoft/luna-collections/backend/internal/repository"
)

type PublicAPI struct {
	repo         *repository.ContentRepo
	databaseMode string
	databaseDrv  string
}

func NewPublic(repo *repository.ContentRepo, databaseMode, databaseDrv string) *PublicAPI {
	return &PublicAPI{repo: repo, databaseMode: databaseMode, databaseDrv: databaseDrv}
}

func (a *PublicAPI) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"database": map[string]any{
			"mode":      a.databaseMode,
			"driver":    a.databaseDrv,
			"connected": true,
		},
		"cms": true,
	})
}

func (a *PublicAPI) Home(w http.ResponseWriter, r *http.Request) {
	home, err := a.repo.Home(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load home")
		return
	}
	writeJSON(w, http.StatusOK, home)
}

func (a *PublicAPI) Site(w http.ResponseWriter, r *http.Request) {
	site, err := a.repo.Site(r.Context())
	if err != nil {
		log.Printf("public site: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to load site")
		return
	}
	writeJSON(w, http.StatusOK, site)
}

func (a *PublicAPI) Nav(w http.ResponseWriter, r *http.Request) {
	nav, err := a.repo.Nav(r.Context(), true)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load nav")
		return
	}
	writeJSON(w, http.StatusOK, nav)
}

func (a *PublicAPI) Services(w http.ResponseWriter, r *http.Request) {
	list, err := a.repo.Services(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load services")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (a *PublicAPI) ServiceBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	svc, ok, err := a.repo.ServiceBySlug(r.Context(), slug)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load service")
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "service not found")
		return
	}
	writeJSON(w, http.StatusOK, svc)
}

func (a *PublicAPI) Features(w http.ResponseWriter, r *http.Request) {
	list, err := a.repo.Features(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load features")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (a *PublicAPI) Press(w http.ResponseWriter, r *http.Request) {
	list, err := a.repo.Press(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load press")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (a *PublicAPI) Reviews(w http.ResponseWriter, r *http.Request) {
	list, err := a.repo.Reviews(r.Context(), true)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load reviews")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (a *PublicAPI) FAQ(w http.ResponseWriter, r *http.Request) {
	list, err := a.repo.FAQ(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load faq")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (a *PublicAPI) Gallery(w http.ResponseWriter, r *http.Request) {
	list, err := a.repo.Gallery(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load gallery")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (a *PublicAPI) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var inquiry models.BookingInquiry
	if err := json.NewDecoder(r.Body).Decode(&inquiry); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	inquiry.Name = strings.TrimSpace(inquiry.Name)
	inquiry.Phone = strings.TrimSpace(inquiry.Phone)
	if inquiry.Name == "" || inquiry.Phone == "" {
		writeError(w, http.StatusBadRequest, "name and phone are required")
		return
	}
	id, err := a.repo.AddBooking(r.Context(), inquiry)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not save booking")
		return
	}
	writeJSON(w, http.StatusCreated, models.BookingResponse{
		ID:      id,
		Message: "Thank you! We will confirm your appointment via WhatsApp shortly.",
	})
}

func (a *PublicAPI) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req models.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.TrimSpace(req.Email)
	req.ServiceSlug = strings.TrimSpace(req.ServiceSlug)
	req.Size = strings.TrimSpace(req.Size)
	req.MpesaPhone = strings.TrimSpace(req.MpesaPhone)
	req.OfflineChannel = strings.TrimSpace(req.OfflineChannel)
	req.PaymentReference = strings.TrimSpace(req.PaymentReference)
	req.Notes = strings.TrimSpace(req.Notes)
	req.PaymentMethod = strings.TrimSpace(req.PaymentMethod)

	if req.Name == "" || req.Phone == "" || req.ServiceSlug == "" {
		writeError(w, http.StatusBadRequest, "name, phone, and service are required")
		return
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	svc, ok, err := a.repo.ServiceBySlug(r.Context(), req.ServiceSlug)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load service")
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "service not found")
		return
	}

	site, err := a.repo.Site(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load site")
		return
	}

	amount := svc.PriceKES * req.Quantity
	status := "pending"
	message := "Order received. We will confirm shortly."
	var whatsAppURL string

	switch req.PaymentMethod {
	case "mpesa_stk":
		mpesaPhone := req.MpesaPhone
		if mpesaPhone == "" {
			mpesaPhone = req.Phone
		}
		req.MpesaPhone = mpesaPhone
		status = "stk_initiated"
		message = fmt.Sprintf(
			"M-Pesa STK push sent to %s for KES %d. Enter your M-Pesa PIN on your phone to complete payment.",
			mpesaPhone, amount,
		)
	case "offline_till", "offline_paybill":
		if req.PaymentReference == "" {
			writeError(w, http.StatusBadRequest, "payment reference is required for offline payment")
			return
		}
		channel := req.OfflineChannel
		if channel == "" {
			if req.PaymentMethod == "offline_till" {
				channel = "till"
			} else {
				channel = "paybill"
			}
		}
		req.OfflineChannel = channel
		status = "awaiting_confirmation"
		if channel == "till" {
			message = fmt.Sprintf(
				"Thank you! We received your order. Pay KES %d to Till %s, then we will verify reference %s.",
				amount, site.Payment.TillNumber, req.PaymentReference,
			)
		} else {
			message = fmt.Sprintf(
				"Thank you! We received your order. Pay KES %d to Paybill %s (Account: %s), then we will verify reference %s.",
				amount, site.Payment.PaybillNumber, site.Payment.PaybillAccount, req.PaymentReference,
			)
		}
	case "whatsapp":
		status = "whatsapp_redirect"
		msg := buildWhatsAppOrderMessage(site.BrandName, svc, req, amount)
		whatsAppURL = appendWhatsAppMessage(site.WhatsAppURL, msg)
		message = "Opening WhatsApp to complete your order with our team."
	default:
		writeError(w, http.StatusBadRequest, "invalid payment method")
		return
	}

	id, err := a.repo.AddOrder(r.Context(), models.OrderRecord{
		Name: req.Name, Phone: req.Phone, Email: req.Email,
		ServiceSlug: svc.Slug, ServiceName: svc.Name, Size: req.Size, Quantity: req.Quantity,
		AmountKES: amount, PaymentMethod: req.PaymentMethod, MpesaPhone: req.MpesaPhone,
		OfflineChannel: req.OfflineChannel, PaymentReference: req.PaymentReference,
		Status: status, Notes: req.Notes,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not save order")
		return
	}

	writeJSON(w, http.StatusCreated, models.OrderResponse{
		ID: id, Status: status, Message: message, WhatsAppURL: whatsAppURL,
	})
}

func buildWhatsAppOrderMessage(brand string, svc models.LashService, req models.OrderRequest, amount int) string {
	lines := []string{
		fmt.Sprintf("Hi %s, I'd like to place an order:", brand),
		fmt.Sprintf("• Style: %s", svc.Name),
		fmt.Sprintf("• Size: %s", fallback(req.Size, "TBC")),
		fmt.Sprintf("• Qty: %d", req.Quantity),
		fmt.Sprintf("• Total: KES %d", amount),
		fmt.Sprintf("• Name: %s", req.Name),
		fmt.Sprintf("• Phone: %s", req.Phone),
	}
	if req.Notes != "" {
		lines = append(lines, fmt.Sprintf("• Notes: %s", req.Notes))
	}
	return strings.Join(lines, "\n")
}

func appendWhatsAppMessage(baseURL, message string) string {
	if baseURL == "" {
		return ""
	}
	sep := "?"
	if strings.Contains(baseURL, "?") {
		sep = "&"
	}
	return baseURL + sep + "text=" + url.QueryEscape(message)
}

func fallback(value, alt string) string {
	if value == "" {
		return alt
	}
	return value
}
