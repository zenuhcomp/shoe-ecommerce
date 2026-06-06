package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/kibsoft/luna-collections/backend/internal/auth"
	"github.com/kibsoft/luna-collections/backend/internal/bootstrap"
	"github.com/kibsoft/luna-collections/backend/internal/config"
	"github.com/kibsoft/luna-collections/backend/internal/db"
	"github.com/kibsoft/luna-collections/backend/internal/handlers"
	"github.com/kibsoft/luna-collections/backend/internal/repository"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()
	database, err := db.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer database.Close()

	repo := repository.NewContentRepo(database.DB, database.Driver)
	if err := bootstrap.Seed(ctx, repo, cfg); err != nil {
		log.Fatalf("bootstrap: %v", err)
	}

	public := handlers.NewPublic(repo, cfg.DatabaseMode, database.Driver)
	admin := handlers.NewAdmin(repo, cfg)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/api/health", public.Health)
	r.Route("/api", func(r chi.Router) {
		r.Get("/home", public.Home)
		r.Get("/site", public.Site)
		r.Get("/nav", public.Nav)
		r.Get("/services", public.Services)
		r.Get("/services/{slug}", public.ServiceBySlug)
		r.Get("/features", public.Features)
		r.Get("/press", public.Press)
		r.Get("/reviews", public.Reviews)
		r.Get("/faq", public.FAQ)
		r.Get("/gallery", public.Gallery)
		r.Post("/bookings", public.CreateBooking)
		r.Post("/orders", public.CreateOrder)
	})

	r.Route("/api/admin", func(r chi.Router) {
		r.Post("/auth/login", admin.Login)
		r.Group(func(r chi.Router) {
			r.Use(auth.Middleware(cfg.JWTSecret))
			r.Get("/auth/me", admin.Me)
			r.Get("/dashboard", admin.Dashboard)

			r.Get("/site", admin.GetSite)
			r.Put("/site", admin.PutSite)

			r.Get("/services", admin.ListServices)
			r.Post("/services", admin.CreateService)
			r.Put("/services/reorder", admin.ReorderServices)
			r.Put("/services/{id}", admin.UpdateService)
			r.Delete("/services/{id}", admin.DeleteService)

			r.Get("/features", admin.ListFeatures)
			r.Post("/features", admin.CreateFeature)
			r.Put("/features/reorder", admin.ReorderFeatures)
			r.Put("/features/{id}", admin.UpdateFeature)
			r.Delete("/features/{id}", admin.DeleteFeature)

			r.Get("/press", admin.ListPress)
			r.Post("/press", admin.CreatePress)
			r.Put("/press/reorder", admin.ReorderPress)
			r.Put("/press/{id}", admin.UpdatePress)
			r.Delete("/press/{id}", admin.DeletePress)

			r.Get("/reviews", admin.ListReviews)
			r.Post("/reviews", admin.CreateReview)
			r.Put("/reviews/reorder", admin.ReorderReviews)
			r.Put("/reviews/{id}", admin.UpdateReview)
			r.Delete("/reviews/{id}", admin.DeleteReview)

			r.Get("/faq", admin.ListFAQ)
			r.Post("/faq", admin.CreateFAQ)
			r.Put("/faq/reorder", admin.ReorderFAQ)
			r.Put("/faq/{id}", admin.UpdateFAQ)
			r.Delete("/faq/{id}", admin.DeleteFAQ)

			r.Get("/gallery", admin.ListGallery)
			r.Post("/gallery", admin.CreateGallery)
			r.Put("/gallery/reorder", admin.ReorderGallery)
			r.Put("/gallery/{id}", admin.UpdateGallery)
			r.Delete("/gallery/{id}", admin.DeleteGallery)

			r.Get("/nav", admin.ListNav)
			r.Post("/nav", admin.CreateNav)
			r.Put("/nav/reorder", admin.ReorderNav)
			r.Put("/nav/{id}", admin.UpdateNav)
			r.Delete("/nav/{id}", admin.DeleteNav)

			r.Get("/bookings", admin.ListBookings)
			r.Delete("/bookings/{id}", admin.DeleteBooking)
		})
	})

	if staticDir := frontendDistDir(); staticDir != "" {
		log.Printf("Serving frontend from %s", staticDir)
		r.NotFound(spaHandler(staticDir))
	}

	addr := ":" + cfg.BackendPort
	srv := &http.Server{Addr: addr, Handler: r}

	go func() {
		log.Printf("Luna Collections API on %s (DATABASE_MODE=%s, CMS enabled)", addr, cfg.DatabaseMode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

func frontendDistDir() string {
	candidates := []string{}
	if p := os.Getenv("FRONTEND_DIST_DIR"); p != "" {
		candidates = append(candidates, p)
	}
	candidates = append(candidates,
		"../frontend/dist/frontend/browser",
		"frontend/dist/frontend/browser",
	)
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(exeDir, "../frontend/dist/frontend/browser"),
			filepath.Join(exeDir, "../../frontend/dist/frontend/browser"),
		)
	}

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		abs, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}
		if info, err := os.Stat(filepath.Join(abs, "index.html")); err == nil && !info.IsDir() {
			return abs
		}
	}
	return ""
}

func spaHandler(staticDir string) http.HandlerFunc {
	fileServer := http.FileServer(http.Dir(staticDir))
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/api" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.NotFound(w, r)
			return
		}

		path := filepath.Clean(strings.TrimPrefix(r.URL.Path, "/"))
		if path == "." {
			path = "index.html"
		}
		fullPath := filepath.Join(staticDir, path)
		if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	}
}
