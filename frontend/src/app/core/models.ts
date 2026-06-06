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

export interface SiteConfig {
  brandName: string;
  brandScript: string;
  tagline: string;
  heroTitle: string;
  heroSubtitle: string;
  heroBadge: string;
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
  about: AboutSection;
  houseCalls: HouseCallsSection;
  aftercare: AftercareSection;
  pricingHighlight: PricingHighlight;
  payment?: PaymentConfig;
}

export interface PaymentConfig {
  tillNumber: string;
  paybillNumber: string;
  paybillAccount: string;
  mpesaShortCode: string;
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
  description: string;
  priceKes: number;
  featured?: boolean;
  images?: string[];
  view360?: string[];
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
