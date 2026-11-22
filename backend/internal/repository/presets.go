package repository

import (
	"context"
	"fmt"
	"github.com/kibsoft/luna-collections/backend/internal/models"
)

func (r *ContentRepo) ListPresets() []models.PresetSummary {
	return []models.PresetSummary{
		{
			ID:           "furniture",
			Name:         "Vanguard Atelier",
			Industry:     "furniture",
			Tagline:      "Sculptural Living & Modern Furniture",
			Description:  "Tailored for furniture studios, interior design showrooms, and home living brands. Includes room dimensions, wood species, 3D room scenes, and custom sizing requests.",
			Icon:         "🛋️",
			PrimaryColor: "#181D19",
			AccentColor:  "#D4A373",
			SampleImages: []string{
				"https://images.unsplash.com/photo-1555041469-a586c61ea9bc?w=600&q=80",
				"https://images.unsplash.com/photo-1617806118233-18e1de247200?w=600&q=80",
				"https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=600&q=80",
			},
		},
		{
			ID:           "vehicles",
			Name:         "Apex Motors Prestige",
			Industry:     "vehicles",
			Tagline:      "Luxury & High-Performance Motor Vehicles",
			Description:  "Built for car dealerships, automotive brokers, and exotic rentals. Includes mileage, engine specs, transmission, test-drive booking, and walkaround videos.",
			Icon:         "🏎️",
			PrimaryColor: "#0D1117",
			AccentColor:  "#38BDF8",
			SampleImages: []string{
				"https://images.unsplash.com/photo-1503376780353-7e6692767b70?w=600&q=80",
				"https://images.unsplash.com/photo-1533473359331-0135ef1b58bf?w=600&q=80",
				"https://images.unsplash.com/photo-1552519507-da3b142c6e3d?w=600&q=80",
			},
		},
		{
			ID:           "spa",
			Name:         "Serena Wellness",
			Industry:     "spa",
			Tagline:      "Holistic Sanctuary & Thermal Spa Retreat",
			Description:  "Crafted for day spas, wellness clinics, retreats, and aesthetic practitioners. Includes treatment durations, therapist details, serene video backdrops, and reservation engine.",
			Icon:         "🌿",
			PrimaryColor: "#2D3A34",
			AccentColor:  "#C2A676",
			SampleImages: []string{
				"https://images.unsplash.com/photo-1544161515-4ab6ce6db874?w=600&q=80",
				"https://images.unsplash.com/photo-1600334089648-b0d9d3028eb2?w=600&q=80",
				"https://images.unsplash.com/photo-1512290900672-1f02e6a0058b?w=600&q=80",
			},
		},
		{
			ID:           "beauty",
			Name:         "Aura Glow Lab",
			Industry:     "beauty",
			Tagline:      "Clinical Botanical Skincare & Clean Beauty",
			Description:  "Created for beauty brands, cosmetics lines, and skincare labs. Includes skin type recommendations, active ingredients, dewy video backgrounds, and instant ordering.",
			Icon:         "✨",
			PrimaryColor: "#2C1B2D",
			AccentColor:  "#E08F84",
			SampleImages: []string{
				"https://images.unsplash.com/photo-1620916566398-39f1143ab7be?w=600&q=80",
				"https://images.unsplash.com/photo-1586495777744-4413f21062fa?w=600&q=80",
				"https://images.unsplash.com/photo-1598440947619-2c35fc9aa908?w=600&q=80",
			},
		},
		{
			ID:           "footwear",
			Name:         "Luna Collections",
			Industry:     "footwear",
			Tagline:      "Curated Ladies Official & Office Footwear",
			Description:  "Designed for premium shoe boutiques and official accessories. Includes heel heights, 360-degree rotation view, leather care guides, and sizing consultations.",
			Icon:         "👠",
			PrimaryColor: "#1A1F3B",
			AccentColor:  "#C4A574",
			SampleImages: []string{
				"https://images.unsplash.com/photo-1543163521-1bf539c55dd2?w=600&q=80",
				"https://images.unsplash.com/photo-1761726065663-6e550c25a6ef?w=600&q=80",
				"https://images.unsplash.com/photo-1596702874230-b5706dfb5bc7?w=600&q=80",
			},
		},
	}
}

func (r *ContentRepo) ApplyPreset(ctx context.Context, presetID string, seedItems bool) error {
	// Preserve existing payment credentials if previously configured
	var existingPayment *models.PaymentConfig
	if existing, err := r.Site(ctx); err == nil {
		p := existing.Payment
		existingPayment = &p
	}

	var err error
	switch presetID {
	case "furniture":
		err = r.applyFurniturePreset(ctx, seedItems)
	case "vehicles":
		err = r.applyVehiclesPreset(ctx, seedItems)
	case "spa":
		err = r.applySpaPreset(ctx, seedItems)
	case "beauty":
		err = r.applyBeautyPreset(ctx, seedItems)
	case "footwear":
		err = r.applyFootwearPreset(ctx, seedItems)
	default:
		return fmt.Errorf("unknown preset: %s", presetID)
	}
	if err != nil {
		return err
	}

	// Restore administrator's configured payment settings completely
	if existingPayment != nil {
		if current, err := r.Site(ctx); err == nil {
			current.Payment = *existingPayment
			_ = r.UpdateSite(ctx, current)
		}
	}
	return nil
}

