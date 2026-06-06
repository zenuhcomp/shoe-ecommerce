package repository

import "github.com/kibsoft/luna-collections/backend/internal/models"

const (
	imageClassicPump = "https://images.unsplash.com/photo-1543163521-1bf539c55dd2?w=800&q=80"
	imageCourtHeel   = "https://images.unsplash.com/photo-1596702874230-b5706dfb5bc7?w=800&q=80"
	imageGlitterPump = "https://images.unsplash.com/photo-1581101767113-1677fc2beaa8?w=800&q=80"
	imageWornHeel    = "https://images.unsplash.com/photo-1651126581286-a6ae76cb5c78?w=800&q=80"
	imageMaryJane    = "https://images.unsplash.com/photo-1757191428048-a1747039f267?w=800&q=80"
	imagePeepToe     = "https://images.unsplash.com/photo-1596702874230-b5706dfb5bc7?w=800&q=80"
	imageLoafer      = "https://images.unsplash.com/photo-1761726065663-6e550c25a6ef?w=800&q=80"
)

func defaultSiteConfig() models.SiteConfig {
	return models.SiteConfig{
		BrandName:       "Luna",
		BrandScript:     "Collections",
		Tagline:         "Step Into Elegance · Own Every Room",
		HeroTitle:       "Ladies Official Shoes, Curated for the Modern Professional",
		HeroSubtitle:    "Premium office-ready footwear — pumps, oxfords, loafers & more, delivered across Kenya",
		HeroBadge:       "◆ State-of-the-Art Shoe Boutique — Nairobi ◆",
		Announcement:    "✦ New arrivals weekly · Sizes 35–42 · Free sizing consult · Nationwide delivery ✦",
		Location:        "Nairobi, Kenya",
		StudioAddress:   "Westlands, Nairobi",
		Phone:           "+254 712 888 420",
		Email:           "hello@lunacollections.co.ke",
		WhatsAppURL:     "https://wa.me/254712888420",
		InstagramURL:    "https://instagram.com/lunacollections.ke",
		TikTokURL:       "https://tiktok.com/@lunacollections.ke",
		PriceRange:      "KES 3,500 – 12,500",
		Copyright:       "© 2026 Luna Collections. Crafted with care in Nairobi, Kenya.",
		MetaDescription: "Luna Collections — Premium ladies official shoes in Nairobi. Pumps, oxfords, loafers, block heels & office footwear with nationwide delivery.",
		Theme: models.Theme{
			Primary: "#1A1F3B", PrimaryDark: "#0F1224", Accent: "#C4A574", AccentLight: "#E8DCC8",
			Background: "#FAF8F5", Surface: "#FFFFFF", Text: "#141824", TextMuted: "#5A6178",
			HeroGradient: "linear-gradient(145deg, #0F1224 0%, #1A1F3B 42%, #2D3561 100%)",
			FontHeading:  "'Playfair Display', Georgia, serif", FontBody: "'DM Sans', system-ui, sans-serif",
			FontScript: "'Cormorant Garamond', Georgia, serif",
		},
		About: models.AboutSection{
			Title: "Our Story", Intro: "Nairobi's destination for polished, professional footwear.",
			Story: "Luna Collections was founded for women who lead — in boardrooms, courtrooms, classrooms, and everywhere in between. Every pair is selected for comfort, durability, and timeless office elegance. From classic pumps to structured oxfords, we stock full size runs and offer expert fitting so you never compromise on the first step of your day.",
			Highlights: []string{
				"Curated official & office-ready styles from trusted makers",
				"Full size range 35–42 with half-size requests on WhatsApp",
				"Premium leather, suede, and vegan options available",
				"Showroom in Westlands plus nationwide delivery",
			},
			ArtistName: "Luna Collections", ArtistTitle: "Premium Ladies Official Footwear",
		},
		HouseCalls: models.HouseCallsSection{
			Title:       "Delivered to Your Door",
			Description: "Order from anywhere in Kenya. We pack every pair with care — tissue, box, and sizing card included — and dispatch within 24 hours on in-stock items.",
			Notes: []string{
				"Nairobi CBD & Westlands: same-day delivery on orders before 2 PM",
				"Nationwide via trusted courier — 1–3 business days",
				"Free delivery on orders above KES 8,000 within Nairobi",
				"Easy exchanges within 7 days for unworn pairs in original packaging",
			},
		},
		Aftercare: models.AftercareSection{
			Title: "Shoe Care & Sizing Guide",
			Intro: "Protect your investment. Every purchase includes care tips and our fit guarantee — message us if the size isn't perfect.",
			Tips: []models.AftercareTip{
				{Title: "Find your fit", Description: "Measure feet in the evening when slightly swollen. Leave ½ cm at the toe for closed pumps; oxfords should feel snug, not pinching."},
				{Title: "Break in gently", Description: "Wear new leather shoes indoors for 30 minutes the first day. Use a shoe horn and alternate pairs during the first week."},
				{Title: "Clean & condition", Description: "Wipe leather after each wear. Condition monthly; use suede brush for nubuck. Store with shoe trees to maintain shape."},
				{Title: "Heel & sole care", Description: "Replace heel caps before they wear through. Resole quality leather pairs — they'll last seasons longer."},
			},
		},
		PricingHighlight: models.PricingHighlight{
			Title: "Transparent Pricing", Subtitle: "Office-ready quality at honest prices — no hidden fees at checkout.",
			PriceFrom: 3500, PriceTo: 12500, Currency: "KES", Badge: "Bestseller range",
			Includes:   "Quality assurance, care card, branded packaging, and 7-day exchange policy.",
			StudioNote: "◆ Westlands showroom · ◆ Sizes 35–42 · ◆ Nationwide delivery",
		},
		Payment: models.PaymentConfig{
			TillNumber: "888420", PaybillNumber: "247247", PaybillAccount: "Luna Collections",
			MpesaShortCode: "888420",
		},
	}
}

