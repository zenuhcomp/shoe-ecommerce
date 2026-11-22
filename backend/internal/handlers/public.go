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

func (a *PublicAPI) Sections(w http.ResponseWriter, r *http.Request) {
	list, err := a.repo.Sections(r.Context(), false)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load sections")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (a *PublicAPI) Presets(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.repo.ListPresets())
}

func (a *PublicAPI) Media(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	list, err := a.repo.MediaAssets(r.Context(), category)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load media")
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
		Message: "Thank you! We will confirm your request via WhatsApp shortly.",
	})
}

func (a *PublicAPI) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req models.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		req.Name = strings.TrimSpace(req.CustomerName)
	}
	if req.Phone == "" {
		req.Phone = strings.TrimSpace(req.CustomerPhone)
	}
	if req.Email == "" {
		req.Email = strings.TrimSpace(req.CustomerEmail)
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)
	req.ServiceSlug = strings.TrimSpace(req.ServiceSlug)
	req.PaymentMethod = strings.TrimSpace(req.PaymentMethod)
	if req.Name == "" || req.Phone == "" || req.ServiceSlug == "" {
		writeError(w, http.StatusBadRequest, "name, phone, and item are required")
		return
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	svc, ok, err := a.repo.ServiceBySlug(r.Context(), req.ServiceSlug)
	if err != nil || !ok {
		writeError(w, http.StatusBadRequest, "selected item was not found")
		return
	}

	site, err := a.repo.Site(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load site settings")
		return
	}

	amount := svc.PriceKES * req.Quantity
	status := "pending_offline"
	message := "Your inquiry has been received. Our team will contact you shortly."
	whatsAppURL := ""

	switch req.PaymentMethod {
	case "whatsapp":
		status = "pending_whatsapp"
		message = "Order inquiry noted. Opening WhatsApp to finalize sizing & delivery."
		whatsAppURL = appendWhatsAppMessage(site.WhatsAppURL,
			buildWhatsAppOrderMessage(site.BrandName, svc, req, amount))
	case "delivery":
		status = "pending_delivery"
		message = "Order recorded for pay on delivery. We will call you to confirm dispatch."
	case "mpesa_express":
		status = "pending_payment"
		message = "Prompt sent. Please complete payment on your phone."
	case "offline":
		status = "pending_offline"
		message = "Order received. Please complete payment using the provided till or paybill."
	default:
		req.PaymentMethod = "whatsapp"
		status = "pending_whatsapp"
		message = "Order recorded. Opening WhatsApp to connect."
		whatsAppURL = appendWhatsAppMessage(site.WhatsAppURL,
			buildWhatsAppOrderMessage(site.BrandName, svc, req, amount))
	}

	orderRecord := models.OrderRecord{
		Name:             req.Name,
		Phone:            req.Phone,
		Email:            req.Email,
		ServiceSlug:      svc.Slug,
		ServiceName:      svc.Name,
		Size:             req.Size,
		Quantity:         req.Quantity,
		AmountKES:        amount,
		PaymentMethod:    req.PaymentMethod,
		MpesaPhone:       req.MpesaPhone,
		OfflineChannel:   req.OfflineChannel,
		PaymentReference: req.PaymentReference,
		Status:           status,
		Notes:            req.Notes,
	}

	id, err := a.repo.AddOrder(r.Context(), orderRecord)
	if err != nil {
		log.Printf("create order: %v", err)
		writeError(w, http.StatusInternalServerError, "could not save order")
		return
	}

	writeJSON(w, http.StatusCreated, models.OrderResponse{
		ID: id, Status: status, Message: message, WhatsAppURL: whatsAppURL,
	})
}

func buildWhatsAppOrderMessage(brand string, svc models.LashService, req models.OrderRequest, amount int) string {
	lines := []string{
		fmt.Sprintf("Hi %s, I'd like to place an inquiry/order:", brand),
		fmt.Sprintf("• Item: %s", svc.Name),
		fmt.Sprintf("• Option/Spec: %s", fallback(req.Size, "Standard")),
		fmt.Sprintf("• Qty: %d", req.Quantity),
		fmt.Sprintf("• Amount: %s %d", svc.Currency, amount),
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