// --- Furniture Preset ---
func (r *ContentRepo) applyFurniturePreset(ctx context.Context, seedItems bool) error {
	site := models.SiteConfig{
		Industry:           "furniture",
		BrandName:          "Vanguard",
		BrandScript:        "Atelier",
		Tagline:            "Sculptural Living · Handcrafted Modernity",
		HeroTitle:          "Architectural Furniture for Spaces with Character",
		HeroSubtitle:       "Bespoke curved sofas, solid white oak dining tables, and sculptural lighting handcrafted to last generations.",
		HeroBadge:          "◆ Handcrafted Solid Wood & Designer Living — Nairobi ◆",
		HeroVideoURL:       "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4",
		HeroBackgroundType: "video",
		HeroOverlayOpacity: 45,
		Announcement:       "✦ Bespoke Craftsmanship · 5-Year Structural Warranty · White Glove Delivery in Nairobi ✦",
		Location:           "Westlands Showroom, Nairobi",
		StudioAddress:      "Delta Corner Annex, 4th Floor, Westlands, Nairobi",
		Phone:              "+254 700 882 119",
		Email:              "concierge@vanguardliving.co.ke",
		WhatsAppURL:        "https://wa.me/254700882119",
		InstagramURL:       "https://instagram.com/vanguardliving.ke",
		TikTokURL:          "https://tiktok.com/@vanguardliving",
		PriceRange:         "KES 25,000 – 180,000",
		Copyright:          "© 2026 Vanguard Atelier. Handcrafted in Kenya.",
		MetaDescription:    "Vanguard Atelier — Handcrafted architectural furniture in Nairobi. Bespoke sofas, oak tables, and accent chairs with white-glove delivery.",
		Theme: models.Theme{
			Primary:      "#181D19",
			PrimaryDark:  "#0F1310",
			Accent:       "#D4A373",
			AccentLight:  "#FAEDCD",
			Background:   "#F9F6F0",
			Surface:      "#FFFFFF",
			Text:         "#1F2421",
			TextMuted:    "#6B705C",
			HeroGradient: "linear-gradient(135deg, #0F1310 0%, #181D19 55%, #28302B 100%)",
			FontHeading:  "'Playfair Display', Georgia, serif",
			FontBody:     "'DM Sans', system-ui, sans-serif",
			FontScript:   "'Cormorant Garamond', Georgia, serif",
		},
		Animations: models.AnimationConfig{
			GlobalEnabled:   true,
			Intensity:       "balanced",
			ScrollReveal:    true,
			DefaultDuration: "0.8s",
			CardHoverEffect: "lift",
		},
		About: models.AboutSection{
			Title: "The Vanguard Philosophy",
			Intro: "We believe true luxury lies in timeless silhouette, honest timber, and tactile perfection.",
			Story: "Born out of an obsession with mid-century lines and African hardwood heritage, Vanguard Atelier collaborates with master joiners and textile weavers to create heirloom furniture. Each piece is constructed from ethically sourced timber, hand-finished with organic oils, and tailored to your living space.",
			Highlights: []string{
				"Ethically sourced solid African teak, mahogany, and European white oak",
				"Custom dimensional tailoring and fabric swatches for every home",
				"5-year structural warranty on all internal hardwood framing",
				"Dedicated interior staging consultations and white-glove installation",
			},
			ArtistName:  "Vanguard Craft Team",
			ArtistTitle: "Master Artisans & Interior Designers",
		},
		HouseCalls: models.HouseCallsSection{
			Title:       "White-Glove Delivery & Installation",
			Description: "Our specialized logistics team unpacks, positions, assembles, and removes all packaging from your home or office space.",
			Notes: []string{
				"Complimentary white-glove delivery within Nairobi for orders over KES 75,000",
				"Fully insured regional dispatch across East Africa",
				"Pre-delivery room dimension verification by our design team",
				"Eco-friendly packaging recycling program",
			},
		},
		Aftercare: models.AftercareSection{
			Title: "Timber Care & Longevity Guide",
			Intro: "Natural hardwood matures gracefully. Follow our simple routine to keep your furniture lustrous for decades.",
			Tips: []models.AftercareTip{
				{Title: "Maintain moisture balance", Description: "Keep solid timber pieces away from direct heat sources and harsh continuous sunlight. Apply natural beeswax conditioner twice a year."},
				{Title: "Bouclé & velvet maintenance", Description: "Vacuum textured upholstery using a soft brush attachment weekly. Blot liquid spills immediately with clean linen — never scrub."},
				{Title: "Natural marble care", Description: "Use coasters for acidic beverages. Clean honed stone tops with warm water and pH-neutral stone soap."},
			},
		},
		PricingHighlight: models.PricingHighlight{
			Title:      "Transparent Atelier Pricing",
			Subtitle:   "Artisanal quality with zero middlemen markup — from workshop straight to your sanctuary.",
			PriceFrom:  25000,
			PriceTo:    180000,
			Currency:   "KES",
			Badge:      "Atelier Direct",
			Includes:   "Craft certificate, fabric care kit, white-glove delivery, and 5-year frame warranty.",
			StudioNote: "◆ Westlands Design Studio · ◆ Bespoke Finishes Available · ◆ 3-4 Week Craft Lead Time",
		},
		Payment: models.PaymentConfig{
			MpesaManualEnabled: false,
			TillNumber:         "",
			PaybillNumber:      "",
			PaybillAccount:     "",
			MpesaShortCode:     "",
			CashEnabled:        true,
			CashInstructions:   "Cash accepted upon showroom collection or white-glove delivery.",
		},
	}
	if err := r.UpdateSite(ctx, site); err != nil {
		return err
	}

	// Update Sections
	sections := []models.Section{
		{
			ID:                "hero",
			Type:              "hero",
			Enabled:           true,
			SortOrder:         0,
			Title:             "Architectural Furniture for Spaces with Character",
			Subtitle:          "Bespoke curved sofas, solid white oak dining tables, and sculptural lighting handcrafted to last generations.",
			Badge:             "◆ Handcrafted Solid Wood & Designer Living ◆",
			BackgroundType:    "video",
			BackgroundVideo:   "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4",
			OverlayOpacity:    45,
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
			PrimaryCtaText:    "Explore Collections",
			PrimaryCtaUrl:     "/collections",
			SecondaryCtaText:  "Consult on WhatsApp",
			SecondaryCtaUrl:   "https://wa.me/254700882119",
		},
		{
			ID:                "features",
			Type:              "features",
			Enabled:           true,
			SortOrder:         1,
			Title:             "Built for <span class='accent'>Enduring Beauty</span>",
			Subtitle:          "Every joint, fabric, and curve is crafted with deliberate intention",
			Badge:             "~ The Vanguard Standard ~",
			BackgroundType:    "default",
			Animation:         "fade-in",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
			HoverEffect:       "lift",
		},
		{
			ID:                "products",
			Type:              "products",
			Enabled:           true,
			SortOrder:         2,
			Title:             "Curated <span class='accent'>Living Pieces</span>",
			Subtitle:          "Handmade to order — explore full dimensions, wood finishes, and 360° details",
			Badge:             "~ Signature Collection ~",
			BackgroundType:    "default",
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
			HoverEffect:       "lift",
			PrimaryCtaText:    "View All Furniture",
			PrimaryCtaUrl:     "/collections",
		},
		{
			ID:                "video_banner",
			Type:              "video_banner",
			Enabled:           true,
			SortOrder:         3,
			Title:             "The Art of <span class='accent'>Joinery & Craft</span>",
			Subtitle:          "Witness how seasoned timber transforms into timeless interior sculptures in our workshop.",
			Badge:             "~ Inside the Workshop ~",
			BackgroundType:    "video",
			BackgroundVideo:   "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4",
			OverlayOpacity:    50,
			Animation:         "zoom-in",
			AnimationDuration: "0.9s",
			AnimationDelay:    "0.1s",
			PrimaryCtaText:    "Watch Full Craft Film",
			PrimaryCtaUrl:     "#craft-film",
		},
		{
			ID:                "gallery",
			Type:              "gallery",
			Enabled:           true,
			SortOrder:         4,
			Title:             "Spaces Crafted <span class='accent'>For Living</span>",
			Subtitle:          "A glimpse inside private residences and boutique spaces furnished by Vanguard",
			Badge:             "~ Interior Portfolio ~",
			BackgroundType:    "default",
			Animation:         "fade-in",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
			PrimaryCtaText:    "Explore Lookbook",
			PrimaryCtaUrl:     "/gallery",
		},
		{
			ID:                "pricing",
			Type:              "pricing",
			Enabled:           true,
			SortOrder:         5,
			Title:             "Transparent Pricing & Bespoke Commissioning",
			Subtitle:          "Honest craftsmanship with complete price transparency.",
			Badge:             "~ Value & Quality ~",
			BackgroundType:    "color",
			BackgroundColor:   "#181D19",
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
			PrimaryCtaText:    "Request Custom Quote",
			PrimaryCtaUrl:     "/pricing",
		},
		{
			ID:                "reviews",
			Type:              "reviews",
			Enabled:           true,
			SortOrder:         6,
			Title:             "Loved by <span class='accent'>Designers & Homeowners</span>",
			Subtitle:          "Read what architects and residents have to say about their Vanguard furniture",
			Badge:             "~ Verified Testimonials ~",
			BackgroundType:    "default",
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
		},
		{
			ID:                "faq",
			Type:              "faq",
			Enabled:           true,
			SortOrder:         7,
			Title:             "Frequently <span class='accent'>Asked Questions</span>",
			Subtitle:          "Custom dimensions, lead times, timber care, and white glove installation",
			Badge:             "~ Everything You Need to Know ~",
			BackgroundType:    "default",
			Animation:         "fade-in",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
		},
		{
			ID:                "inquiry",
			Type:              "inquiry",
			Enabled:           true,
			SortOrder:         8,
			Title:             "Order or <span class='accent'>Commission a Piece</span>",
			Subtitle:          "Send us your room dimensions or schedule a complimentary private studio tour",
			Badge:             "~ Start Your Project ~",
			BackgroundType:    "default",
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
		},
	}
	if err := r.DeleteAllSections(ctx); err != nil {
		return err
	}
	for i, s := range sections {
		if err := r.CreateSection(ctx, s, i); err != nil {
			return err
		}
	}

	if seedItems {
		items := []models.LashService{
			{
				ID: "furn-oslo-sofa", Name: "Oslo Sculptural 3-Seater Sofa", Slug: "oslo-sculptural-3-seater-sofa",
				Emoji: "🛋️", Style: "Living Room", Length: "3-Seater", Category: "Sofas & Seating", Badge: "Best Seller",
				Description: "Organic curved silhouette upholstered in Italian bouclé with an internal kiln-dried teak framework. Designed for deep comfort and visual poetry.",
				PriceKES:    95000, OriginalPriceKES: 110000, Currency: "KES", Featured: true, InStock: true,
				Specs:    map[string]string{"Material": "Premium Italian Bouclé & High-Density Memory Foam", "Dimensions": "W: 240cm × D: 95cm × H: 78cm", "Frame": "FSC-Certified Kiln-Dried Teak", "Assembly": "White-Glove Delivery Included"},
				Images:   []string{"https://images.unsplash.com/photo-1555041469-a586c61ea9bc?w=800&q=80", "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=800&q=80"},
				View360:  []string{"https://images.unsplash.com/photo-1555041469-a586c61ea9bc?w=800&q=80"},
				VideoURL: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4",
			},
			{
				ID: "furn-kyoto-table", Name: "Kyoto Solid White Oak Dining Table", Slug: "kyoto-solid-white-oak-dining-table",
				Emoji: "🪵", Style: "Dining", Length: "8-10 Seater", Category: "Dining Tables", Badge: "Artisan Made",
				Description: "Minimalist Japanese-Scandinavian joinery showcasing continuous end-grain flow with a live-edge contour and zero synthetic veneers.",
				PriceKES:    120000, OriginalPriceKES: 135000, Currency: "KES", Featured: true, InStock: true,
				Specs:   map[string]string{"Timber": "American White Oak (Solid 45mm Slab)", "Finish": "Organic Matte Hardwax Oil (Food-Safe)", "Dimensions": "L: 260cm × W: 100cm × H: 76cm", "Seats": "8 – 10 Guests Comfortably"},
				Images:  []string{"https://images.unsplash.com/photo-1615066390971-03e4e1c36ddf?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1615066390971-03e4e1c36ddf?w=800&q=80"},
			},
			{
				ID: "furn-aura-armchair", Name: "Aura Walnut Lounge Chair & Ottoman", Slug: "aura-walnut-lounge-chair-ottoman",
				Emoji: "🪑", Style: "Accent Chair", Length: "Standard", Category: "Accent Chairs", Badge: "Architect's Choice",
				Description: "Bent American walnut veneer shells lined with hand-buffed saddle leather. Precision balanced 360-degree swivel mechanism on die-cast aluminium base.",
				PriceKES:    48000, OriginalPriceKES: 55000, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Wood": "7-Ply Molded American Walnut", "Leather": "Full-Grain Saddle Tan Aniline", "Base": "Die-Cast Powder-Coated Aluminium Swivel", "Includes": "Matching Ergonomic Footstool"},
				Images:  []string{"https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=800&q=80"},
			},
			{
				ID: "furn-carrara-table", Name: "Nordic Fluted Carrara Marble Coffee Table", Slug: "nordic-fluted-carrara-marble-coffee-table",
				Emoji: "🏛️", Style: "Living Room", Length: "Round 90cm", Category: "Coffee & Side Tables", Badge: "Natural Stone",
				Description: "Honed Italian Carrara marble disc anchored on a fluted natural white ash plinth. Sealed against wine and coffee stains.",
				PriceKES:    38000, OriginalPriceKES: 45000, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Stone": "Authentic Italian Carrara Marble Top (20mm)", "Base": "Solid Ash Wood Fluted Column", "Diameter": "90cm Diameter × 40cm Height", "Treatment": "Deep-Penetrating Hydrophobic Sealant"},
				Images:  []string{"https://images.unsplash.com/photo-1533090161767-e6ffed986c88?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1533090161767-e6ffed986c88?w=800&q=80"},
			},
			{
				ID: "furn-bergen-bed", Name: "Bergen Floating Teak Platform Bed (King)", Slug: "bergen-floating-teak-platform-bed-king",
				Emoji: "🛏️", Style: "Bedroom", Length: "6x6 King", Category: "Beds & Bedroom", Badge: "Integrated Lighting",
				Description: "Cantilevered floating aesthetic crafted in solid plantation teak with an integrated LED headboard and dual soft-close bedside ledges.",
				PriceKES:    145000, OriginalPriceKES: 165000, Currency: "KES", Featured: true, InStock: true,
				Specs:   map[string]string{"Timber": "100% Solid Plantation Teak", "Mattress Compatibility": "Standard King (180cm × 200cm)", "Lighting": "Warm 2700K Dimmable Headboard Ambient LED", "Joinery": "Concealed Japanese Mortise & Tenon"},
				Images:  []string{"https://images.unsplash.com/photo-1505693416388-ac5ce068fe85?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1505693416388-ac5ce068fe85?w=800&q=80"},
			},
			{
				ID: "furn-malmo-credenza", Name: "Malmö Ribbed Oak Credenza & Sideboard", Slug: "malmo-ribbed-oak-credenza-sideboard",
				Emoji: "🚪", Style: "Storage", Length: "200cm Console", Category: "Storage & Consoles", Badge: "Soft-Close",
				Description: "Architectural tambour ribbed sliding doors concealing adjustable shelving, integrated media cable passthroughs, and matte black steel hairpin legs.",
				PriceKES:    78000, OriginalPriceKES: 88000, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Exterior": "Ribbed Solid Oak Tambour Sliding Doors", "Hardware": "Blum Soft-Close Concealed Drawer Runners", "Dimensions": "W: 200cm × D: 48cm × H: 75cm", "Legs": "Matte Powder-Coated Architectural Steel"},
				Images:  []string{"https://images.unsplash.com/photo-1538688525198-9b88f6f53126?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1538688525198-9b88f6f53126?w=800&q=80"},
			},
			{
				ID: "furn-stockholm-desk", Name: "Stockholm Executive Solid Ash Writing Desk", Slug: "stockholm-executive-solid-ash-writing-desk",
				Emoji: "💼", Style: "Workspace", Length: "160cm Desk", Category: "Desks & Office", Badge: "Cable Concealment",
				Description: "Sleek Scandinavian executive work desk in blonde European ash with integrated leather desk mat inlay and concealed dual USB/power ports.",
				PriceKES:    68000, OriginalPriceKES: 78000, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Timber": "Hand-Planed European Ash", "Inlay": "Vegetable-Tanned Saddle Leather Blotter", "Dimensions": "L: 160cm × W: 75cm × H: 75cm", "Storage": "Dual Concealed Felt-Lined Stationery Drawers"},
				Images:  []string{"https://images.unsplash.com/photo-1518455027359-f3f8164ba6bd?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1518455027359-f3f8164ba6bd?w=800&q=80"},
			},
			{
				ID: "furn-uppsala-bookshelf", Name: "Uppsala Geometric Hardwood Bookshelf", Slug: "uppsala-geometric-hardwood-bookshelf",
				Emoji: "📚", Style: "Shelving", Length: "Tall 190cm", Category: "Storage & Consoles", Badge: "Modular",
				Description: "Asymmetrical open-grid display shelf made from solid mahogany and brass uprights. Perfect as a room divider or statement wall display.",
				PriceKES:    54000, OriginalPriceKES: 62000, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Wood": "Seasoned Solid African Mahogany", "Accents": "Brushed Solid Brass Corner Brackets", "Dimensions": "W: 120cm × D: 35cm × H: 190cm", "Configuration": "8 Staggered Display Cubbies"},
				Images:  []string{"https://images.unsplash.com/photo-1594938298603-c8148c4dae35?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1594938298603-c8148c4dae35?w=800&q=80"},
			},
			{
				ID: "furn-gothenburg-bench", Name: "Gothenburg Woven Danish Paper Cord Bench", Slug: "gothenburg-woven-danish-paper-cord-bench",
				Emoji: "🪵", Style: "Entryway & Dining", Length: "140cm Bench", Category: "Dining Tables", Badge: "Hand-Woven",
				Description: "A heritage Danish design featuring over 400 meters of unbleached continuous paper cord hand-woven over an oil-finished solid oak frame.",
				PriceKES:    29000, OriginalPriceKES: 34000, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Seat": "3-Ply Unbleached Danish Paper Cord (Hand-Laced)", "Frame": "Solid European White Oak", "Dimensions": "L: 140cm × D: 42cm × H: 46cm", "Capacity": "Tested to 300 kg Static Load"},
				Images:  []string{"https://images.unsplash.com/photo-1503602642458-232111445657?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1503602642458-232111445657?w=800&q=80"},
			},
			{
				ID: "furn-lund-counter", Name: "Lund Ergonomic Solid Walnut Bar Stool (Pair)", Slug: "lund-ergonomic-solid-walnut-bar-stool-pair",
				Emoji: "🍸", Style: "Kitchen & Bar", Length: "Counter Height 65cm", Category: "Accent Chairs", Badge: "Set of 2",
				Description: "Sculpted tractor-style wooden seat contours to the body without padding. Seamless mortise joint legs with a brushed brass footrest protection ring.",
				PriceKES:    36000, OriginalPriceKES: 42000, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Material": "Solid American Black Walnut", "Seat Height": "65cm (Designed for Kitchen Islands)", "Footrest": "Brushed Brass Anti-Scuff Protector", "Packaging": "Sold as a Pair (2 Stools)"},
				Images:  []string{"https://images.unsplash.com/photo-1503602642458-232111445657?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1503602642458-232111445657?w=800&q=80"},
			},
		}
		if err := r.DeleteAllServices(ctx); err != nil {
			return err
		}
		for i, item := range items {
			if err := r.CreateService(ctx, item, i); err != nil {
				return err
			}
		}
	}
	return nil
}

