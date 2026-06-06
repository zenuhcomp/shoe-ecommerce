package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kibsoft/luna-collections/backend/internal/auth"
	"github.com/kibsoft/luna-collections/backend/internal/config"
	"github.com/kibsoft/luna-collections/backend/internal/models"
	"github.com/kibsoft/luna-collections/backend/internal/repository"
)

type AdminAPI struct {
	repo *repository.ContentRepo
	cfg  config.Config
}

func NewAdmin(repo *repository.ContentRepo, cfg config.Config) *AdminAPI {
	return &AdminAPI{repo: repo, cfg: cfg}
}

func (a *AdminAPI) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	email := strings.TrimSpace(strings.ToLower(req.Email))
	row, err := a.repo.UserByEmail(r.Context(), email)
	if err != nil || !auth.CheckPassword(row.PasswordHash, req.Password) {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	token, err := auth.IssueToken(a.cfg.JWTSecret, row.ID, row.Email, row.Role,
		time.Duration(a.cfg.JWTExpiryHours)*time.Hour)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "token error")
		return
	}
	writeJSON(w, http.StatusOK, models.LoginResponse{
		Token: token,
		User: models.AdminUser{
			ID: row.ID, Email: row.Email, DisplayName: row.DisplayName,
			Role: row.Role, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		},
	})
}

func (a *AdminAPI) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	user, err := a.repo.UserByID(r.Context(), claims.UserID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (a *AdminAPI) Dashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	services, _ := a.repo.Services(ctx)
	reviews, _ := a.repo.Reviews(ctx, false)
	faq, _ := a.repo.FAQ(ctx)
	gallery, _ := a.repo.Gallery(ctx)
	bookings, _ := a.repo.ListBookings(ctx, 5)
	writeJSON(w, http.StatusOK, map[string]any{
		"counts": map[string]int{
			"services": len(services),
			"reviews":  len(reviews),
			"faq":      len(faq),
			"gallery":  len(gallery),
		},
		"recentBookings": bookings,
	})
}

func (a *AdminAPI) GetSite(w http.ResponseWriter, r *http.Request) {
	site, err := a.repo.Site(r.Context())
	if err != nil {
		log.Printf("admin site: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to load site")
		return
	}
	writeJSON(w, http.StatusOK, site)
}

func (a *AdminAPI) PutSite(w http.ResponseWriter, r *http.Request) {
	var site models.SiteConfig
	if err := json.NewDecoder(r.Body).Decode(&site); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := a.repo.UpdateSite(r.Context(), site); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save site")
		return
	}
	writeJSON(w, http.StatusOK, site)
}

func (a *AdminAPI) ListServices(w http.ResponseWriter, r *http.Request) {
	list, err := a.repo.Services(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (a *AdminAPI) CreateService(w http.ResponseWriter, r *http.Request) {
	var s models.LashService
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if s.ID == "" {
		s.ID = "svc-" + uuid.NewString()[:8]
	}
	if s.Slug == "" {
		s.Slug = slugify(s.Name)
	}
	order, _ := a.repo.NextServiceSort(r.Context())
	if err := a.repo.CreateService(r.Context(), s, order); err != nil {
		writeError(w, http.StatusConflict, "could not create service")
		return
	}
	writeJSON(w, http.StatusCreated, s)
}

func (a *AdminAPI) UpdateService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var s models.LashService
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	s.ID = id
	if err := a.repo.UpdateService(r.Context(), s); err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, s)
}

func (a *AdminAPI) DeleteService(w http.ResponseWriter, r *http.Request) {
	a.deleteID(w, r, a.repo.DeleteService)
}

func (a *AdminAPI) ReorderServices(w http.ResponseWriter, r *http.Request) {
	a.reorderIDs(w, r, a.repo.ReorderServices)
}

