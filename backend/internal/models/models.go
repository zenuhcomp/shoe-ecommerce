package models

type AnimationConfig struct {
	GlobalEnabled   bool   `json:"globalEnabled"`
	Intensity       string `json:"intensity"`       // subtle, balanced, energetic, off
	ScrollReveal    bool   `json:"scrollReveal"`
	DefaultDuration string `json:"defaultDuration"` // e.g. "0.8s"
	CardHoverEffect string `json:"cardHoverEffect"` // lift, glow, scale, none
}

type SiteConfig struct {
	Industry           string            `json:"industry"` // furniture, vehicles, spa, beauty, footwear, custom
	BrandName          string            `json:"brandName"`
	BrandScript        string            `json:"brandScript"`
	Tagline            string            `json:"tagline"`
	HeroTitle          string            `json:"heroTitle"`
	HeroSubtitle       string            `json:"heroSubtitle"`
	HeroBadge          string            `json:"heroBadge"`
	LogoURL            string            `json:"logoUrl"`
	HeroVideoURL       string            `json:"heroVideoUrl"`
	HeroBackgroundType string            `json:"heroBackgroundType"` // gradient, video, image, color
	HeroOverlayOpacity int               `json:"heroOverlayOpacity"`
	Announcement       string            `json:"announcement"`
	Location           string            `json:"location"`
	StudioAddress      string            `json:"studioAddress"`
	Phone              string            `json:"phone"`
	Email              string            `json:"email"`
	WhatsAppURL        string            `json:"whatsAppUrl"`
	InstagramURL       string            `json:"instagramUrl"`
	TikTokURL          string            `json:"tiktokUrl"`
	PriceRange         string            `json:"priceRange"`
	Copyright          string            `json:"copyright"`
	MetaDescription    string            `json:"metaDescription"`
	Theme              Theme             `json:"theme"`
	Animations         AnimationConfig   `json:"animations"`
	About              AboutSection      `json:"about"`
	HouseCalls         HouseCallsSection `json:"houseCalls"`
	Aftercare          AftercareSection  `json:"aftercare"`
	PricingHighlight   PricingHighlight  `json:"pricingHighlight"`
	Payment            PaymentConfig     `json:"payment"`
}

type PaymentConfig struct {
	// Lipa Na M-Pesa Online (Daraja STK Push)
	MpesaExpressEnabled bool   `json:"mpesaExpressEnabled"`
	ConsumerKey         string `json:"consumerKey"`
	ConsumerSecret      string `json:"consumerSecret"`
	PassKey             string `json:"passKey"`
	ShortCode           string `json:"shortCode"` // Business ShortCode (e.g. 174379 sandbox or live paybill)
	PartyB              string `json:"partyB"`
	AccountReference    string `json:"accountReference"`
	Environment         string `json:"environment"` // "sandbox" or "production"
	CallbackURL         string `json:"callbackUrl"`

	// Lipa Na M-Pesa Offline (Till & Paybill Manual Verification)
	MpesaManualEnabled bool   `json:"mpesaManualEnabled"`
	TillNumber         string `json:"tillNumber"`
	PaybillNumber      string `json:"paybillNumber"`
	PaybillAccount     string `json:"paybillAccount"`
	MpesaShortCode     string `json:"mpesaShortCode"`
	MpesaInstructions  string `json:"mpesaInstructions"`

	// Credit & Debit Card (Stripe / IntaSend / Direct Gateway)
	CardEnabled      bool   `json:"cardEnabled"`
	CardGateway      string `json:"cardGateway"` // "stripe", "intasend", "custom"
	CardPublicKey    string `json:"cardPublicKey"`
	CardSecretKey    string `json:"cardSecretKey"`
	CardInstructions string `json:"cardInstructions"`

	// Cash on Delivery / In-Store Pickup
	CashEnabled      bool   `json:"cashEnabled"`
	CashInstructions string `json:"cashInstructions"`

	// Bank Wire / RTGS / KCB Buni / Commercial Banks
	BankEnabled      bool   `json:"bankEnabled"`
	BankName         string `json:"bankName"`
	BankAccountName  string `json:"bankAccountName"`
	BankAccountNo    string `json:"bankAccountNo"`
	BankBranch       string `json:"bankBranch"`
	BankSwiftCode    string `json:"bankSwiftCode"`
	BankInstructions string `json:"bankInstructions"`
}

