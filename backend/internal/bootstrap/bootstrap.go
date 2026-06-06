package bootstrap

import (
	"context"
	"log"

	"github.com/kibsoft/luna-collections/backend/internal/auth"
	"github.com/kibsoft/luna-collections/backend/internal/config"
	"github.com/kibsoft/luna-collections/backend/internal/repository"
)

func Seed(ctx context.Context, repo *repository.ContentRepo, cfg config.Config) error {
	if err := repo.SeedIfEmpty(ctx); err != nil {
		return err
	}
	if err := repo.BackfillServiceMedia(ctx); err != nil {
		return err
	}
	if err := repo.BackfillGalleryMedia(ctx); err != nil {
		return err
	}
	if err := repo.BackfillSitePayment(ctx); err != nil {
		return err
	}
	n, err := repo.CountUsers(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	hash, err := auth.HashPassword(cfg.AdminPassword)
	if err != nil {
		return err
	}
	user, err := repo.CreateUser(ctx, cfg.AdminEmail, hash, cfg.AdminName)
	if err != nil {
		return err
	}
	log.Printf("CMS admin created: %s (change CMS_ADMIN_PASSWORD in .env)", user.Email)
	return nil
}