func shoeMedia(main string, angles ...string) ([]string, []string) {
	images := []string{main}
	for _, a := range angles {
		images = append(images, a)
	}
	view360 := append([]string{}, images...)
	if len(view360) < 4 {
		view360 = append(view360, main)
	}
	return images, view360
}

func defaultServices() []models.LashService {
	type seed struct {
		models.LashService
		main string
		more []string
	}
	raw := []seed{
		{models.LashService{ID: "shoe-pump-classic", Name: "Classic Court Pump", Slug: "classic-court-pump", Emoji: "👠", Style: "Official", Length: "Mid Heel", Description: "Timeless pointed-toe pump — boardroom essential in patent or matte leather", PriceKES: 4500, Featured: true},
			imageClassicPump,
			[]string{imageCourtHeel, imageGlitterPump, imageWornHeel}},
		{models.LashService{ID: "shoe-pump-block", Name: "Block Heel Pump", Slug: "block-heel-pump", Emoji: "👡", Style: "Comfort", Length: "Low–Mid", Description: "All-day stability with a refined silhouette for long office days", PriceKES: 5200},
			imageCourtHeel,
			[]string{imageClassicPump, imageGlitterPump}},
		{models.LashService{ID: "shoe-stiletto", Name: "Pointed Stiletto", Slug: "pointed-stiletto", Emoji: "✨", Style: "Statement", Length: "High Heel", Description: "Sharp toe and slim heel for presentations, events, and power dressing", PriceKES: 6800},
			imageWornHeel,
			[]string{imageClassicPump, imageCourtHeel}},
		{models.LashService{ID: "shoe-oxford", Name: "Ladies Oxford", Slug: "ladies-oxford", Emoji: "👞", Style: "Official", Length: "Flat", Description: "Perforated brogue detailing — structured, polished, and unmistakably professional", PriceKES: 5900, Featured: true},
			imageLoafer,
			[]string{imageMaryJane, imageClassicPump}},
		{models.LashService{ID: "shoe-loafer", Name: "Penny Loafer", Slug: "penny-loafer", Emoji: "🥿", Style: "Smart Casual", Length: "Flat", Description: "Slip-on ease with tailored lines — perfect with trousers and blazers", PriceKES: 4800},
			imageLoafer,
			[]string{imageMaryJane, imageClassicPump}},
		{models.LashService{ID: "shoe-slingback", Name: "Slingback Heel", Slug: "slingback-heel", Emoji: "💫", Style: "Official", Length: "Mid Heel", Description: "Open back, secure strap — breathable elegance for warm Nairobi days", PriceKES: 5500},
			imageWornHeel,
			[]string{imageCourtHeel, imageClassicPump}},
		{models.LashService{ID: "shoe-kitten", Name: "Kitten Heel", Slug: "kitten-heel", Emoji: "🌸", Style: "Classic", Length: "Low Heel", Description: "Subtle lift with maximum poise — ideal for conservative dress codes", PriceKES: 4200},
			imageGlitterPump,
			[]string{imageCourtHeel, imageWornHeel}},
		{models.LashService{ID: "shoe-mule", Name: "Heeled Mule", Slug: "heeled-mule", Emoji: "👢", Style: "Modern", Length: "Mid Heel", Description: "Backless sophistication that transitions from desk to dinner", PriceKES: 5100},
			imageClassicPump,
			[]string{imageCourtHeel, imageWornHeel}},
		{models.LashService{ID: "shoe-ankle-boot", Name: "Ankle Boot", Slug: "ankle-boot", Emoji: "🖤", Style: "Official", Length: "Low Heel", Description: "Structured Chelsea-inspired boot for rainy season and AC-heavy offices", PriceKES: 7500},
			imageLoafer,
			[]string{imageMaryJane, imageClassicPump}},
		{models.LashService{ID: "shoe-mary-jane", Name: "Mary Jane", Slug: "mary-jane", Emoji: "🎀", Style: "Classic", Length: "Low–Mid", Description: "Strap-front charm with a professional finish — vintage meets modern office", PriceKES: 4600},
			imageMaryJane,
			[]string{imageClassicPump, imageLoafer}},
		{models.LashService{ID: "shoe-peep-toe", Name: "Peep-Toe Pump", Slug: "peep-toe-pump", Emoji: "🌺", Style: "Seasonal", Length: "Mid Heel", Description: "Open toe detail for Friday dress-down and client-facing Fridays", PriceKES: 5400},
			imagePeepToe,
			[]string{imageGlitterPump, imageWornHeel}},
		{models.LashService{ID: "shoe-platform", Name: "Platform Loafer", Slug: "platform-loafer", Emoji: "👑", Style: "Trend", Length: "Platform", Description: "Elevated sole, grounded comfort — height without the wobble", PriceKES: 6200},
			imageLoafer,
			[]string{imageMaryJane, imageClassicPump}},
	}
	out := make([]models.LashService, 0, len(raw))
	for _, item := range raw {
		imgs, v360 := shoeMedia(item.main, item.more...)
		svc := item.LashService
		svc.Images = imgs
		svc.View360 = v360
		out = append(out, svc)
	}
	return out
}