// --- Motor Vehicles Preset ---
func (r *ContentRepo) applyVehiclesPreset(ctx context.Context, seedItems bool) error {
	site := models.SiteConfig{
		Industry:           "vehicles",
		BrandName:          "Apex Motors",
		BrandScript:        "Prestige",
		Tagline:            "Uncompromising Power · Curated Exotics",
		HeroTitle:          "Luxury, Performance & Bespoke Motor Vehicles",
		HeroSubtitle:       "150-point certified exotics, executive SUVs, and custom vehicle sourcing delivered with concierge care in Kenya.",
		HeroBadge:          "◆ Certified Luxury Dealership & Concierge Sourcing — Nairobi ◆",
		HeroVideoURL:       "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4",
		HeroBackgroundType: "video",
		HeroOverlayOpacity: 60,
		Announcement:       "✦ 150-Point Certified Inspection · Trade-ins Welcome · In-House Asset Financing Available ✦",
		Location:           "Apex Automotive Showroom, Westlands, Nairobi",
		StudioAddress:      "Apex Plaza, Chiromo Road, Westlands, Nairobi",
		Phone:              "+254 711 999 888",
		Email:              "sales@apexmotors.co.ke",
		WhatsAppURL:        "https://wa.me/254711999888",
		InstagramURL:       "https://instagram.com/apexmotors.ke",
		TikTokURL:          "https://tiktok.com/@apexmotors",
		PriceRange:         "KES 6,500,000 – 45,000,000",
		Copyright:          "© 2026 Apex Motors Prestige. All rights reserved.",
		MetaDescription:    "Apex Motors Prestige — Nairobi's premier luxury vehicle dealer. Porsche, Range Rover, BMW M, Mercedes-AMG with verified inspections and trade-ins.",
		Theme: models.Theme{
			Primary:      "#0D1117",
			PrimaryDark:  "#05070A",
			Accent:       "#38BDF8",
			AccentLight:  "#BAE6FD",
			Background:   "#0B0F19",
			Surface:      "#161F30",
			Text:         "#F1F5F9",
			TextMuted:    "#94A3B8",
			HeroGradient: "linear-gradient(135deg, #05070A 0%, #0D1117 60%, #1A2333 100%)",
			FontHeading:  "'Montserrat', system-ui, sans-serif",
			FontBody:     "'Inter', system-ui, sans-serif",
			FontScript:   "'Cinzel', Georgia, serif",
		},
		Animations: models.AnimationConfig{
			GlobalEnabled:   true,
			Intensity:       "energetic",
			ScrollReveal:    true,
			DefaultDuration: "0.7s",
			CardHoverEffect: "glow",
		},
		About: models.AboutSection{
			Title: "The Apex Standard",
			Intro: "Engineering precision, documented pedigree, and exhilarating road presence.",
			Story: "Apex Motors Prestige was built for automotive connoisseurs who settle for nothing less than absolute excellence. Every vehicle in our showroom undergoes rigorous computerized diagnostics, mileage verification, and structural assessment. Whether you are acquiring your first performance sports car or expanding a bespoke collection, our private client advisors guide every milestone.",
			Highlights: []string{
				"150-point comprehensive mechanical & electrical diagnostic report",
				"Verified genuine mileage with global Carfax & Japanese auction pedigrees",
				"Bespoke international vehicle sourcing from UK, Japan, and Germany in 21 days",
				"Nationwide enclosed trailer delivery and tailored warranty coverage",
			},
			ArtistName:  "Apex Concierge Team",
			ArtistTitle: "Senior Automotive Advisors",
		},
		HouseCalls: models.HouseCallsSection{
			Title:       "Enclosed Concierge Delivery & Home Test Drives",
			Description: "We deliver vehicles for private viewing and test drives directly to your residence or corporate offices in Nairobi.",
			Notes: []string{
				"Discreet home viewing by appointment with our product specialist",
				"Enclosed transport to ensure zero road mileage or stone chips",
				"On-site computerized diagnostic handover check with the client",
				"Immediate registration and logbook transfer assistance",
			},
		},
		Aftercare: models.AftercareSection{
			Title: "Ownership & Maintenance Protocol",
			Intro: "High-performance machines demand meticulous maintenance. Our certified technicians keep your vehicle at peak output.",
			Tips: []models.AftercareTip{
				{Title: "Engine warm-up & cool-down", Description: "Allow turbocharged engines to idle for 60 seconds after spirited driving before shutdown to preserve turbo bearings."},
				{Title: "Ceramic coat protection", Description: "Every Apex vehicle receives a 3-year hydrophobic ceramic coating. Wash only with pH-neutral automotive foam."},
				{Title: "Scheduled telemetry checks", Description: "Bring your vehicle in every 7,500 km for complimentary fluid analysis and software updates."},
			},
		},
		PricingHighlight: models.PricingHighlight{
			Title:      "Transparent Vehicle Pricing",
			Subtitle:   "Fully duty-paid, CIF Nairobi pricing with zero concealed charges.",
			PriceFrom:  6500000,
			PriceTo:    45000000,
			Currency:   "KES",
			Badge:      "Certified Pedigree",
			Includes:   "Duty paid, genuine logbook, 1-year warranty, full service, and detailing.",
			StudioNote: "◆ Chiromo Road Showroom · ◆ Trade-in Evaluations Daily · ◆ Nationwide Delivery",
		},
		Payment: models.PaymentConfig{
			MpesaManualEnabled: false,
			TillNumber:         "",
			PaybillNumber:      "",
			PaybillAccount:     "",
			MpesaShortCode:     "",
			CashEnabled:        true,
			CashInstructions:   "Payment in full at our Westlands showroom upon handover.",
		},
	}
	if err := r.UpdateSite(ctx, site); err != nil {
		return err
	}

	sections := []models.Section{
		{
			ID:                "hero",
			Type:              "hero",
			Enabled:           true,
			SortOrder:         0,
			Title:             "Luxury, Performance & Bespoke Motor Vehicles",
			Subtitle:          "150-point certified exotics, executive SUVs, and custom vehicle sourcing delivered with concierge care in Kenya.",
			Badge:             "◆ Certified Luxury Dealership & Concierge Sourcing ◆",
			BackgroundType:    "video",
			BackgroundVideo:   "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4",
			OverlayOpacity:    60,
			Animation:         "zoom-in",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
			PrimaryCtaText:    "Browse Inventory",
			PrimaryCtaUrl:     "/collections",
			SecondaryCtaText:  "Book a Test Drive",
			SecondaryCtaUrl:   "/contact",
		},
		{
			ID:                "features",
			Type:              "features",
			Enabled:           true,
			SortOrder:         1,
			Title:             "Why Discerning Drivers <span class='accent'>Choose Apex</span>",
			Subtitle:          "Every vehicle rigorously verified, flawlessly detailed, and ready for the open road",
			Badge:             "~ The Apex Guarantee ~",
			BackgroundType:    "default",
			Animation:         "fade-in",
			AnimationDuration: "0.7s",
			AnimationDelay:    "0.1s",
			HoverEffect:       "glow",
		},
		{
			ID:                "products",
			Type:              "products",
			Enabled:           true,
			SortOrder:         2,
			Title:             "Showroom <span class='accent'>Inventory</span>",
			Subtitle:          "Verified mileage, certified engine diagnostics, and 360° walkarounds",
			Badge:             "~ Available Vehicles ~",
			BackgroundType:    "default",
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
			HoverEffect:       "lift",
			PrimaryCtaText:    "View Full Inventory",
			PrimaryCtaUrl:     "/collections",
		},
		{
			ID:                "video_banner",
			Type:              "video_banner",
			Enabled:           true,
			SortOrder:         3,
			Title:             "Pure Adrenaline <span class='accent'>On Demand</span>",
			Subtitle:          "Experience the growl of precision-tuned V8s and instantaneous electric torque.",
			Badge:             "~ Cinematic Experience ~",
			BackgroundType:    "video",
			BackgroundVideo:   "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4",
			OverlayOpacity:    55,
			Animation:         "zoom-in",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
			PrimaryCtaText:    "Schedule Track Day",
			PrimaryCtaUrl:     "#track-day",
		},
		{
			ID:                "reviews",
			Type:              "reviews",
			Enabled:           true,
			SortOrder:         4,
			Title:             "Driver <span class='accent'>Endorsements</span>",
			Subtitle:          "Feedback from verified clients and corporate leaders across East Africa",
			Badge:             "~ Client Stories ~",
			BackgroundType:    "default",
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
		},
		{
			ID:                "faq",
			Type:              "faq",
			Enabled:           true,
			SortOrder:         5,
			Title:             "Dealership <span class='accent'>Inquiries</span>",
			Subtitle:          "Import clearances, trade-in valuations, bank asset financing, and warranties",
			Badge:             "~ FAQ ~",
			BackgroundType:    "default",
			Animation:         "fade-in",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
		},
		{
			ID:                "inquiry",
			Type:              "inquiry",
			Enabled:           true,
			SortOrder:         6,
			Title:             "Book a Test Drive or <span class='accent'>Request Sourcing</span>",
			Subtitle:          "Our automotive advisors will prepare the vehicle and welcome you at our Westlands showroom",
			Badge:             "~ Reserve Experience ~",
			BackgroundType:    "default",
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
		},
	}
	if err := r.DeleteAllSections(ctx); err != nil {
		return err
	}
	for i, s := range sections {
		if err := r.CreateSection(ctx, s, i); err != nil {
			return err
		}
	}

	if seedItems {
		items := []models.LashService{
			{
				ID: "veh-porsche-911", Name: "Porsche 911 Carrera 4S (992)", Slug: "porsche-911-carrera-4s-992",
				Emoji: "🏎️", Style: "Sports Coupe", Length: "AWD", Category: "Sports & Performance", Badge: "Low Mileage",
				Description: "Agate Grey Metallic over Bordeaux Red leather. Equipped with Sport Chrono Package, Sport Exhaust System, and PASM Sport Suspension. Full Porsche service record.",
				PriceKES:    18500000, OriginalPriceKES: 19800000, Currency: "KES", Featured: true, InStock: true,
				Specs:    map[string]string{"Year": "2023", "Mileage": "8,400 km", "Engine": "3.0L Twin-Turbo Flat-6 (443 HP)", "Transmission": "8-Speed Dual-Clutch (PDK)", "0-100 km/h": "3.4 seconds", "Drivetrain": "All-Wheel Drive (AWD)"},
				Images:   []string{"https://images.unsplash.com/photo-1503376780353-7e6692767b70?w=800&q=80", "https://images.unsplash.com/photo-1552519507-da3b142c6e3d?w=800&q=80"},
				View360:  []string{"https://images.unsplash.com/photo-1503376780353-7e6692767b70?w=800&q=80"},
				VideoURL: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4",
			},
			{
				ID: "veh-range-autobio", Name: "Range Rover Autobiography LWB", Slug: "range-rover-autobiography-lwb",
				Emoji: "🚙", Style: "Executive SUV", Length: "Long Wheelbase", Category: "Executive SUVs", Badge: "Flagship Luxury",
				Description: "Charente Grey with Semi-Aniline Perlino leather interior. Executive Class rear seating with massage, Meridian Signature 1600W Sound System, and rear-axle steering.",
				PriceKES:    26500000, OriginalPriceKES: 28000000, Currency: "KES", Featured: true, InStock: true,
				Specs:    map[string]string{"Year": "2024", "Mileage": "4,200 km", "Engine": "4.4L Twin-Turbo V8 (523 HP)", "Transmission": "8-Speed Automatic", "Suspension": "Electronic Air with Dynamic Response", "Seating": "Executive Class 4-Seat"},
				Images:   []string{"https://images.unsplash.com/photo-1533473359331-0135ef1b58bf?w=800&q=80"},
				View360:  []string{"https://images.unsplash.com/photo-1533473359331-0135ef1b58bf?w=800&q=80"},
				VideoURL: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4",
			},
			{
				ID: "veh-bmw-m4", Name: "BMW M4 Competition Coupe xDrive", Slug: "bmw-m4-competition-coupe-xdrive",
				Emoji: "⚡", Style: "Track Coupe", Length: "M-xDrive", Category: "Sports & Performance", Badge: "Carbon Package",
				Description: "Isle of Man Green over Kyalami Orange Merino leather. Carbon fiber bucket seats, M Carbon ceramic brakes, and M Driver's package with unrestricted top speed.",
				PriceKES:    14800000, OriginalPriceKES: 16000000, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Year": "2023", "Mileage": "11,500 km", "Engine": "3.0L BMW M TwinPower Turbo Inline-6 (503 HP)", "Transmission": "8-Speed M Steptronic", "0-100 km/h": "3.5 seconds", "Brakes": "M Carbon Ceramic"},
				Images:  []string{"https://images.unsplash.com/photo-1552519507-da3b142c6e3d?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1552519507-da3b142c6e3d?w=800&q=80"},
			},
			{
				ID: "veh-amg-g63", Name: "Mercedes-AMG G63 BiTurbo V8", Slug: "mercedes-amg-g63-biturbo-v8",
				Emoji: "🔥", Style: "Super SUV", Length: "4MATIC", Category: "Executive SUVs", Badge: "Night Package",
				Description: "Obsidian Black Metallic with G manufaktur Bengal Red Nappa leather. Burmester Surround Sound, AMG Ride Control, and switchable side-exit sports exhaust.",
				PriceKES:    32000000, OriginalPriceKES: 34500000, Currency: "KES", Featured: true, InStock: true,
				Specs:   map[string]string{"Year": "2024", "Mileage": "3,100 km", "Engine": "4.0L Handcrafted AMG V8 BiTurbo (577 HP)", "Transmission": "AMG SPEEDSHIFT TCT 9-Speed", "0-100 km/h": "4.5 seconds", "Exhaust": "AMG Performance Dual Side Pipes"},
				Images:  []string{"https://images.unsplash.com/photo-1520050206274-a1ae44613e6d?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1520050206274-a1ae44613e6d?w=800&q=80"},
			},
			{
				ID: "veh-audi-rs7", Name: "Audi RS7 Sportback Carbon Black", Slug: "audi-rs7-sportback-carbon-black",
				Emoji: "🏁", Style: "Sportback GT", Length: "Quattro", Category: "Grand Tourers", Badge: "RS Dynamic Plus",
				Description: "Nardo Grey over Valcona Black honeycomb leather. Dynamic all-wheel steering, RS Sport suspension plus with Dynamic Ride Control, and Matrix Laser LED lighting.",
				PriceKES:    16500000, OriginalPriceKES: 17800000, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Year": "2023", "Mileage": "14,200 km", "Engine": "4.0L Twin-Turbo V8 MHEV (591 HP)", "Transmission": "8-Speed Tiptronic Quattro", "0-100 km/h": "3.6 seconds", "Top Speed": "305 km/h (Derestricted)"},
				Images:  []string{"https://images.unsplash.com/photo-1603584173870-7f23fdae1b7a?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1603584173870-7f23fdae1b7a?w=800&q=80"},
			},
			{
				ID: "veh-bentley-gt", Name: "Bentley Continental GT V8 Mulliner", Slug: "bentley-continental-gt-v8-mulliner",
				Emoji: "💎", Style: "Grand Tourer", Length: "AWD", Category: "Bespoke Luxury", Badge: "Mulliner Driving Spec",
				Description: "Glacier White over Hotspur and Beluga dual-tone hide with Grand Black veneers. Bentley Rotating Display, Naim 2200W Audio, and diamond-in-diamond quilting.",
				PriceKES:    38500000, OriginalPriceKES: 41000000, Currency: "KES", Featured: true, InStock: true,
				Specs:   map[string]string{"Year": "2024", "Mileage": "2,400 km", "Engine": "4.0L Twin-Turbo V8 (542 HP)", "Transmission": "8-Speed Dual Clutch", "Audio": "Naim for Bentley 2,200W", "Interior": "Hand-Stitched Mulliner Hide"},
				Images:  []string{"https://images.unsplash.com/photo-1541348263662-e0c8de4259ba?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1541348263662-e0c8de4259ba?w=800&q=80"},
			},
			{
				ID: "veh-lc300-zx", Name: "Toyota Land Cruiser 300 ZX VIP", Slug: "toyota-land-cruiser-300-zx-vip",
				Emoji: "🏔️", Style: "Luxury 4WD", Length: "7-Seater", Category: "Executive SUVs", Badge: "Armored Ready",
				Description: "Precious White Pearl with Neutral Beige leather. Modellista aero styling kit, JBL 14-speaker premium sound, rear seat entertainment, and Multi-Terrain Select.",
				PriceKES:    21500000, OriginalPriceKES: 22800000, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Year": "2024", "Mileage": "6,800 km", "Engine": "3.3L Twin-Turbo V6 Diesel (304 HP / 700 Nm)", "Transmission": "10-Speed Direct Shift Automatic", "Suspension": "Adaptive Variable Suspension (AVS)", "Styling": "Full Modellista Aero Package"},
				Images:  []string{"https://images.unsplash.com/photo-1594502184342-2e12f877aa73?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1594502184342-2e12f877aa73?w=800&q=80"},
			},
			{
				ID: "veh-aston-dbx707", Name: "Aston Martin DBX707 Super SUV", Slug: "aston-martin-dbx707-super-suv",
				Emoji: "🏆", Style: "Super SUV", Length: "AWD", Category: "Sports & Performance", Badge: "707 Horsepower",
				Description: "Satin Titanium Grey over Onyx Black Semi-Aniline leather with Lime contrast stitching. Carbon fiber aero package, sport plus seats, and soft-close doors.",
				PriceKES:    36000000, OriginalPriceKES: 39000000, Currency: "KES", Featured: true, InStock: true,
				Specs:   map[string]string{"Year": "2024", "Mileage": "1,900 km", "Engine": "4.0L Twin-Turbocharged V8 (697 HP / 900 Nm)", "Transmission": "9-Speed Wet-Clutch Automatic", "0-100 km/h": "3.3 seconds", "Top Speed": "310 km/h"},
				Images:  []string{"https://images.unsplash.com/photo-1542282088-72c9c27ed0cd?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1542282088-72c9c27ed0cd?w=800&q=80"},
			},
			{
				ID: "veh-taycan-turbo", Name: "Porsche Taycan Turbo S Electric", Slug: "porsche-taycan-turbo-s-electric",
				Emoji: "⚡", Style: "Electric GT", Length: "AWD", Category: "Electric Performance", Badge: "Zero Emissions",
				Description: "Carrara White Metallic with Olea Club leather in Basalt Black. Performance Battery Plus (93.4 kWh), Porsche Ceramic Composite Brakes, and Panoramic Fixed Glass Roof.",
				PriceKES:    19800000, OriginalPriceKES: 21500000, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Year": "2023", "Mileage": "9,200 km", "Powertrain": "Dual Permanent Magnet Synchronous Motors (750 HP)", "0-100 km/h": "2.8 seconds (Launch Control)", "Range": "412 km WLTP", "Charging": "270 kW DC Fast Charging (5-80% in 22 mins)"},
				Images:  []string{"https://images.unsplash.com/photo-1614162692292-7ac56d7f7f1e?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1614162692292-7ac56d7f7f1e?w=800&q=80"},
			},
			{
				ID: "veh-maybach-s680", Name: "Mercedes-Maybach S680 V12 First Class", Slug: "mercedes-maybach-s680-v12-first-class",
				Emoji: "👑", Style: "State Limousine", Length: "Extended Wheelbase", Category: "Bespoke Luxury", Badge: "V12 Flagship",
				Description: "Two-tone High-Tech Silver over Nautical Blue. Individual rear executive suites, refrigerated champagne compartment with silver flutes, and active road noise cancellation.",
				PriceKES:    42000000, OriginalPriceKES: 45000000, Currency: "KES", Featured: true, InStock: true,
				Specs:   map[string]string{"Year": "2024", "Mileage": "1,200 km", "Engine": "6.0L Handcrafted Twin-Turbo V12 (621 HP / 1000 Nm)", "Transmission": "9G-TRONIC 4MATIC AWD", "Rear Cabin": "First-Class Rear Suite with Leg Rests & Calf Massage", "Chassis": "E-ACTIVE BODY CONTROL with Curve Tilting"},
				Images:  []string{"https://images.unsplash.com/photo-1563720223185-11003d516935?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1563720223185-11003d516935?w=800&q=80"},
			},
		}
		if err := r.DeleteAllServices(ctx); err != nil {
			return err
		}
		for i, item := range items {
			if err := r.CreateService(ctx, item, i); err != nil {
				return err
			}
		}
	}
	return nil
}

