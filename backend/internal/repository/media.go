package repository

import (
	"context"
	"fmt"

	"github.com/kibsoft/luna-collections/backend/internal/models"
)

const mediaSelectCols = `id, name, category, media_type, url, thumbnail_url, alt_text, tags`

func (r *ContentRepo) MediaAssets(ctx context.Context, category string) ([]models.MediaAsset, error) {
	q := fmt.Sprintf(`SELECT %s FROM media_assets`, mediaSelectCols)
	var args []any
	if category != "" && category != "all" {
		q += ` WHERE category = ?`
		args = append(args, category)
	}
	q += ` ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, rebind(r.driver, q), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.MediaAsset
	for rows.Next() {
		var m models.MediaAsset
		if err := rows.Scan(&m.ID, &m.Name, &m.Category, &m.MediaType, &m.URL, &m.ThumbnailURL, &m.AltText, &m.Tags); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

func (r *ContentRepo) CreateMediaAsset(ctx context.Context, m models.MediaAsset) error {
	ts := nowUTC()
	q := rebind(r.driver, `
INSERT INTO media_assets (id, name, category, media_type, url, thumbnail_url, alt_text, tags, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	_, err := r.db.ExecContext(ctx, q,
		m.ID, m.Name, m.Category, m.MediaType, m.URL, m.ThumbnailURL, m.AltText, m.Tags, ts, ts,
	)
	return err
}

func (r *ContentRepo) DeleteMediaAsset(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, rebind(r.driver, `DELETE FROM media_assets WHERE id = ?`), id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("media asset not found")
	}
	return nil
}

