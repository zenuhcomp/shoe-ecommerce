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

	// Always sync admin credentials from configuration / .env so login never fails
	if cfg.AdminEmail != "" && cfg.AdminPassword != "" {
		hash, err := auth.HashPassword(cfg.AdminPassword)
		if err != nil {
			return err
		}
		if _, err := repo.UpsertAdminUser(ctx, cfg.AdminEmail, hash, cfg.AdminName); err != nil {
			log.Printf("sync admin user (%s): %v", cfg.AdminEmail, err)
		} else {
			log.Printf("CMS admin synced from env: %s", cfg.AdminEmail)
		}

		// Also ensure admin@lunacollections.local exists with same password if different
		if cfg.AdminEmail != "admin@lunacollections.local" {
			_, _ = repo.UpsertAdminUser(ctx, "admin@lunacollections.local", hash, "Default Admin")
		}
		// Also ensure aniinyoo@gmail.com exists with same password so either email works
		if cfg.AdminEmail != "aniinyoo@gmail.com" {
			_, _ = repo.UpsertAdminUser(ctx, "aniinyoo@gmail.com", hash, "Site Admin")
		}
	}

	// Always sync Safaricom production credentials from .env and eliminate hardcoded Till/Paybill
	site, err := repo.Site(ctx)
	if err == nil {
		changed := false
		if cfg.ConsumerKey != "" && site.Payment.ConsumerKey != cfg.ConsumerKey {
			site.Payment.ConsumerKey = cfg.ConsumerKey
			changed = true
		}
		if cfg.ConsumerSecret != "" && site.Payment.ConsumerSecret != cfg.ConsumerSecret {
			site.Payment.ConsumerSecret = cfg.ConsumerSecret
			changed = true
		}
		if cfg.PassKey != "" && site.Payment.PassKey != cfg.PassKey {
			site.Payment.PassKey = cfg.PassKey
			changed = true
		}
		if cfg.ShortCode != "" && site.Payment.ShortCode != cfg.ShortCode {
			site.Payment.ShortCode = cfg.ShortCode
			changed = true
		}
		if cfg.CallbackURL != "" && site.Payment.CallbackURL != cfg.CallbackURL {
			site.Payment.CallbackURL = cfg.CallbackURL
			changed = true
		}
		if site.Payment.Environment != "production" {
			site.Payment.Environment = "production"
			changed = true
		}
		if cfg.ConsumerKey != "" && cfg.PassKey != "" && !site.Payment.MpesaExpressEnabled {
			site.Payment.MpesaExpressEnabled = true
			changed = true
		}

		// Wipe out legacy hardcoded Till numbers & Paybill if they match demo seeds
		demoTills := map[string]bool{"902411": true, "808500": true, "888420": true, "710200": true, "612300": true}
		if demoTills[site.Payment.TillNumber] {
			site.Payment.TillNumber = ""
			site.Payment.MpesaShortCode = ""
			site.Payment.MpesaManualEnabled = false
			changed = true
		}
		if site.Payment.PaybillNumber == "247247" {
			site.Payment.PaybillNumber = ""
			site.Payment.PaybillAccount = ""
			site.Payment.MpesaManualEnabled = false
			changed = true
		}
		if site.Payment.ShortCode == "174379" && cfg.ShortCode != "" {
			site.Payment.ShortCode = cfg.ShortCode
			changed = true
		}

		if changed {
			_ = repo.UpdateSite(ctx, site)
			log.Printf("Safaricom production credentials synced (ShortCode=%s, Env=%s)", site.Payment.ShortCode, site.Payment.Environment)
		}
	}

	return nil
}