// --- Spa Preset ---
func (r *ContentRepo) applySpaPreset(ctx context.Context, seedItems bool) error {
	site := models.SiteConfig{
		Industry:           "spa",
		BrandName:          "Serena",
		BrandScript:        "Sanctuary",
		Tagline:            "Breathe · Reconnect · Restore",
		HeroTitle:          "Holistic Botanical Healing & Thermal Wellness Sanctuary",
		HeroSubtitle:       "Awaken your senses with restorative hot stone rituals, organic phyto-facials, and private hydrotherapy suites in Nairobi.",
		HeroBadge:          "◆ Award-Winning Wellness Sanctuary & Day Spa — Nairobi ◆",
		HeroVideoURL:       "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerJoyBlazes.mp4",
		HeroBackgroundType: "video",
		HeroOverlayOpacity: 40,
		Announcement:       "✦ 100% Organic Botanical Extracts · Private Couples Suites · Complimentary Herbal Elixir ✦",
		Location:           "Serena Sanctuary, Peponi Road, Westlands, Nairobi",
		StudioAddress:      "Peponi Gardens Estate, Villa 14, Westlands, Nairobi",
		Phone:              "+254 733 445 566",
		Email:              "wellness@serenasanctuary.co.ke",
		WhatsAppURL:        "https://wa.me/254733445566",
		InstagramURL:       "https://instagram.com/serenasanctuary.ke",
		TikTokURL:          "https://tiktok.com/@serenasanctuary",
		PriceRange:         "KES 4,500 – 28,000",
		Copyright:          "© 2026 Serena Sanctuary. Crafted for wellbeing.",
		MetaDescription:    "Serena Sanctuary — Premier holistic spa in Nairobi. Hot stone massages, hydrotherapy, botanical facials, and couples wellness packages.",
		Theme: models.Theme{
			Primary:      "#2D3A34",
			PrimaryDark:  "#1E2723",
			Accent:       "#C2A676",
			AccentLight:  "#F3EDE2",
			Background:   "#FBF9F5",
			Surface:      "#FFFFFF",
			Text:         "#242B27",
			TextMuted:    "#66726C",
			HeroGradient: "linear-gradient(135deg, #1E2723 0%, #2D3A34 60%, #3D4E46 100%)",
			FontHeading:  "'Cormorant Garamond', Georgia, serif",
			FontBody:     "'DM Sans', system-ui, sans-serif",
			FontScript:   "'Playfair Display', Georgia, serif",
		},
		Animations: models.AnimationConfig{
			GlobalEnabled:   true,
			Intensity:       "subtle",
			ScrollReveal:    true,
			DefaultDuration: "1.0s",
			CardHoverEffect: "lift",
		},
		About: models.AboutSection{
			Title: "Our Sanctuary Ritual",
			Intro: "Rooted in African flora, thermal mineral waters, and ancient restorative wisdom.",
			Story: "Serena Sanctuary was born as a tranquil haven away from the rhythm of the city. We blend potent organic botanical oils extracted from Kenyan highlands with time-honored European hydrotherapy and Ayurvedic pressure rituals. Every session is personalized by master practitioners to realign mind, body, and spirit.",
			Highlights: []string{
				"Custom blended cold-pressed baobab, marula, and frankincense oils",
				"Private eucalyptus steam rooms and cedarwood soaking baths",
				"Master aestheticians with international CIDESCO accreditation",
				"Bespoke retreat packages for couples, brides, and corporate recharge",
			},
			ArtistName:  "Serena Sanctuary Team",
			ArtistTitle: "Holistic Health Practitioners",
		},
		HouseCalls: models.HouseCallsSection{
			Title:       "Private Residence & Executive Spa Visits",
			Description: "Experience the sanctuary in your private home, embassy residence, or luxury boutique suite.",
			Notes: []string{
				"Fully portable heated orthopedic treatment tables and Egyptian linen",
				"Ambient soundscapes and aromatherapy diffusers set up in minutes",
				"Available for individuals, private couples, and bridal parties",
				"Discreet booking with 24-hour advance confirmation",
			},
		},
		Aftercare: models.AftercareSection{
			Title: "Post-Treatment Wellbeing Guidelines",
			Intro: "Maximize the therapeutic benefits of your massage or facial ritual.",
			Tips: []models.AftercareTip{
				{Title: "Hydrate deeply", Description: "Drink at least 2 liters of room temperature water or herbal infusion to support your body's lymphatic flushing."},
				{Title: "Allow botanical oils to absorb", Description: "Wait at least 4 hours before taking a hot shower so the nourishing plant actives penetrate deeply into your dermis."},
				{Title: "Restful evening", Description: "Avoid strenuous workouts or caffeine immediately following thermal bodywork."},
			},
		},
		PricingHighlight: models.PricingHighlight{
			Title:      "Transparent Wellness Investment",
			Subtitle:   "Pure botanical luxury with zero hidden surcharges.",
			PriceFrom:  4500,
			PriceTo:    28000,
			Currency:   "KES",
			Badge:      "All-Inclusive Rituals",
			Includes:   "Organic oils, eucalyptus steam bath, cedar sauna, and herbal tea ceremony.",
			StudioNote: "◆ Peponi Road Villa · ◆ Couples Suites Available · ◆ Complimentary Steam Bath",
		},
		Payment: models.PaymentConfig{
			MpesaManualEnabled: false,
			TillNumber:         "",
			PaybillNumber:      "",
			PaybillAccount:     "",
			MpesaShortCode:     "",
			CashEnabled:        true,
			CashInstructions:   "Cash accepted at the reception desk upon check-in or after treatment.",
		},
	}
	if err := r.UpdateSite(ctx, site); err != nil {
		return err
	}

	sections := []models.Section{
		{
			ID:                "hero",
			Type:              "hero",
			Enabled:           true,
			SortOrder:         0,
			Title:             "Holistic Botanical Healing & Thermal Wellness Sanctuary",
			Subtitle:          "Awaken your senses with restorative hot stone rituals, organic phyto-facials, and private hydrotherapy suites in Nairobi.",
			Badge:             "◆ Award-Winning Wellness Sanctuary & Day Spa ◆",
			BackgroundType:    "video",
			BackgroundVideo:   "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerJoyBlazes.mp4",
			OverlayOpacity:    40,
			Animation:         "fade-in",
			AnimationDuration: "1.0s",
			AnimationDelay:    "0.1s",
			PrimaryCtaText:    "View Rituals & Treatments",
			PrimaryCtaUrl:     "/collections",
			SecondaryCtaText:  "Reserve on WhatsApp",
			SecondaryCtaUrl:   "https://wa.me/254733445566",
		},
		{
			ID:                "features",
			Type:              "features",
			Enabled:           true,
			SortOrder:         1,
			Title:             "The Serena <span class='accent'>Sanctuary Pillars</span>",
			Subtitle:          "Ancient wisdom harmonized with modern botanical skin science",
			Badge:             "~ Our Wellness Ethos ~",
			BackgroundType:    "default",
			Animation:         "fade-in",
			AnimationDuration: "0.9s",
			AnimationDelay:    "0.1s",
			HoverEffect:       "lift",
		},
		{
			ID:                "products",
			Type:              "products",
			Enabled:           true,
			SortOrder:         2,
			Title:             "Signature <span class='accent'>Sanctuary Rituals</span>",
			Subtitle:          "Individually customized treatment durations, organic essential oils, and restorative benefits",
			Badge:             "~ Treatment Menu ~",
			BackgroundType:    "default",
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
			HoverEffect:       "lift",
			PrimaryCtaText:    "View Full Treatment Menu",
			PrimaryCtaUrl:     "/collections",
		},
		{
			ID:                "video_banner",
			Type:              "video_banner",
			Enabled:           true,
			SortOrder:         3,
			Title:             "Immerse in <span class='accent'>Gentle Serenity</span>",
			Subtitle:          "Step into a world of cascading water, warm volcanic stones, and calming aromatic steam.",
			Badge:             "~ Sensory Retreat ~",
			BackgroundType:    "video",
			BackgroundVideo:   "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerJoyBlazes.mp4",
			OverlayOpacity:    45,
			Animation:         "fade-in",
			AnimationDuration: "1.0s",
			AnimationDelay:    "0.1s",
			PrimaryCtaText:    "Discover Sanctuary Suites",
			PrimaryCtaUrl:     "#sanctuary-suites",
		},
		{
			ID:                "reviews",
			Type:              "reviews",
			Enabled:           true,
			SortOrder:         4,
			Title:             "Client <span class='accent'>Reflections</span>",
			Subtitle:          "Read how our treatments have helped guests find tranquility and restoration",
			Badge:             "~ Verified Testimonials ~",
			BackgroundType:    "default",
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
		},
		{
			ID:                "faq",
			Type:              "faq",
			Enabled:           true,
			SortOrder:         5,
			Title:             "Sanctuary <span class='accent'>Etiquette & FAQ</span>",
			Subtitle:          "Arrival times, private suite reservations, cancellation guidelines, and custom gift vouchers",
			Badge:             "~ Helpful Information ~",
			BackgroundType:    "default",
			Animation:         "fade-in",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
		},
		{
			ID:                "inquiry",
			Type:              "inquiry",
			Enabled:           true,
			SortOrder:         6,
			Title:             "Reserve Your <span class='accent'>Sanctuary Session</span>",
			Subtitle:          "Select your preferred treatment, date, and let our practitioners prepare your private suite",
			Badge:             "~ Book Appointment ~",
			BackgroundType:    "default",
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
		},
	}
	if err := r.DeleteAllSections(ctx); err != nil {
		return err
	}
	for i, s := range sections {
		if err := r.CreateSection(ctx, s, i); err != nil {
			return err
		}
	}

	if seedItems {
		items := []models.LashService{
			{
				ID: "spa-imperial-royal", Name: "The Imperial Royal Ritual", Slug: "the-imperial-royal-ritual",
				Emoji: "👑", Style: "Body Ritual", Length: "120 Mins", Category: "Signature Rituals", Badge: "Most Popular",
				Description: "A lavish full-body hot stone massage accompanied by Moroccan eucalyptus black soap scrub, private cedar soaking tub, and custom 24K gold hydration mask.",
				PriceKES:    12500, OriginalPriceKES: 14500, Currency: "KES", Featured: true, InStock: true,
				Specs:    map[string]string{"Duration": "120 Minutes", "Inclusions": "Hydrotherapy Soak, Hot Stones, Scalp Treatment", "Aromatherapy": "Organic Damask Rose & Wild Sandalwood", "Therapist": "Senior Master Holistic Therapist"},
				Images:   []string{"https://images.unsplash.com/photo-1540555700478-4be289fbecef?w=800&q=80"},
				View360:  []string{"https://images.unsplash.com/photo-1540555700478-4be289fbecef?w=800&q=80"},
				VideoURL: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4",
			},
			{
				ID: "spa-ancestral-herbal", Name: "Ancestral Heated Herbal Compress Massage", Slug: "ancestral-heated-herbal-compress-massage",
				Emoji: "🌿", Style: "Holistic Therapy", Length: "90 Mins", Category: "Holistic Massages", Badge: "Deep Recovery",
				Description: "Steamed pouches of lemongrass, turmeric, camphor, and indigenous botanical roots pressed along meridian pathways to dissolve chronic tension.",
				PriceKES:    7800, OriginalPriceKES: 9000, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Duration": "90 Minutes", "Technique": "Thai Herbal Poultice + Swedish Acupressure", "Benefits": "Lymphatic Circulation & Muscle Inflammation Relief", "Includes": "Warm Ginger Botanical Elixir"},
				Images:  []string{"https://images.unsplash.com/photo-1519823551278-64ac92734fb1?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1519823551278-64ac92734fb1?w=800&q=80"},
			},
			{
				ID: "spa-botanical-facial", Name: "Botanical Radiance 24K Gold Facial", Slug: "botanical-radiance-24k-gold-facial",
				Emoji: "✨", Style: "Clinical Facial", Length: "75 Mins", Category: "Clinical Facials", Badge: "Cellular Glow",
				Description: "Ultrasonic exfoliation followed by pure 24-karat gold leaf sheeting, bio-placenta serum infusion, and chilled rose quartz contour sculpting.",
				PriceKES:    9500, OriginalPriceKES: 11000, Currency: "KES", Featured: true, InStock: true,
				Specs:   map[string]string{"Duration": "75 Minutes", "Actives": "Pure 24K Gold Leaf & Triple Collagen Peptides", "Device": "Ultrasonic Skin Scrubber + Cryo Wand", "Skin Type": "Dull, Dehydrated, or Mature Skin"},
				Images:  []string{"https://images.unsplash.com/photo-1570172619644-dfd03ed5d881?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1570172619644-dfd03ed5d881?w=800&q=80"},
			},
			{
				ID: "spa-deep-tissue-cbd", Name: "Deep Tissue Recovery & CBD Muscle Melt", Slug: "deep-tissue-recovery-cbd-muscle-melt",
				Emoji: "💆‍♂️", Style: "Therapeutic", Length: "90 Mins", Category: "Holistic Massages", Badge: "Athletic Release",
				Description: "Targeted myofascial release combined with broad-spectrum organic CBD balm to penetrate knotted fascia, soothe joint soreness, and restore mobility.",
				PriceKES:    8500, OriginalPriceKES: 9800, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Duration": "90 Minutes", "Pressure": "Firm to Deep Myofascial Release", "Formulation": "500mg Organic Broad-Spectrum CBD Oil", "Focus": "Lower Back, Shoulders, Neck & Hamstrings"},
				Images:  []string{"https://images.unsplash.com/photo-1600334129128-685c5582fd35?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1600334129128-685c5582fd35?w=800&q=80"},
			},
			{
				ID: "spa-couples-sanctuary", Name: "Signature Couples Sanctuary Indulgence", Slug: "signature-couples-sanctuary-indulgence",
				Emoji: "🕯️", Style: "Couples Package", Length: "150 Mins", Category: "Couples Packages", Badge: "Private Suite",
				Description: "Exclusive reservation of our private twin sanctuary suite. Includes dual dry body brushing, 90-min synchronized candle massages, private whirlpool bath, and Moët & Chandon champagne.",
				PriceKES:    18500, OriginalPriceKES: 22000, Currency: "KES", Featured: true, InStock: true,
				Specs:   map[string]string{"Duration": "150 Minutes (2.5 Hours)", "Suite": "Private Deluxe Suite with Whirlpool & Shower", "Inclusions": "Chilled Champagne & Belgian Chocolate Strawberries", "Therapists": "Two Synchronized Dedicated Practitioners"},
				Images:  []string{"https://images.unsplash.com/photo-1540555700478-4be289fbecef?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1540555700478-4be289fbecef?w=800&q=80"},
			},
			{
				ID: "spa-salt-hydrotherapy", Name: "Himalayan Salt Detox Hydrotherapy Wrap", Slug: "himalayan-salt-detox-hydrotherapy-wrap",
				Emoji: "🧂", Style: "Body Wrap", Length: "60 Mins", Category: "Body Treatments", Badge: "Detoxifying",
				Description: "Full body exfoliation using pink Himalayan crystal salts infused with grapefruit and cypress essential oils, followed by a thermal cocoon wrap.",
				PriceKES:    8200, OriginalPriceKES: 9500, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Duration": "60 Minutes", "Exfoliant": "Crushed 84-Mineral Himalayan Pink Salt", "Wrap": "Thermal Infrared Detoxifying Cocoon", "Post-Care": "Ultra-Rich Shea & Baobab Body Butter Glaze"},
				Images:  []string{"https://images.unsplash.com/photo-1515377905703-c4788e51af15?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1515377905703-c4788e51af15?w=800&q=80"},
			},
			{
				ID: "spa-shirodhara-flow", Name: "Ayurvedic Shirodhara & Scalp Elixir Journey", Slug: "ayurvedic-shirodhara-scalp-elixir-journey",
				Emoji: "🧘‍♀️", Style: "Ayurvedic", Length: "75 Mins", Category: "Holistic Massages", Badge: "Stress Eraser",
				Description: "A continuous rhythmic stream of warm herbal Brahmi oil poured onto the forehead chakra (Ajna) to calm the nervous system, quiet mental chatter, and nourish hair roots.",
				PriceKES:    7500, OriginalPriceKES: 8500, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Duration": "75 Minutes", "Protocol": "Gentle Head & Neck Marma Massage + 45-Min Shirodhara Stream", "Oil": "Authentic Organic Brahmi & Sesame Infusion", "Outcome": "Deep Sleep Induction & Complete Nervous Reset"},
				Images:  []string{"https://images.unsplash.com/photo-1519823551278-64ac92734fb1?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1519823551278-64ac92734fb1?w=800&q=80"},
			},
			{
				ID: "spa-moroccan-hammam", Name: "Traditional Moroccan Rhassoul Hammam Ceremony", Slug: "traditional-moroccan-rhassoul-hammam-ceremony",
				Emoji: "🏺", Style: "Hammam Ritual", Length: "90 Mins", Category: "Signature Rituals", Badge: "Authentic Bath",
				Description: "Heated marble slab treatment with black olive Beldi soap cleansing, vigorous Kessa glove exfoliation, and pure Atlas Mountain rhassoul clay body mask.",
				PriceKES:    11000, OriginalPriceKES: 13000, Currency: "KES", Featured: true, InStock: true,
				Specs:   map[string]string{"Duration": "90 Minutes", "Room": "Heated Steam Chamber with Marble Hammam Slab", "Scrub": "Authentic Moroccan Kessa Exfoliating Mitten", "Mud": "Atlas Mountain Mineral-Rich Rhassoul Lava Clay"},
				Images:  []string{"https://images.unsplash.com/photo-1540555700478-4be289fbecef?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1540555700478-4be289fbecef?w=800&q=80"},
			},
			{
				ID: "spa-cryo-contour", Name: "Cryo-Firming Lymphatic Drainage Contour", Slug: "cryo-firming-lymphatic-drainage-contour",
				Emoji: "❄️", Style: "High-Tech Facial", Length: "60 Mins", Category: "Clinical Facials", Badge: "Instant Sculpt",
				Description: "Sub-zero controlled cryo-stimulation to tighten pores, eliminate water retention, and stimulate micro-circulation for an immediate red-carpet lifted jawline.",
				PriceKES:    9200, OriginalPriceKES: 10500, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Duration": "60 Minutes", "Technology": "Medical-Grade Cryo-Probe (-10°C to -18°C)", "Drainage": "Manual Vodder Lymphatic Drainage Technique", "Immediate Effect": "Reduced Puffiness, Visible Jawline Definition"},
				Images:  []string{"https://images.unsplash.com/photo-1570172619644-dfd03ed5d881?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1570172619644-dfd03ed5d881?w=800&q=80"},
			},
			{
				ID: "spa-sound-bath", Name: "Chakra Balancing & Sound Bath Meditation Ritual", Slug: "chakra-balancing-sound-bath-meditation-ritual",
				Emoji: "🔔", Style: "Sound Healing", Length: "60 Mins", Category: "Mind & Spirit", Badge: "Total Serenity",
				Description: "Immersion in sound waves from 7 tuned Tibetan singing bowls and German crystal bowls, paired with frankincense anointing on pulse points.",
				PriceKES:    6500, OriginalPriceKES: 7500, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Duration": "60 Minutes", "Instruments": "432Hz Quartz Crystal & 7-Metal Tibetan Bowls", "Complementary": "Smudge Cleanse with Organic White Sage & Sweetgrass", "Setting": "Low-Light Candlelit Acoustic Sanctuary"},
				Images:  []string{"https://images.unsplash.com/photo-1506126613408-eca07ce68773?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1506126613408-eca07ce68773?w=800&q=80"},
			},
		}
		if err := r.DeleteAllServices(ctx); err != nil {
			return err
		}
		for i, item := range items {
			if err := r.CreateService(ctx, item, i); err != nil {
				return err
			}
		}
	}
	return nil
}

