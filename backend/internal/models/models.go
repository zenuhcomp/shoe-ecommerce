package models

type SiteConfig struct {
	BrandName        string            `json:"brandName"`
	BrandScript      string            `json:"brandScript"`
	Tagline          string            `json:"tagline"`
	HeroTitle        string            `json:"heroTitle"`
	HeroSubtitle     string            `json:"heroSubtitle"`
	HeroBadge        string            `json:"heroBadge"`
	Announcement     string            `json:"announcement"`
	Location         string            `json:"location"`
	StudioAddress    string            `json:"studioAddress"`
	Phone            string            `json:"phone"`
	Email            string            `json:"email"`
	WhatsAppURL      string            `json:"whatsAppUrl"`
	InstagramURL     string            `json:"instagramUrl"`
	TikTokURL        string            `json:"tiktokUrl"`
	PriceRange       string            `json:"priceRange"`
	Copyright        string            `json:"copyright"`
	MetaDescription  string            `json:"metaDescription"`
	Theme            Theme             `json:"theme"`
	About            AboutSection      `json:"about"`
	HouseCalls       HouseCallsSection `json:"houseCalls"`
	Aftercare        AftercareSection  `json:"aftercare"`
	PricingHighlight PricingHighlight  `json:"pricingHighlight"`
	Payment          PaymentConfig     `json:"payment"`
}

type PaymentConfig struct {
	TillNumber     string `json:"tillNumber"`
	PaybillNumber  string `json:"paybillNumber"`
	PaybillAccount string `json:"paybillAccount"`
	MpesaShortCode string `json:"mpesaShortCode"`
}

type Theme struct {
	Primary       string `json:"primary"`
	PrimaryDark   string `json:"primaryDark"`
	Accent        string `json:"accent"`
	AccentLight   string `json:"accentLight"`
	Background    string `json:"background"`
	Surface       string `json:"surface"`
	Text          string `json:"text"`
	TextMuted     string `json:"textMuted"`
	HeroGradient  string `json:"heroGradient"`
	FontHeading   string `json:"fontHeading"`
	FontBody      string `json:"fontBody"`
	FontScript    string `json:"fontScript"`
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
	Title string       `json:"title"`
	Intro string       `json:"intro"`
	Tips  []AftercareTip `json:"tips"`
}

type AftercareTip struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type PricingHighlight struct {
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	PriceFrom   int    `json:"priceFrom"`
	PriceTo     int    `json:"priceTo"`
	Currency    string `json:"currency"`
	Badge       string `json:"badge"`
	Includes    string `json:"includes"`
	StudioNote  string `json:"studioNote"`
}

type LashService struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Emoji       string   `json:"emoji"`
	Style       string   `json:"style"`
	Length      string   `json:"length"`
	Description string   `json:"description"`
	PriceKES    int      `json:"priceKes"`
	Featured    bool     `json:"featured"`
	Images      []string `json:"images"`
	View360     []string `json:"view360"`
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
	Name               string `json:"name"`
	Phone              string `json:"phone"`
	Email              string `json:"email"`
	ServiceSlug        string `json:"serviceSlug"`
	Size               string `json:"size"`
	Quantity           int    `json:"quantity"`
	PaymentMethod      string `json:"paymentMethod"`
	MpesaPhone         string `json:"mpesaPhone"`
	OfflineChannel     string `json:"offlineChannel"`
	PaymentReference   string `json:"paymentReference"`
	Notes              string `json:"notes"`
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
}
