export interface Theme {
  primary: string;
  primaryDark: string;
  accent: string;
  accentLight: string;
  background: string;
  surface: string;
  text: string;
  textMuted: string;
  heroGradient: string;
  fontHeading: string;
  fontBody: string;
  fontScript: string;
}

export interface AnimationConfig {
  globalEnabled: boolean;
  intensity: 'subtle' | 'balanced' | 'energetic' | 'off';
  scrollReveal: boolean;
  defaultDuration: string;
  cardHoverEffect: 'lift' | 'glow' | 'scale' | 'none';
}

export interface SiteConfig {
  industry?: string; // furniture, vehicles, spa, beauty, footwear, custom
  brandName: string;
  brandScript: string;
  tagline: string;
  heroTitle: string;
  heroSubtitle: string;
  heroBadge: string;
  logoUrl?: string;
  heroVideoUrl?: string;
  heroBackgroundType?: 'gradient' | 'video' | 'image' | 'color';
  heroOverlayOpacity?: number;
  announcement: string;
  location: string;
  studioAddress: string;
  phone: string;
  email: string;
  whatsAppUrl: string;
  instagramUrl: string;
  tiktokUrl: string;
  priceRange: string;
  copyright: string;
  metaDescription: string;
  theme: Theme;
  animations?: AnimationConfig;
  about: AboutSection;
  houseCalls: HouseCallsSection;
  aftercare: AftercareSection;
  pricingHighlight: PricingHighlight;
  payment?: PaymentConfig;
}

export interface PaymentConfig {
  // M-Pesa Online (Daraja STK Push)
  mpesaExpressEnabled: boolean;
  consumerKey: string;
  consumerSecret: string;
  passKey: string;
  shortCode: string;
  partyB: string;
  accountReference: string;
  environment: string; // 'sandbox' | 'production'
  callbackUrl: string;
  // M-Pesa Offline
  mpesaManualEnabled: boolean;
  tillNumber: string;
  paybillNumber: string;
  paybillAccount: string;
  mpesaShortCode: string;
  mpesaInstructions: string;
  // Credit / Debit Card
  cardEnabled: boolean;
  cardGateway: string;
  cardPublicKey: string;
  cardSecretKey: string;
  cardInstructions: string;
  // Cash
  cashEnabled: boolean;
  cashInstructions: string;
  // Bank Transfer
  bankEnabled: boolean;
  bankName: string;
  bankAccountName: string;
  bankAccountNo: string;
  bankBranch: string;
  bankSwiftCode: string;
  bankInstructions: string;
}

export interface PaymentMethodInfo {
  id: string;
  name: string;
  type: string; // 'online' | 'manual' | 'card' | 'cash' | 'bank'
  enabled: boolean;
  description: string;
  publicParams?: Record<string, string>;
}

export interface PaymentInitiateRequest {
  orderId: string;
  serviceSlug: string;
  customerName: string;
  customerPhone: string;
  customerEmail?: string;
  amountKes: number;
  method: string;
  mpesaPhone?: string;
  notes?: string;
  metadata?: Record<string, string>;
}

export interface PaymentInitiateResponse {
  success: boolean;
  status: string; // 'pending' | 'completed' | 'failed'
  message: string;
  transactionId?: string;
  promptRequired: boolean;
  details?: Record<string, string>;
  whatsAppUrl?: string;
}

export interface PaymentVerifyRequest {
  orderId: string;
  method: string;
  reference: string;
}

export interface AboutSection {
  title: string;
  intro: string;
  story: string;
  highlights: string[];
  artistName: string;
  artistTitle: string;
}

export interface HouseCallsSection {
  title: string;
  description: string;
  notes: string[];
}

export interface AftercareSection {
  title: string;
  intro: string;
  tips: { title: string; description: string }[];
}

export interface PricingHighlight {
  title: string;
  subtitle: string;
  priceFrom: number;
  priceTo: number;
  currency: string;
  badge: string;
  includes: string;
  studioNote: string;
}

export interface LashService {
  id: string;
  name: string;
  slug: string;
  emoji: string;
  style: string;
  length: string;
  category?: string;
  badge?: string;
  description: string;
  priceKes: number;
  originalPriceKes?: number;
  currency?: string;
  specs?: Record<string, string>;
  featured?: boolean;
  inStock?: boolean;
  images?: string[];
  view360?: string[];
  videoUrl?: string;
}

export interface Section {
  id: string;
  type: string; // hero, features, products, video_banner, media_360, gallery, pricing, reviews, faq, about, inquiry, cta, custom_block
  enabled: boolean;
  sortOrder: number;
  title: string;
  subtitle: string;
  badge: string;
  content: string;
  backgroundType: 'default' | 'color' | 'gradient' | 'image' | 'video';
  backgroundColor: string;
  backgroundImage: string;
  backgroundVideo: string;
  overlayOpacity: number; // 0 to 100
  animation: string; // none, fade-in, slide-up, slide-left, slide-right, zoom-in, float, pulse
  animationDuration: string;
  animationDelay: string;
  hoverEffect: 'lift' | 'scale' | 'glow' | 'none';
  primaryCtaText: string;
  primaryCtaUrl: string;
  secondaryCtaText: string;
  secondaryCtaUrl: string;
  extraData?: string;
}

export interface MediaAsset {
  id: string;
  name: string;
  category: string;
  mediaType: 'image' | 'video';
  url: string;
  thumbnailUrl?: string;
  altText?: string;
  tags?: string;
}

export interface IndustryPreset {
  id: string;
  name: string;
  industry: string;
  tagline: string;
  description: string;
  icon: string;
  primaryColor: string;
  accentColor: string;
  sampleImages: string[];
}

export interface Feature {
  id?: string;
  icon: string;
  title: string;
  description: string;
}

export interface PressMention {
  id?: string;
  name: string;
  note: string;
}

export interface Review {
  id: string;
  name: string;
  handle: string;
  initials: string;
  rating: number;
  quote: string;
  published?: boolean;
}

export interface FAQ {
  id: number;
  question: string;
  answer: string;
}

export interface GalleryItem {
  id: string;
  title: string;
  category: string;
  imageUrl: string;
  alt: string;
}

export interface NavLink {
  id?: string;
  label: string;
  path: string;
  visible?: boolean;
}

export interface HomePage {
  site: SiteConfig;
  services: LashService[];
  features: Feature[];
  press: PressMention[];
  reviews: Review[];
  faq: FAQ[];
  nav: NavLink[];
  sections?: Section[];
}

export interface BookingInquiry {
  name: string;
  phone: string;
  email?: string;
  serviceSlug?: string;
  preferredAt?: string;
  message?: string;
}

export interface BookingResponse {
  id: string;
  message: string;
}

export type PaymentMethod = 'mpesa_stk' | 'offline_till' | 'offline_paybill' | 'whatsapp';

export interface OrderRequest {
  name: string;
  phone: string;
  email?: string;
  serviceSlug: string;
  size?: string;
  quantity: number;
  paymentMethod: PaymentMethod;
  mpesaPhone?: string;
  offlineChannel?: 'till' | 'paybill';
  paymentReference?: string;
  notes?: string;
}

export interface OrderResponse {
  id: string;
  status: string;
  message: string;
  whatsAppUrl?: string;
}

export interface OrderRecord {
  id: string;
  name: string;
  phone: string;
  email?: string;
  serviceSlug: string;
  serviceName: string;
  size?: string;
  quantity: number;
  amountKes: number;
  paymentMethod: string;
  mpesaPhone?: string;
  offlineChannel?: string;
  paymentReference?: string;
  status: string;
  notes?: string;
  createdAt: string;
}