// --- Beauty & Cosmetics Preset ---
func (r *ContentRepo) applyBeautyPreset(ctx context.Context, seedItems bool) error {
	site := models.SiteConfig{
		Industry:           "beauty",
		BrandName:          "Aura Glow",
		BrandScript:        "Beauty Lab",
		Tagline:            "Clean Clinical Actives · Radiant Skin",
		HeroTitle:          "High-Performance Clean Skincare & Glow Cosmetics",
		HeroSubtitle:       "Formulated with biocompatible peptides, multi-molecular hyaluronic acid, and botanical antioxidants for glass-skin luminosity.",
		HeroBadge:          "◆ Clean Science · Cruelty-Free · Dermatologist Formulated ◆",
		HeroVideoURL:       "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/WeAreGoingOnBullrun.mp4",
		HeroBackgroundType: "video",
		HeroOverlayOpacity: 35,
		Announcement:       "✦ Complimentary Dewy Cleanser on orders over KES 5,000 · Nationwide Delivery ✦",
		Location:           "Aura Beauty Studio, Sarit Centre, Westlands, Nairobi",
		StudioAddress:      "Sarit Centre, Lower Ground Wing, Westlands, Nairobi",
		Phone:              "+254 718 223 344",
		Email:              "glow@auraglowlab.co.ke",
		WhatsAppURL:        "https://wa.me/254718223344",
		InstagramURL:       "https://instagram.com/auraglowlab.ke",
		TikTokURL:          "https://tiktok.com/@auraglowlab",
		PriceRange:         "KES 2,200 – 16,500",
		Copyright:          "© 2026 Aura Glow Lab. Clean Clinical Beauty.",
		MetaDescription:    "Aura Glow Lab — Clean clinical skincare and cosmetics in Nairobi. Hydrating peptide serums, barrier creams, and dewy lip tints with fast delivery.",
		Theme: models.Theme{
			Primary:      "#2C1B2D",
			PrimaryDark:  "#1A0F1B",
			Accent:       "#E08F84",
			AccentLight:  "#FDECEB",
			Background:   "#FDF8F7",
			Surface:      "#FFFFFF",
			Text:         "#291823",
			TextMuted:    "#7A6172",
			HeroGradient: "linear-gradient(135deg, #1A0F1B 0%, #2C1B2D 55%, #4D2D49 100%)",
			FontHeading:  "'Playfair Display', Georgia, serif",
			FontBody:     "'DM Sans', system-ui, sans-serif",
			FontScript:   "'Cormorant Garamond', Georgia, serif",
		},
		Animations: models.AnimationConfig{
			GlobalEnabled:   true,
			Intensity:       "balanced",
			ScrollReveal:    true,
			DefaultDuration: "0.8s",
			CardHoverEffect: "lift",
		},
		About: models.AboutSection{
			Title: "The Science of Glow",
			Intro: "Biocompatible formulas engineered for African sun and high-altitude living.",
			Story: "Aura Glow Lab was founded by cosmetic scientists and clinical dermatologists who wanted clean, potent formulations that deliver visible epidermal radiance without irritation. Every formula is free from silicones, artificial parabens, and microplastics, relying instead on encapsulated vitamins, bio-fermented botanicals, and cellular peptides.",
			Highlights: []string{
				"100% cruelty-free, vegan certified, and dermatologist tested",
				"Barrier-repairing multi-ceramides and multi-molecular hyaluronic weights",
				"Eco-friendly recyclable frosted glass packaging",
				"Virtual skin consultations with licensed aesthetic specialists",
			},
			ArtistName:  "Aura Formulation Lab",
			ArtistTitle: "Clinical Skincare Specialists",
		},
		HouseCalls: models.HouseCallsSection{
			Title:       "Express Nationwide Delivery & Consultations",
			Description: "We deliver directly to your doorstep in protective thermal packaging to ensure active ingredients maintain optimal potency.",
			Notes: []string{
				"Same-day dispatch for orders within Nairobi placed before 2 PM",
				"Insured door-to-door courier across Kenya (1-2 business days)",
				"Complimentary deluxe travel samples with every order",
				"Hassle-free 14-day formula exchange guarantee",
			},
		},
		Aftercare: models.AftercareSection{
			Title: "Skin Routine & Active Layering Protocol",
			Intro: "Maximize your active serum results with our dermatologist-approved layering sequence.",
			Tips: []models.AftercareTip{
				{Title: "Layer by viscosity", Description: "Apply your hydrating watery toners first, followed by active peptide serums, and lock in moisture with our ceramide cream."},
				{Title: "SPF is mandatory", Description: "Always finish your morning skincare routine with broad-spectrum SPF 50+ to protect newly revealed skin cells."},
				{Title: "Patch test active acids", Description: "Introduce chemical exfoliants gradually — starting twice weekly in the evening."},
			},
		},
		PricingHighlight: models.PricingHighlight{
			Title:      "Transparent Clean Beauty",
			Subtitle:   "Clinical-grade active ingredients at honest, direct-to-consumer prices.",
			PriceFrom:  2200,
			PriceTo:    16500,
			Currency:   "KES",
			Badge:      "Clinical Direct",
			Includes:   "Fresh batch guarantee, thermal packaging, and personalized skin consult.",
			StudioNote: "◆ Sarit Centre Studio · ◆ Free Consultations on WhatsApp · ◆ Same-Day Nairobi Delivery",
		},
		Payment: models.PaymentConfig{
			MpesaManualEnabled: false,
			TillNumber:         "",
			PaybillNumber:      "",
			PaybillAccount:     "",
			MpesaShortCode:     "",
			CashEnabled:        true,
			CashInstructions:   "Cash accepted upon in-store collection or courier handover.",
		},
	}
	if err := r.UpdateSite(ctx, site); err != nil {
		return err
	}

	sections := []models.Section{
		{
			ID:                "hero",
			Type:              "hero",
			Enabled:           true,
			SortOrder:         0,
			Title:             "High-Performance Clean Skincare & Glow Cosmetics",
			Subtitle:          "Formulated with biocompatible peptides, multi-molecular hyaluronic acid, and botanical antioxidants for glass-skin luminosity.",
			Badge:             "◆ Clean Science · Cruelty-Free · Dermatologist Formulated ◆",
			BackgroundType:    "video",
			BackgroundVideo:   "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/WeAreGoingOnBullrun.mp4",
			OverlayOpacity:    35,
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
			PrimaryCtaText:    "Shop Radiant Actives",
			PrimaryCtaUrl:     "/collections",
			SecondaryCtaText:  "Skin Consultation",
			SecondaryCtaUrl:   "https://wa.me/254718223344",
		},
		{
			ID:                "features",
			Type:              "features",
			Enabled:           true,
			SortOrder:         1,
			Title:             "The Clean <span class='accent'>Formulation Promise</span>",
			Subtitle:          "Potent biocompatible ingredients with zero fillers or irritants",
			Badge:             "~ Our Standards ~",
			BackgroundType:    "default",
			Animation:         "fade-in",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
			HoverEffect:       "lift",
		},
		{
			ID:                "products",
			Type:              "products",
			Enabled:           true,
			SortOrder:         2,
			Title:             "Daily <span class='accent'>Skin Essentials</span>",
			Subtitle:          "Targeted serums, restorative barrier moisturizers, and dewy finishes",
			Badge:             "~ Active Formulas ~",
			BackgroundType:    "default",
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
			HoverEffect:       "lift",
			PrimaryCtaText:    "Browse All Products",
			PrimaryCtaUrl:     "/collections",
		},
		{
			ID:                "video_banner",
			Type:              "video_banner",
			Enabled:           true,
			SortOrder:         3,
			Title:             "Glass Skin in <span class='accent'>Every Drop</span>",
			Subtitle:          "Watch the immediate hydration bounce and light-refracting glow on clean skin.",
			Badge:             "~ Dewy Glow in Motion ~",
			BackgroundType:    "video",
			BackgroundVideo:   "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/WeAreGoingOnBullrun.mp4",
			OverlayOpacity:    40,
			Animation:         "zoom-in",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
			PrimaryCtaText:    "View Application Ritual",
			PrimaryCtaUrl:     "#skin-ritual",
		},
		{
			ID:                "reviews",
			Type:              "reviews",
			Enabled:           true,
			SortOrder:         4,
			Title:             "Loved by <span class='accent'>Glow Seekers</span>",
			Subtitle:          "Real verified transformations and dermatologist reviews",
			Badge:             "~ Real Results ~",
			BackgroundType:    "default",
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
		},
		{
			ID:                "faq",
			Type:              "faq",
			Enabled:           true,
			SortOrder:         5,
			Title:             "Skincare <span class='accent'>Consultation & FAQ</span>",
			Subtitle:          "Ingredients, layering, sensitive skin advice, and expiration dates",
			Badge:             "~ Skin Advice ~",
			BackgroundType:    "default",
			Animation:         "fade-in",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
		},
		{
			ID:                "inquiry",
			Type:              "inquiry",
			Enabled:           true,
			SortOrder:         6,
			Title:             "Order or Request <span class='accent'>Custom Skin Plan</span>",
			Subtitle:          "Tell us your skin goals or place your order directly for prompt same-day dispatch",
			Badge:             "~ Order Desk ~",
			BackgroundType:    "default",
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
		},
	}
	if err := r.DeleteAllSections(ctx); err != nil {
		return err
	}
	for i, s := range sections {
		if err := r.CreateSection(ctx, s, i); err != nil {
			return err
		}
	}

	if seedItems {
		items := []models.LashService{
			{
				ID: "beau-peptide-glaze", Name: "Multi-Molecular Peptide Glaze Serum", Slug: "multi-molecular-peptide-glaze-serum",
				Emoji: "✨", Style: "Serum", Length: "50ml", Category: "Face Serums", Badge: "Cult Favorite",
				Description: "Concentrated matrixyl 3000, copper peptides, and 5 weights of hyaluronic acid delivering plump, glass-like bouncy radiance in seconds.",
				PriceKES:    4800, OriginalPriceKES: 5500, Currency: "KES", Featured: true, InStock: true,
				Specs:    map[string]string{"Volume": "50ml Pipette Dropper Bottle", "Key Actives": "Copper Tripeptide-1, Matrixyl 3000, Hyaluronic Acid (5 Molecular Weights)", "Skin Type": "All Skin Types, Dehydrated, Aging", "Usage": "Morning and Night on Damp Skin"},
				Images:   []string{"https://images.unsplash.com/photo-1620916566398-39f1143ab7be?w=800&q=80"},
				View360:  []string{"https://images.unsplash.com/photo-1620916566398-39f1143ab7be?w=800&q=80"},
				VideoURL: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4",
			},
			{
				ID: "beau-ceramide-cream", Name: "Barrier Restorative Ceramides Recovery Cream", Slug: "barrier-restorative-ceramides-recovery-cream",
				Emoji: "🧴", Style: "Moisturiser", Length: "65ml", Category: "Moisturisers", Badge: "Dermatologist Tested",
				Description: "Bio-identical ceramide complex (EOP, NP, AP) with squalane, colloidal oat extract, and cholesterol to rebuild compromised skin barriers overnight.",
				PriceKES:    4200, OriginalPriceKES: 4900, Currency: "KES", Featured: true, InStock: true,
				Specs:   map[string]string{"Volume": "65ml Airless Pump Jar", "Actives": "Ceramides 1, 3, 6-II, 5% Squalane, 2% Colloidal Oatmeal", "Texture": "Rich Velvet Whip, Zero Greasiness", "Fragrance": "100% Fragrance-Free & Hypoallergenic"},
				Images:  []string{"https://images.unsplash.com/photo-1598440947619-2c35fc9aa908?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1598440947619-2c35fc9aa908?w=800&q=80"},
			},
			{
				ID: "beau-matte-lipstick", Name: "Velvet Satin Cloud Matte Lipstick", Slug: "velvet-satin-cloud-matte-lipstick",
				Emoji: "💄", Style: "Lip Color", Length: "3.8g", Category: "Lip Formulations", Badge: "Clean Formula",
				Description: "Weightless, transfer-resistant matte lip colour infused with jojoba oil and vitamin E. Delivers saturated pigment with an airy cloud finish.",
				PriceKES:    2200, OriginalPriceKES: 2600, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Weight": "3.8g Magnetic Bullet", "Finish": "Soft Blur Velvet Matte", "Wear Time": "Up to 12 Hours Non-Drying Comfort", "Key Ingredients": "Cold-Pressed Jojoba Oil, Candelilla Wax, Vitamin E"},
				Images:  []string{"https://images.unsplash.com/photo-1586495777744-4413f21062fa?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1586495777744-4413f21062fa?w=800&q=80"},
			},
			{
				ID: "beau-vitamin-c", Name: "20% Vitamin C + Ferulic Illuminating Concentrate", Slug: "20-percent-vitamin-c-ferulic-illuminating-concentrate",
				Emoji: "🍊", Style: "Serum", Length: "30ml", Category: "Face Serums", Badge: "Potent Antioxidant",
				Description: "Pure L-Ascorbic Acid stabilized with 1% Ferulic Acid and Vitamin E. Erases hyperpigmentation, shields against UV damage, and evens skin tone.",
				PriceKES:    5400, OriginalPriceKES: 6200, Currency: "KES", Featured: true, InStock: true,
				Specs:   map[string]string{"Volume": "30ml Dark Amber Dropper", "Actives": "20% Pure L-Ascorbic Acid, 1% Ferulic Acid, 1% Alpha-Tocopherol", "pH": "3.2 for Maximum Bioavailability", "Benefit": "Clinically Proven Hyperpigmentation Fading"},
				Images:  []string{"https://images.unsplash.com/photo-1608248597359-5777717d23d8?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1608248597359-5777717d23d8?w=800&q=80"},
			},
			{
				ID: "beau-retinol-active", Name: "Retinol 1.0% Micro-Encapsulated Night Elixir", Slug: "retinol-1-micro-encapsulated-night-elixir",
				Emoji: "🌙", Style: "Night Serum", Length: "30ml", Category: "Face Serums", Badge: "Clinical Grade",
				Description: "Slow-release lipid sphere retinol that stimulates fibroblast collagen synthesis without the redness or flaking of traditional vitamin A.",
				PriceKES:    5800, OriginalPriceKES: 6600, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Volume": "30ml Airless Vacuum Dispenser", "Active": "1.0% Micro-Encapsulated Bio-Retinol + Bakuchiol", "Mechanism": "Gradual 8-Hour Timed Release", "Use": "Nightly Application Before Creams"},
				Images:  []string{"https://images.unsplash.com/photo-1601049541289-9b1b7bbbfe19?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1601049541289-9b1b7bbbfe19?w=800&q=80"},
			},
			{
				ID: "beau-aha-bha-peel", Name: "Glycolic & Salicylic Radiance Resurfacing Peel", Slug: "glycolic-salicylic-radiance-resurfacing-peel",
				Emoji: "🧪", Style: "Peel Solution", Length: "30ml", Category: "Cleansers & Exfoliants", Badge: "Weekly Treatment",
				Description: "10% Glycolic + 2% Salicylic Acid solution that gently dissolves dead surface keratin, clears congested pores, and reveals a glass skin glow.",
				PriceKES:    3900, OriginalPriceKES: 4500, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Volume": "30ml Glass Bottle", "Composition": "10% Glycolic Acid (AHA) + 2% Salicylic Acid (BHA) + Tasmanian Pepperberry", "Application": "10-Minute Wash-Off Mask (1-2x Weekly)", "Benefit": "Blackhead Removal & Texture Smoothing"},
				Images:  []string{"https://images.unsplash.com/photo-1556228720-195a672e8a03?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1556228720-195a672e8a03?w=800&q=80"},
			},
			{
				ID: "beau-mineral-spf50", Name: "Invisible Zinc Broad Spectrum Mineral Shield SPF 50+", Slug: "invisible-zinc-broad-spectrum-mineral-shield-spf-50",
				Emoji: "☀️", Style: "Sun Defense", Length: "50ml", Category: "Sun & Environmental Defense", Badge: "Zero White Cast",
				Description: "Non-nano zinc oxide sunscreen with blue light defense, niacinamide, and ectoin. Leaves a non-greasy dewy satin base under makeup.",
				PriceKES:    3600, OriginalPriceKES: 4200, Currency: "KES", Featured: true, InStock: true,
				Specs:   map[string]string{"Volume": "50ml Matte Tube", "Filter": "21% Non-Nano Micronized Zinc Oxide (Mineral Only)", "Protection": "PA++++ Broad Spectrum UVA/UVB + HEV Blue Light", "Finish": "Transparent Dewy Satin, Zero Ashy Cast on Deep Skin Tones"},
				Images:  []string{"https://images.unsplash.com/photo-1567928815117-ac75a8946761?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1567928815117-ac75a8946761?w=800&q=80"},
			},
			{
				ID: "beau-eye-contour", Name: "Phyto-Stem Cell Caffeine Depuffing Eye Matrix", Slug: "phyto-stem-cell-caffeine-depuffing-eye-matrix",
				Emoji: "👁️", Style: "Eye Treatment", Length: "15ml", Category: "Eye & Lip Care", Badge: "Cooling Ceramic Tip",
				Description: "5% Green tea caffeine complex, Swiss apple stem cells, and light-reflecting mica that visually lifts eyelids and fades dark under-eye shadows.",
				PriceKES:    4100, OriginalPriceKES: 4700, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Volume": "15ml Tube with Cryo-Ceramic Cooling Wand Applicator", "Actives": "5% Caffeine, Malus Domestica Fruit Cell Culture, Haloxyl Peptide", "Target": "Dark Circles, Fluid Retention, Fine Expression Lines", "Use": "Morning & Evening with Gentle Lymphatic Tapping"},
				Images:  []string{"https://images.unsplash.com/photo-1512290900672-1f41639d48b7?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1512290900672-1f41639d48b7?w=800&q=80"},
			},
			{
				ID: "beau-rosehip-oil", Name: "Organic Cold-Pressed Chilean Rosehip Squalane Oil", Slug: "organic-cold-pressed-chilean-rosehip-squalane-oil",
				Emoji: "🌹", Style: "Face Oil", Length: "30ml", Category: "Facial Oils", Badge: "100% Pure Organic",
				Description: "Virgin Chilean rosehip seed oil rich in provitamin A and omegas 3, 6, and 9. Blended with olive squalane for rapid absorption and supple elasticity.",
				PriceKES:    3400, OriginalPriceKES: 3900, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Volume": "30ml UV-Protected Glass Bottle", "Source": "Wild-Harvested Chilean Rosa Canina Seeds (Cold-Pressed)", "Lipids": "Linoleic Acid (44%), Alpha-Linolenic Acid (34%), Oleic Acid (14%)", "Comedogenic Rating": "1 (Non-Clogging Dry Oil)"},
				Images:  []string{"https://images.unsplash.com/photo-1608248597359-5777717d23d8?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1608248597359-5777717d23d8?w=800&q=80"},
			},
			{
				ID: "beau-cleansing-balm", Name: "Botanical Enzyme Melting Cleansing Balm", Slug: "botanical-enzyme-melting-cleansing-balm",
				Emoji: "🍯", Style: "Cleanser", Length: "100g", Category: "Cleansers & Exfoliants", Badge: "Double Cleanse",
				Description: "Solid sherbet balm that melts into a silken oil upon skin contact, dissolving waterproof mascara and SPF before emulsifying clean with warm water.",
				PriceKES:    3200, OriginalPriceKES: 3700, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Volume": "100g Aluminum Jar with Magnetic Spatula", "Base": "Sunflower Seed Oil, Moringa Butter, Sweet Almond Oil", "Enzymes": "Papaya Fruit Enzyme & Green Tea Leaf Extract", "Rinse": "100% Water-Soluble Emulsification, Zero Waxy Residue"},
				Images:  []string{"https://images.unsplash.com/photo-1556228720-195a672e8a03?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1556228720-195a672e8a03?w=800&q=80"},
			},
		}
		if err := r.DeleteAllServices(ctx); err != nil {
			return err
		}
		for i, item := range items {
			if err := r.CreateService(ctx, item, i); err != nil {
				return err
			}
		}
	}
	return nil
}