type PaymentMethodInfo struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Type         string         `json:"type"` // "online", "manual", "card", "cash", "bank"
	Enabled      bool           `json:"enabled"`
	Description  string         `json:"description"`
	PublicParams map[string]any `json:"publicParams,omitempty"`
}

type PaymentInitiateRequest struct {
	OrderID       string            `json:"orderId"`
	ServiceSlug   string            `json:"serviceSlug"`
	CustomerName  string            `json:"customerName"`
	CustomerPhone string            `json:"customerPhone"`
	CustomerEmail string            `json:"customerEmail"`
	AmountKES     int               `json:"amountKes"`
	Method        string            `json:"method"`
	MpesaPhone    string            `json:"mpesaPhone,omitempty"`
	Notes         string            `json:"notes,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

type PaymentInitiateResponse struct {
	Success        bool              `json:"success"`
	Status         string            `json:"status"` // "pending", "completed", "failed"
	Message        string            `json:"message"`
	TransactionID  string            `json:"transactionId,omitempty"`
	PromptRequired bool              `json:"promptRequired"` // e.g. true for STK push
	Details        map[string]string `json:"details,omitempty"`
	WhatsAppURL    string            `json:"whatsAppUrl,omitempty"`
}

type PaymentVerifyRequest struct {
	OrderID   string `json:"orderId"`
	Method    string `json:"method"`
	Reference string `json:"reference"`
}

type Theme struct {
	Primary      string `json:"primary"`
	PrimaryDark  string `json:"primaryDark"`
	Accent       string `json:"accent"`
	AccentLight  string `json:"accentLight"`
	Background   string `json:"background"`
	Surface      string `json:"surface"`
	Text         string `json:"text"`
	TextMuted    string `json:"textMuted"`
	HeroGradient string `json:"heroGradient"`
	FontHeading  string `json:"fontHeading"`
	FontBody     string `json:"fontBody"`
	FontScript   string `json:"fontScript"`
}

type AboutSection struct {
	Title       string   `json:"title"`
	Intro       string   `json:"intro"`
	Story       string   `json:"story"`
	Highlights  []string `json:"highlights"`
	ArtistName  string   `json:"artistName"`
	ArtistTitle string   `json:"artistTitle"`
}

type HouseCallsSection struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Notes       []string `json:"notes"`
}

type AftercareSection struct {
	Title string         `json:"title"`
	Intro string         `json:"intro"`
	Tips  []AftercareTip `json:"tips"`
}

type AftercareTip struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type PricingHighlight struct {
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	PriceFrom  int    `json:"priceFrom"`
	PriceTo    int    `json:"priceTo"`
	Currency   string `json:"currency"`
	Badge      string `json:"badge"`
	Includes   string `json:"includes"`
	StudioNote string `json:"studioNote"`
}

type LashService struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Slug             string            `json:"slug"`
	Emoji            string            `json:"emoji"`
	Style            string            `json:"style"`
	Length           string            `json:"length"`
	Category         string            `json:"category"`
	Badge            string            `json:"badge"`
	Description      string            `json:"description"`
	PriceKES         int               `json:"priceKes"`
	OriginalPriceKES int               `json:"originalPriceKes"`
	Currency         string            `json:"currency"`
	Specs            map[string]string `json:"specs"`
	Featured         bool              `json:"featured"`
	InStock          bool              `json:"inStock"`
	Images           []string          `json:"images"`
	View360          []string          `json:"view360"`
	VideoURL         string            `json:"videoUrl"`
}

type Section struct {
	ID                string `json:"id"`
	Type              string `json:"type"` // hero, features, products, video_banner, media_360, gallery, pricing, reviews, faq, about, inquiry, cta, custom_block
	Enabled           bool   `json:"enabled"`
	SortOrder         int    `json:"sortOrder"`
	Title             string `json:"title"`
	Subtitle          string `json:"subtitle"`
	Badge             string `json:"badge"`
	Content           string `json:"content"`
	BackgroundType    string `json:"backgroundType"` // default, color, gradient, image, video
	BackgroundColor   string `json:"backgroundColor"`
	BackgroundImage   string `json:"backgroundImage"`
	BackgroundVideo   string `json:"backgroundVideo"`
	OverlayOpacity    int    `json:"overlayOpacity"` // 0 - 100
	Animation         string `json:"animation"`      // none, fade-in, slide-up, slide-left, slide-right, zoom-in, float, pulse
	AnimationDuration string `json:"animationDuration"`
	AnimationDelay    string `json:"animationDelay"`
	HoverEffect       string `json:"hoverEffect"` // lift, scale, glow, none
	PrimaryCtaText    string `json:"primaryCtaText"`
	PrimaryCtaUrl     string `json:"primaryCtaUrl"`
	SecondaryCtaText  string `json:"secondaryCtaText"`
	SecondaryCtaUrl   string `json:"secondaryCtaUrl"`
	ExtraData         string `json:"extraData"` // JSON string
}

type MediaAsset struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Category     string `json:"category"`  // furniture, vehicles, spa, beauty, footwear, general
	MediaType    string `json:"mediaType"` // image, video
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
	AltText      string `json:"altText"`
	Tags         string `json:"tags"`
}

type PresetSummary struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Industry     string   `json:"industry"`
	Tagline      string   `json:"tagline"`
	Description  string   `json:"description"`
	Icon         string   `json:"icon"`
	PrimaryColor string   `json:"primaryColor"`
	AccentColor  string   `json:"accentColor"`
	SampleImages []string `json:"sampleImages"`
}

type ApplyPresetRequest struct {
	PresetID      string `json:"presetId"`
	PresetIDSnake string `json:"preset_id"`
	ID            string `json:"id"`
	SeedItems     bool   `json:"seedItems"`
}

type Feature struct {
	ID          string `json:"id,omitempty"`
	Icon        string `json:"icon"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type PressMention struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
	Note string `json:"note"`
}