func (r *ContentRepo) DeleteAllMedia(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM media_assets`)
	return err
}

func (r *ContentRepo) SeedDefaultMediaIfEmpty(ctx context.Context) error {
	var count int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM media_assets`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	samples := []models.MediaAsset{
		// Furniture
		{ID: "med-furn-1", Name: "Velvet Sculptural Sofa", Category: "furniture", MediaType: "image", URL: "https://images.unsplash.com/photo-1555041469-a586c61ea9bc?w=1200&q=80", AltText: "Emerald velvet sofa in modern Scandinavian living room", Tags: "furniture, sofa, living room, modern"},
		{ID: "med-furn-2", Name: "Nordic Dining Table Set", Category: "furniture", MediaType: "image", URL: "https://images.unsplash.com/photo-1617806118233-18e1de247200?w=1200&q=80", AltText: "Minimalist oak dining table with curved chairs", Tags: "furniture, dining, oak, table"},
		{ID: "med-furn-3", Name: "Architectural Lounge Armchair", Category: "furniture", MediaType: "image", URL: "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=1200&q=80", AltText: "Designer armchair with textured boucle cushion", Tags: "furniture, armchair, lounge"},
		{ID: "med-furn-4", Name: "Warm Living Space Showcase", Category: "furniture", MediaType: "video", URL: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4", AltText: "Interior living space walk-through video", Tags: "furniture, video, interior, showcase"},

		// Motor Vehicles
		{ID: "med-veh-1", Name: "Matte Black Performance Coupe", Category: "vehicles", MediaType: "image", URL: "https://images.unsplash.com/photo-1503376780353-7e6692767b70?w=1200&q=80", AltText: "Black luxury sports car front view", Tags: "vehicles, sports, coupe, luxury"},
		{ID: "med-veh-2", Name: "Urban Premium Electric SUV", Category: "vehicles", MediaType: "image", URL: "https://images.unsplash.com/photo-1533473359331-0135ef1b58bf?w=1200&q=80", AltText: "Modern silver SUV on mountain road", Tags: "vehicles, suv, modern, luxury"},
		{ID: "med-veh-3", Name: "Classic Executive Sedan", Category: "vehicles", MediaType: "image", URL: "https://images.unsplash.com/photo-1552519507-da3b142c6e3d?w=1200&q=80", AltText: "High-end luxury sedan in showroom lighting", Tags: "vehicles, sedan, executive"},
		{ID: "med-veh-4", Name: "Supercar Cinematic Walkaround", Category: "vehicles", MediaType: "video", URL: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4", AltText: "Automotive high-speed track and walkaround video", Tags: "vehicles, video, track, automotive"},

		// Spa & Wellness
		{ID: "med-spa-1", Name: "Hot Stone Aromatherapy Ritual", Category: "spa", MediaType: "image", URL: "https://images.unsplash.com/photo-1544161515-4ab6ce6db874?w=1200&q=80", AltText: "Relaxing massage with volcanic hot stones", Tags: "spa, massage, wellness, stone"},
		{ID: "med-spa-2", Name: "Serene Bamboo Bath Sanctuary", Category: "spa", MediaType: "image", URL: "https://images.unsplash.com/photo-1600334089648-b0d9d3028eb2?w=1200&q=80", AltText: "Japanese inspired wooden cedar bath tub with petals", Tags: "spa, bath, relaxation, sanctuary"},
		{ID: "med-spa-3", Name: "Holistic Herbal Facial Treatment", Category: "spa", MediaType: "image", URL: "https://images.unsplash.com/photo-1512290900672-1f02e6a0058b?w=1200&q=80", AltText: "Facial treatment with organic botanicals", Tags: "spa, facial, skincare, beauty"},
		{ID: "med-spa-4", Name: "Cascading Waterfall & Zen Ambience", Category: "spa", MediaType: "video", URL: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerJoyBlazes.mp4", AltText: "Soothing natural water and zen wellness backdrop video", Tags: "spa, video, water, zen, relaxation"},

		// Beauty & Cosmetics
		{ID: "med-beau-1", Name: "Rose Petal Glow Serum Drops", Category: "beauty", MediaType: "image", URL: "https://images.unsplash.com/photo-1620916566398-39f1143ab7be?w=1200&q=80", AltText: "Glass dropper bottle with botanical hydrating serum", Tags: "beauty, skincare, serum, cosmetics"},
		{ID: "med-beau-2", Name: "Velvet Lip Color Palette", Category: "beauty", MediaType: "image", URL: "https://images.unsplash.com/photo-1586495777744-4413f21062fa?w=1200&q=80", AltText: "Luxury lip colors with silky textures", Tags: "beauty, makeup, lipstick, glamour"},
		{ID: "med-beau-3", Name: "Hydra-Luminous Complexion Cream", Category: "beauty", MediaType: "image", URL: "https://images.unsplash.com/photo-1598440947619-2c35fc9aa908?w=1200&q=80", AltText: "Organic moisturizing face cream with golden lid", Tags: "beauty, moisturizer, luxury, organic"},
		{ID: "med-beau-4", Name: "Radiant Skin Application Demo", Category: "beauty", MediaType: "video", URL: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/WeAreGoingOnBullrun.mp4", AltText: "Beauty editorial application and glow showcase video", Tags: "beauty, video, tutorial, glow"},

		// Footwear / Fashion
		{ID: "med-foot-1", Name: "Italian Pointed Court Pump", Category: "footwear", MediaType: "image", URL: "https://images.unsplash.com/photo-1543163521-1bf539c55dd2?w=1200&q=80", AltText: "Classic patent navy pump with elegant stiletto", Tags: "footwear, pump, official, shoes"},
		{ID: "med-foot-2", Name: "Polished Oxford Brogue", Category: "footwear", MediaType: "image", URL: "https://images.unsplash.com/photo-1761726065663-6e550c25a6ef?w=1200&q=80", AltText: "Handcrafted calfskin leather oxford shoe", Tags: "footwear, oxford, leather, office"},
		{ID: "med-foot-3", Name: "Structured Block Heel Mule", Category: "footwear", MediaType: "image", URL: "https://images.unsplash.com/photo-1596702874230-b5706dfb5bc7?w=1200&q=80", AltText: "Minimalist block heel comfortable office shoe", Tags: "footwear, block heel, comfortable"},
	}

	for _, s := range samples {
		_ = r.CreateMediaAsset(ctx, s)
	}
	return nil
}