// --- Footwear Preset ---
func (r *ContentRepo) applyFootwearPreset(ctx context.Context, seedItems bool) error {
	site := defaultSiteConfig()
	site.Industry = "footwear"
	site.HeroVideoURL = "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4"
	site.HeroBackgroundType = "gradient"
	site.HeroOverlayOpacity = 30
	site.Animations = models.AnimationConfig{
		GlobalEnabled:   true,
		Intensity:       "balanced",
		ScrollReveal:    true,
		DefaultDuration: "0.8s",
		CardHoverEffect: "lift",
	}
	if err := r.UpdateSite(ctx, site); err != nil {
		return err
	}

	sections := []models.Section{
		{
			ID:                "hero",
			Type:              "hero",
			Enabled:           true,
			SortOrder:         0,
			Title:             site.HeroTitle,
			Subtitle:          site.HeroSubtitle,
			Badge:             site.HeroBadge,
			BackgroundType:    "gradient",
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
			PrimaryCtaText:    "Browse Collections",
			PrimaryCtaUrl:     "/collections",
			SecondaryCtaText:  "Shop on WhatsApp",
			SecondaryCtaUrl:   site.WhatsAppURL,
		},
		{
			ID:                "press",
			Type:              "features",
			Enabled:           true,
			SortOrder:         1,
			Title:             "The Shoes <span class='accent'>Professionals Trust</span>",
			Subtitle:          "Curated styles that look sharp from nine to five and beyond",
			Badge:             "~ As Featured In ~",
			BackgroundType:    "default",
			Animation:         "fade-in",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
		},
		{
			ID:                "features",
			Type:              "features",
			Enabled:           true,
			SortOrder:         2,
			Title:             "What Makes Us <span class='accent'>Different</span>",
			Subtitle:          "Every pair chosen for the woman who shows up ready to lead",
			Badge:             "~ Why Luna Collections ~",
			BackgroundType:    "default",
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
			HoverEffect:       "lift",
		},
		{
			ID:                "products",
			Type:              "products",
			Enabled:           true,
			SortOrder:         3,
			Title:             "Find Your <span class='accent'>Perfect Pair</span>",
			Subtitle:          "From classic court pumps to structured oxfords — beautifully made and office-ready",
			Badge:             "~ Best Sellers ~",
			BackgroundType:    "default",
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
			HoverEffect:       "lift",
			PrimaryCtaText:    "View All Collections",
			PrimaryCtaUrl:     "/collections",
		},
		{
			ID:                "gallery",
			Type:              "gallery",
			Enabled:           true,
			SortOrder:         4,
			Title:             "Styles That <span class='accent'>Lead the Room</span>",
			Subtitle:          "A glimpse of our official footwear lookbook",
			Badge:             "~ Lookbook ~",
			BackgroundType:    "default",
			Animation:         "fade-in",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
			PrimaryCtaText:    "View Full Lookbook",
			PrimaryCtaUrl:     "/gallery",
		},
		{
			ID:                "pricing",
			Type:              "pricing",
			Enabled:           true,
			SortOrder:         5,
			Title:             site.PricingHighlight.Title,
			Subtitle:          site.PricingHighlight.Subtitle,
			Badge:             "~ Honest Pricing ~",
			BackgroundType:    "color",
			BackgroundColor:   "#1A1F3B",
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
			PrimaryCtaText:    "View Full Price List",
			PrimaryCtaUrl:     "/pricing",
		},
		{
			ID:                "reviews",
			Type:              "reviews",
			Enabled:           true,
			SortOrder:         6,
			Title:             "Steps That <span class='accent'>Speak for Themselves</span>",
			Subtitle:          "Real reviews from executives and professionals who trust Luna Collections",
			Badge:             "~ Happy Customers ~",
			BackgroundType:    "default",
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
		},
		{
			ID:                "faq",
			Type:              "faq",
			Enabled:           true,
			SortOrder:         7,
			Title:             "Frequently <span class='accent'>Asked Questions</span>",
			Subtitle:          "Sizing, exchanges, nationwide delivery, and leather care",
			Badge:             "~ Got Questions? ~",
			BackgroundType:    "default",
			Animation:         "fade-in",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
		},
		{
			ID:                "inquiry",
			Type:              "inquiry",
			Enabled:           true,
			SortOrder:         8,
			Title:             "Order Your Pair <span class='accent'>Today</span>",
			Subtitle:          "WhatsApp us for sizing help — walk-ins welcome at our Westlands showroom",
			Badge:             "~ Ready to Step Up? ~",
			BackgroundType:    "default",
			Animation:         "slide-up",
			AnimationDuration: "0.8s",
			AnimationDelay:    "0.1s",
		},
	}
	if err := r.DeleteAllSections(ctx); err != nil {
		return err
	}
	for i, s := range sections {
		if err := r.CreateSection(ctx, s, i); err != nil {
			return err
		}
	}

	if seedItems {
		items := []models.LashService{
			{
				ID: "shoe-pump-classic", Name: "Classic Pointed Court Pump in Full-Grain Nappa", Slug: "classic-pointed-court-pump-nappa",
				Emoji: "👠", Style: "Court Pump", Length: "75mm Heel", Category: "Pumps & Heels", Badge: "Boardroom Staple",
				Description: "Supple Italian calfskin nappa upper with memory foam cushioned footbed. Structured for effortless 10-hour days without pinching.",
				PriceKES:    4800, OriginalPriceKES: 5500, Currency: "KES", Featured: true, InStock: true,
				Specs:    map[string]string{"Material": "Full-Grain Italian Calfskin Nappa", "Insole": "Dual-Density Ergonomic Memory Foam", "Heel Height": "75mm Wrapped Stiletto Heel", "Sizes": "EU 36 – 42"},
				Images:   []string{"https://images.unsplash.com/photo-1543163521-1bf539c55dd2?w=800&q=80"},
				View360:  []string{"https://images.unsplash.com/photo-1543163521-1bf539c55dd2?w=800&q=80"},
				VideoURL: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4",
			},
			{
				ID: "shoe-pump-block", Name: "Block Heel Office Pump in Burnished Chestnut", Slug: "block-heel-office-pump-burnished-chestnut",
				Emoji: "👡", Style: "Block Pump", Length: "55mm Block", Category: "Pumps & Heels", Badge: "All-Day Comfort",
				Description: "Architectural block heel providing maximum ground stability across tiled offices. Hand-burnished toe cap with breathable leather lining.",
				PriceKES:    5200, OriginalPriceKES: 5900, Currency: "KES", Featured: true, InStock: true,
				Specs:   map[string]string{"Material": "Hand-Finished Argentine Vegetable Leather", "Sole": "Non-Slip Anti-Skid Rubber Pod Sole", "Heel": "55mm Architectural Block Heel", "Sizes": "EU 35 – 42"},
				Images:  []string{"https://images.unsplash.com/photo-1596702874230-b5706dfb5bc7?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1596702874230-b5706dfb5bc7?w=800&q=80"},
			},
			{
				ID: "shoe-oxford-brogue", Name: "Kensington Full-Grain Leather Wingtip Oxford", Slug: "kensington-full-grain-leather-wingtip-oxford",
				Emoji: "👞", Style: "Brogue", Length: "25mm Heel", Category: "Oxfords & Derbies", Badge: "Executive Line",
				Description: "Heritage perforated wingtip brogue detailing crafted in stiffened box calf leather. Formal, razor-sharp, and unmistakably executive.",
				PriceKES:    6500, OriginalPriceKES: 7500, Currency: "KES", Featured: true, InStock: true,
				Specs:   map[string]string{"Construction": "Goodyear-Welted Stitched Sole", "Leather": "Full-Grain French Box Calfskin", "Lining": "Natural Breathable Vachetta", "Sizes": "EU 36 – 43"},
				Images:  []string{"https://images.unsplash.com/photo-1761726065663-6e550c25a6ef?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1761726065663-6e550c25a6ef?w=800&q=80"},
			},
			{
				ID: "shoe-stiletto-power", Name: "Sovereign Pointed Stiletto 90mm in Onyx Patent", Slug: "sovereign-pointed-stiletto-90mm-onyx-patent",
				Emoji: "✨", Style: "Power Stiletto", Length: "90mm High Heel", Category: "Pumps & Heels", Badge: "Red-Carpet",
				Description: "High-gloss patent leather with a sharp blade heel. Engineered weight distribution prevents forward foot slide during high-profile presentations.",
				PriceKES:    6800, OriginalPriceKES: 7800, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Upper": "Mirror-Finish Italian Patent Leather", "Heel": "90mm Reinforced Steel-Core Stiletto", "Sole": "Crimson Lacquered Dress Leather Outsole", "Fit": "True to European Size (Width: B/Medium)"},
				Images:  []string{"https://images.unsplash.com/photo-1543163521-1bf539c55dd2?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1543163521-1bf539c55dd2?w=800&q=80"},
			},
			{
				ID: "shoe-loafer-penny", Name: "Regent Street Saddle Penny Loafer", Slug: "regent-street-saddle-penny-loafer",
				Emoji: "🥿", Style: "Penny Loafer", Length: "Flat 20mm", Category: "Loafers & Flats", Badge: "Smart Casual",
				Description: "Slip-on ease with precision stitched saddle strap. Beautifully proportioned for ankle-length tailored trousers and executive blazers.",
				PriceKES:    5400, OriginalPriceKES: 6200, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Upper": "Waxed Full-Grain Pull-Up Leather", "Footbed": "Arch-Support Cushioned Leather Insole", "Construction": "Moccasin Hand-Stitched Vamp", "Sizes": "EU 36 – 42"},
				Images:  []string{"https://images.unsplash.com/photo-1533867617858-e7b97e060509?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1533867617858-e7b97e060509?w=800&q=80"},
			},
			{
				ID: "shoe-monk-strap", Name: "Hand-Patinated Double Monk Strap Shoe", Slug: "hand-patinated-double-monk-strap-shoe",
				Emoji: "🎖️", Style: "Monk Strap", Length: "30mm Heel", Category: "Oxfords & Derbies", Badge: "Hand-Burnished",
				Description: "Dual polished brass buckles over an asymmetric saddle. Hand-painted museum patina finish where no two pairs are identical.",
				PriceKES:    7200, OriginalPriceKES: 8200, Currency: "KES", Featured: true, InStock: true,
				Specs:   map[string]string{"Hardware": "Solid Brushed Brass Double Buckles", "Patina": "Artisan Hand-Dyed Tobacco Museum Calf", "Sole": "Double Leather Sole with Brass Toe Taps", "Sizes": "EU 37 – 43"},
				Images:  []string{"https://images.unsplash.com/photo-1614252235316-8c857d38b5f4?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1614252235316-8c857d38b5f4?w=800&q=80"},
			},
			{
				ID: "shoe-slingback", Name: "Mayfair Square-Toe Block Slingback Pump", Slug: "mayfair-square-toe-block-slingback-pump",
				Emoji: "💫", Style: "Slingback", Length: "50mm Block", Category: "Pumps & Heels", Badge: "Breathable Poise",
				Description: "Contrasting black patent cap-toe over cream lambskin leather with an elasticated adjustable slingback heel strap.",
				PriceKES:    5500, OriginalPriceKES: 6200, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Upper": "Dual-Tone Soft Lambskin & Patent Leather Cap", "Strap": "Hidden Elastic Stretch Slingback", "Heel": "50mm Flared Block Heel", "Sizes": "EU 35 – 42"},
				Images:  []string{"https://images.unsplash.com/photo-1596702874230-b5706dfb5bc7?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1596702874230-b5706dfb5bc7?w=800&q=80"},
			},
			{
				ID: "shoe-ankle-chelsea", Name: "Chelsea Burnished Calfskin Leather Ankle Boot", Slug: "chelsea-burnished-calfskin-leather-ankle-boot",
				Emoji: "👢", Style: "Ankle Boot", Length: "35mm Stacked", Category: "Boots", Badge: "Rain & AC Ready",
				Description: "Structured dress boot with twin elastic gore side panels and woven pull tabs. Effortlessly slips on under dress trousers.",
				PriceKES:    7800, OriginalPriceKES: 8900, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Upper": "Water-Resistant Full-Grain Waxed Calfskin", "Panels": "Heavy-Duty Ribbed Elastic Gore", "Sole": "Dainite Studded Rubber Commando Outsole", "Sizes": "EU 36 – 43"},
				Images:  []string{"https://images.unsplash.com/photo-1543163521-1bf539c55dd2?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1543163521-1bf539c55dd2?w=800&q=80"},
			},
			{
				ID: "shoe-kitten-heel", Name: "Conservative Pointed Kitten Heel 45mm", Slug: "conservative-pointed-kitten-heel-45mm",
				Emoji: "🌸", Style: "Kitten Heel", Length: "45mm Low Heel", Category: "Pumps & Heels", Badge: "Diplomatic Fit",
				Description: "A delicate curved low heel designed for conservative dress codes and all-day courtroom or hospital rounds mobility.",
				PriceKES:    4600, OriginalPriceKES: 5200, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Leather": "Soft Grain Cowhide Leather", "Heel": "45mm Contoured Curved Kitten Heel", "Padding": "Orthotic Metatarsal Support Pad", "Sizes": "EU 35 – 41"},
				Images:  []string{"https://images.unsplash.com/photo-1543163521-1bf539c55dd2?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1543163521-1bf539c55dd2?w=800&q=80"},
			},
			{
				ID: "shoe-tassel-loafer", Name: "St. James Goodyear-Welted Tassel Loafer", Slug: "st-james-goodyear-welted-tassel-loafer",
				Emoji: "👞", Style: "Tassel Loafer", Length: "Flat 25mm", Category: "Loafers & Flats", Badge: "Artisan Leather",
				Description: "Braided leather collar laced through brass eyelets with dual hand-cut leather tassels. Polished to a deep mirror shine.",
				PriceKES:    5800, OriginalPriceKES: 6500, Currency: "KES", Featured: false, InStock: true,
				Specs:   map[string]string{"Detailing": "Braided Rawhide Lacing & Double Tassels", "Material": "French Full-Grain Crust Leather", "Construction": "Traditional Goodyear Welt", "Sizes": "EU 36 – 43"},
				Images:  []string{"https://images.unsplash.com/photo-1533867617858-e7b97e060509?w=800&q=80"},
				View360: []string{"https://images.unsplash.com/photo-1533867617858-e7b97e060509?w=800&q=80"},
			},
		}
		if err := r.DeleteAllServices(ctx); err != nil {
			return err
		}
		for i, item := range items {
			if err := r.CreateService(ctx, item, i); err != nil {
				return err
			}
		}
	}
	return nil
}