func (a *AdminAPI) ListFeatures(w http.ResponseWriter, r *http.Request) {
	list, err := a.repo.Features(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (a *AdminAPI) CreateFeature(w http.ResponseWriter, r *http.Request) {
	var f models.Feature
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if f.ID == "" {
		f.ID = "feat-" + uuid.NewString()[:8]
	}
	order, _ := a.repo.MaxSort(r.Context(), "features")
	if err := a.repo.CreateFeature(r.Context(), f, order); err != nil {
		writeError(w, http.StatusConflict, "could not create")
		return
	}
	writeJSON(w, http.StatusCreated, f)
}

func (a *AdminAPI) UpdateFeature(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var f models.Feature
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	f.ID = id
	if err := a.repo.UpdateFeature(r.Context(), f); err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, f)
}

func (a *AdminAPI) DeleteFeature(w http.ResponseWriter, r *http.Request) {
	a.deleteID(w, r, a.repo.DeleteFeature)
}

func (a *AdminAPI) ReorderFeatures(w http.ResponseWriter, r *http.Request) {
	a.reorderIDs(w, r, a.repo.ReorderFeatures)
}

func (a *AdminAPI) ListPress(w http.ResponseWriter, r *http.Request) {
	list, err := a.repo.Press(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (a *AdminAPI) CreatePress(w http.ResponseWriter, r *http.Request) {
	var p models.PressMention
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if p.ID == "" {
		p.ID = "press-" + uuid.NewString()[:8]
	}
	order, _ := a.repo.MaxSort(r.Context(), "press_mentions")
	if err := a.repo.CreatePress(r.Context(), p, order); err != nil {
		writeError(w, http.StatusConflict, "could not create")
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (a *AdminAPI) UpdatePress(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var p models.PressMention
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	p.ID = id
	if err := a.repo.UpdatePress(r.Context(), p); err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (a *AdminAPI) DeletePress(w http.ResponseWriter, r *http.Request) {
	a.deleteID(w, r, a.repo.DeletePress)
}

func (a *AdminAPI) ReorderPress(w http.ResponseWriter, r *http.Request) {
	a.reorderIDs(w, r, a.repo.ReorderPress)
}

func (a *AdminAPI) ListReviews(w http.ResponseWriter, r *http.Request) {
	list, err := a.repo.Reviews(r.Context(), false)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (a *AdminAPI) CreateReview(w http.ResponseWriter, r *http.Request) {
	var rev models.Review
	if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if rev.ID == "" {
		rev.ID = "rev-" + uuid.NewString()[:8]
	}
	if rev.Initials == "" && rev.Name != "" {
		rev.Initials = initials(rev.Name)
	}
	order, _ := a.repo.MaxSort(r.Context(), "reviews")
	pub := rev.Published
	if err := a.repo.CreateReview(r.Context(), rev, order); err != nil {
		writeError(w, http.StatusConflict, "could not create")
		return
	}
	if !pub {
		_ = a.repo.UpdateReview(r.Context(), rev, false)
		rev.Published = false
	} else {
		rev.Published = true
	}
	writeJSON(w, http.StatusCreated, rev)
}

func (a *AdminAPI) UpdateReview(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var rev models.Review
	if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	rev.ID = id
	if err := a.repo.UpdateReview(r.Context(), rev, rev.Published); err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, rev)
}

func (a *AdminAPI) DeleteReview(w http.ResponseWriter, r *http.Request) {
	a.deleteID(w, r, a.repo.DeleteReview)
}

func (a *AdminAPI) ReorderReviews(w http.ResponseWriter, r *http.Request) {
	a.reorderIDs(w, r, a.repo.ReorderReviews)
}

func (a *AdminAPI) ListFAQ(w http.ResponseWriter, r *http.Request) {
	list, err := a.repo.FAQ(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (a *AdminAPI) CreateFAQ(w http.ResponseWriter, r *http.Request) {
	var f models.FAQ
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if f.ID == 0 {
		f.ID, _ = a.repo.NextFAQID(r.Context())
	}
	order, _ := a.repo.MaxSort(r.Context(), "faq_items")
	if err := a.repo.CreateFAQ(r.Context(), f, order); err != nil {
		writeError(w, http.StatusConflict, "could not create")
		return
	}
	writeJSON(w, http.StatusCreated, f)
}

func (a *AdminAPI) UpdateFAQ(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var f models.FAQ
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	f.ID = id
	if err := a.repo.UpdateFAQ(r.Context(), f); err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, f)
}

func (a *AdminAPI) DeleteFAQ(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := a.repo.DeleteFAQ(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *AdminAPI) ReorderFAQ(w http.ResponseWriter, r *http.Request) {
	var req models.ReorderIntRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := a.repo.ReorderFAQ(r.Context(), req.IDs); err != nil {
		writeError(w, http.StatusInternalServerError, "reorder failed")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *AdminAPI) ListGallery(w http.ResponseWriter, r *http.Request) {
	list, err := a.repo.Gallery(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (a *AdminAPI) CreateGallery(w http.ResponseWriter, r *http.Request) {
	var g models.GalleryItem
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if g.ID == "" {
		g.ID = "g-" + uuid.NewString()[:8]
	}
	order, _ := a.repo.MaxSort(r.Context(), "gallery_items")
	if err := a.repo.CreateGallery(r.Context(), g, order); err != nil {
		writeError(w, http.StatusConflict, "could not create")
		return
	}
	writeJSON(w, http.StatusCreated, g)
}

func (a *AdminAPI) UpdateGallery(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var g models.GalleryItem
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	g.ID = id
	if err := a.repo.UpdateGallery(r.Context(), g); err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, g)
}

func (a *AdminAPI) DeleteGallery(w http.ResponseWriter, r *http.Request) {
	a.deleteID(w, r, a.repo.DeleteGallery)
}

func (a *AdminAPI) ReorderGallery(w http.ResponseWriter, r *http.Request) {
	a.reorderIDs(w, r, a.repo.ReorderGallery)
}

func (a *AdminAPI) ListNav(w http.ResponseWriter, r *http.Request) {
	list, err := a.repo.Nav(r.Context(), false)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (a *AdminAPI) CreateNav(w http.ResponseWriter, r *http.Request) {
	var n models.NavLink
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if n.ID == "" {
		n.ID = "nav-" + uuid.NewString()[:8]
	}
	order, _ := a.repo.MaxSort(r.Context(), "nav_links")
	if err := a.repo.CreateNav(r.Context(), n, order); err != nil {
		writeError(w, http.StatusConflict, "could not create")
		return
	}
	writeJSON(w, http.StatusCreated, n)
}

func (a *AdminAPI) UpdateNav(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var n models.NavLink
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	n.ID = id
	if err := a.repo.UpdateNav(r.Context(), n); err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, n)
}

func (a *AdminAPI) DeleteNav(w http.ResponseWriter, r *http.Request) {
	a.deleteID(w, r, a.repo.DeleteNav)
}

func (a *AdminAPI) ReorderNav(w http.ResponseWriter, r *http.Request) {
	a.reorderIDs(w, r, a.repo.ReorderNav)
}

func (a *AdminAPI) ListBookings(w http.ResponseWriter, r *http.Request) {
	list, err := a.repo.ListBookings(r.Context(), 200)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load bookings")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (a *AdminAPI) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	a.deleteID(w, r, a.repo.DeleteBooking)
}

func (a *AdminAPI) deleteID(w http.ResponseWriter, r *http.Request, fn func(context.Context, string) error) {
	id := chi.URLParam(r, "id")
	if err := fn(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *AdminAPI) reorderIDs(w http.ResponseWriter, r *http.Request, fn func(context.Context, []string) error) {
	var req models.ReorderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := fn(r.Context(), req.IDs); err != nil {
		writeError(w, http.StatusInternalServerError, "reorder failed")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	return strings.ReplaceAll(s, " ", "-")
}

func initials(name string) string {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "?"
	}
	if len(parts) == 1 {
		r := []rune(parts[0])
		if len(r) >= 2 {
			return strings.ToUpper(string(r[:2]))
		}
		return strings.ToUpper(string(r))
	}
	a := []rune(parts[0])
	b := []rune(parts[len(parts)-1])
	return strings.ToUpper(string(a[0])) + strings.ToUpper(string(b[0]))
}