func defaultFeatures() []models.Feature {
	return []models.Feature{
		{ID: "feat-1", Icon: "👠", Title: "Curated Official Styles", Description: "Every style vetted for office dress codes — from corporate to creative professional."},
		{ID: "feat-2", Icon: "✦", Title: "Expert Sizing", Description: "WhatsApp fit consults and in-showroom measuring so you order once, correctly."},
		{ID: "feat-3", Icon: "📦", Title: "Nationwide Delivery", Description: "Secure packaging and fast dispatch to Nairobi, Mombasa, Kisumu, and beyond."},
		{ID: "feat-4", Icon: "🔄", Title: "7-Day Exchange", Description: "Unworn pairs in original box? We'll swap sizes or styles hassle-free."},
	}
}

func defaultPress() []models.PressMention {
	return []models.PressMention{
		{ID: "press-1", Name: "VOGUE", Note: "Workwear Edit"},
		{ID: "press-2", Name: "ELLE", Note: "Shoe Guide 2026"},
		{ID: "press-3", Name: "HARPER'S", Note: "Office Style"},
		{ID: "press-4", Name: "GRAZIA", Note: "Nairobi Picks"},
		{ID: "press-5", Name: "FORBES", Note: "Women Who Lead"},
	}
}

func defaultReviews() []models.Review {
	return []models.Review{
		{ID: "rev-1", Name: "Wanjiku M.", Handle: "@wanjiku_m", Initials: "WM", Rating: 5, Quote: "Finally pumps that survive a full day in the office without killing my feet. The block heels are my new uniform."},
		{ID: "rev-2", Name: "Sarah O.", Handle: "@saraho_law", Initials: "SO", Rating: 5, Quote: "Ordered oxfords for court days — impeccable leather and true to size. Delivery to Kisumu took two days."},
		{ID: "rev-3", Name: "Amina K.", Handle: "@aminak_style", Initials: "AK", Rating: 5, Quote: "Luna's sizing consult on WhatsApp saved me. The slingbacks fit perfectly and look expensive — because they are."},
		{ID: "rev-4", Name: "Grace N.", Handle: "@gracen_ke", Initials: "GN", Rating: 5, Quote: "Best ladies official shoe shop in Nairobi. Loafers, pumps, boots — my entire work wardrobe starts here."},
	}
}