type Review struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Handle    string `json:"handle"`
	Initials  string `json:"initials"`
	Rating    int    `json:"rating"`
	Quote     string `json:"quote"`
	Published bool   `json:"published,omitempty"`
}

type FAQ struct {
	ID       int    `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type GalleryItem struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	ImageURL string `json:"imageUrl"`
	Alt      string `json:"alt"`
}

type BookingInquiry struct {
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	ServiceSlug string `json:"serviceSlug"`
	PreferredAt string `json:"preferredAt"`
	Message     string `json:"message"`
}

type BookingResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type OrderRequest struct {
	Name             string `json:"name"`
	CustomerName     string `json:"customerName,omitempty"`
	Phone            string `json:"phone"`
	CustomerPhone    string `json:"customerPhone,omitempty"`
	Email            string `json:"email"`
	CustomerEmail    string `json:"customerEmail,omitempty"`
	ServiceSlug      string `json:"serviceSlug"`
	Size             string `json:"size"`
	Quantity         int    `json:"quantity"`
	PaymentMethod    string `json:"paymentMethod"`
	MpesaPhone       string `json:"mpesaPhone"`
	OfflineChannel   string `json:"offlineChannel"`
	PaymentReference string `json:"paymentReference"`
	Notes            string `json:"notes"`
}

type OrderResponse struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	Message     string `json:"message"`
	WhatsAppURL string `json:"whatsAppUrl,omitempty"`
}

type NavLink struct {
	ID      string `json:"id,omitempty"`
	Label   string `json:"label"`
	Path    string `json:"path"`
	Visible bool   `json:"visible,omitempty"`
}

type HomePage struct {
	Site     SiteConfig     `json:"site"`
	Services []LashService  `json:"services"`
	Features []Feature      `json:"features"`
	Press    []PressMention `json:"press"`
	Reviews  []Review       `json:"reviews"`
	FAQ      []FAQ          `json:"faq"`
	Nav      []NavLink      `json:"nav"`
	Sections []Section      `json:"sections"`
}