func defaultFAQ() []models.FAQ {
	return []models.FAQ{
		{ID: 1, Question: "What sizes do you stock?", Answer: "We carry EU sizes 35–42. Half sizes and wide-fit requests are available on select styles — message us on WhatsApp before ordering."},
		{ID: 2, Question: "How do I know my size?", Answer: "Send a foot tracing or your usual size in UK/EU/US. We offer free sizing guidance and recommend measuring in the evening for the best fit."},
		{ID: 3, Question: "Do you deliver outside Nairobi?", Answer: "Yes — nationwide via trusted couriers. Nairobi orders before 2 PM often arrive same day; other counties typically 1–3 business days."},
		{ID: 4, Question: "Can I exchange if the size is wrong?", Answer: "Unworn pairs in original packaging can be exchanged within 7 days. Contact us on WhatsApp to arrange a swap."},
		{ID: 5, Question: "Are your shoes genuine leather?", Answer: "Most official styles use full-grain or patent leather. Vegan leather options are labeled clearly on each product — ask us for details anytime."},
	}
}

func defaultGallery() []models.GalleryItem {
	return []models.GalleryItem{
		{ID: "g1", Title: "Classic court pump", Category: "pumps", ImageURL: imageClassicPump, Alt: "Black classic court pump on display"},
		{ID: "g2", Title: "Ladies oxford brogue", Category: "oxfords", ImageURL: imageLoafer, Alt: "Black official loafer shoes"},
		{ID: "g3", Title: "Block heel office pump", Category: "pumps", ImageURL: imageCourtHeel, Alt: "Black office pump heel"},
		{ID: "g4", Title: "Penny loafer", Category: "loafers", ImageURL: imageLoafer, Alt: "Black penny loafer flat"},
		{ID: "g5", Title: "Slingback heel", Category: "heels", ImageURL: imageWornHeel, Alt: "Elegant pointed heel"},
		{ID: "g6", Title: "Ankle boot official", Category: "boots", ImageURL: imageLoafer, Alt: "Structured black official shoe"},
		{ID: "g7", Title: "Pointed stiletto", Category: "heels", ImageURL: imageWornHeel, Alt: "Pointed stiletto heel"},
		{ID: "g8", Title: "Kitten heel pump", Category: "pumps", ImageURL: imageGlitterPump, Alt: "Polished low heel pump"},
		{ID: "g9", Title: "Mary Jane official", Category: "classic", ImageURL: imageMaryJane, Alt: "Black Mary Jane shoe"},
	}
}

func defaultNav() []models.NavLink {
	return []models.NavLink{
		{ID: "nav-1", Label: "Collections", Path: "/collections", Visible: true},
		{ID: "nav-2", Label: "Pricing", Path: "/pricing", Visible: true},
		{ID: "nav-3", Label: "Gallery", Path: "/gallery", Visible: true},
		{ID: "nav-4", Label: "About", Path: "/about", Visible: true},
		{ID: "nav-5", Label: "Delivery", Path: "/delivery", Visible: true},
		{ID: "nav-6", Label: "Reviews", Path: "/reviews", Visible: true},
		{ID: "nav-7", Label: "Contact", Path: "/contact", Visible: true},
	}
}
