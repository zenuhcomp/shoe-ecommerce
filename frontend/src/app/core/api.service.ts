import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable, shareReplay, tap } from 'rxjs';
import {
  BookingInquiry,
  BookingResponse,
  OrderRequest,
  OrderResponse,
  FAQ,
  GalleryItem,
  HomePage,
  LashService,
  Review,
  SiteConfig,
  Section,
  IndustryPreset,
  MediaAsset,
  PaymentMethodInfo,
  PaymentInitiateRequest,
  PaymentInitiateResponse,
  PaymentVerifyRequest,
} from './models';

@Injectable({ providedIn: 'root' })
export class ApiService {
  private readonly http = inject(HttpClient);
  private readonly base = '/api';

  private siteCache$?: Observable<SiteConfig>;
  private homeCache$?: Observable<HomePage>;

  clearCache(): void {
    this.siteCache$ = undefined;
    this.homeCache$ = undefined;
  }

  getSite(): Observable<SiteConfig> {
    if (!this.siteCache$) {
      this.siteCache$ = this.http.get<SiteConfig>(`${this.base}/site`).pipe(shareReplay(1));
    }
    return this.siteCache$;
  }

  getHome(): Observable<HomePage> {
    if (!this.homeCache$) {
      this.homeCache$ = this.http.get<HomePage>(`${this.base}/home`).pipe(shareReplay(1));
    }
    return this.homeCache$;
  }

  getSections(): Observable<Section[]> {
    return this.http.get<Section[]>(`${this.base}/sections`);
  }

  getPresets(): Observable<IndustryPreset[]> {
    return this.http.get<IndustryPreset[]>(`${this.base}/presets`);
  }

  getMedia(category?: string): Observable<MediaAsset[]> {
    const url = category ? `${this.base}/media?category=${encodeURIComponent(category)}` : `${this.base}/media`;
    return this.http.get<MediaAsset[]>(url);
  }

  getServices(): Observable<LashService[]> {
    return this.http.get<LashService[]>(`${this.base}/services`);
  }

  getReviews(): Observable<Review[]> {
    return this.http.get<Review[]>(`${this.base}/reviews`);
  }

  getFaq(): Observable<FAQ[]> {
    return this.http.get<FAQ[]>(`${this.base}/faq`);
  }

  getGallery(): Observable<GalleryItem[]> {
    return this.http.get<GalleryItem[]>(`${this.base}/gallery`);
  }

  submitBooking(inquiry: BookingInquiry): Observable<BookingResponse> {
    return this.http.post<BookingResponse>(`${this.base}/bookings`, inquiry);
  }

  getService(slug: string): Observable<LashService> {
    return this.http.get<LashService>(`${this.base}/services/${slug}`);
  }

  submitOrder(order: OrderRequest): Observable<OrderResponse> {
    return this.http.post<OrderResponse>(`${this.base}/orders`, order);
  }

  // ── Payment Strategy API ──────────────────────────────────────────────────
  getActivePaymentMethods(): Observable<{ methods: PaymentMethodInfo[] }> {
    return this.http.get<{ methods: PaymentMethodInfo[] }>(`${this.base}/payments/methods`);
  }

  initiatePayment(req: PaymentInitiateRequest): Observable<PaymentInitiateResponse> {
    return this.http.post<PaymentInitiateResponse>(`${this.base}/payments/initiate`, req);
  }

  verifyPayment(req: PaymentVerifyRequest): Observable<PaymentInitiateResponse> {
    return this.http.post<PaymentInitiateResponse>(`${this.base}/payments/verify`, req);
  }

  getOrderStatus(orderId: string): Observable<{ id: string; status: string; paymentMethod: string; paymentReference?: string; amountKes: number; notes?: string }> {
    return this.http.get<{ id: string; status: string; paymentMethod: string; paymentReference?: string; amountKes: number; notes?: string }>(
      `${this.base}/payments/order-status?orderId=${encodeURIComponent(orderId)}`
    );
  }

  simulateMpesaCallback(payload: { orderId: string; checkoutRequestId?: string; resultCode?: number; resultDesc?: string; receiptNumber?: string; amount?: number; phone?: string }): Observable<any> {
    return this.http.post<any>(`${this.base}/payments/mpesa/simulate-callback`, payload);
  }

  applyTheme(site: SiteConfig): void {
    const t = site.theme;
    const root = document.documentElement;
    if (t.primary) root.style.setProperty('--color-primary', t.primary);
    if (t.primaryDark) root.style.setProperty('--color-primary-dark', t.primaryDark);
    if (t.accent) root.style.setProperty('--color-accent', t.accent);
    if (t.accentLight) root.style.setProperty('--color-accent-light', t.accentLight);
    if (t.background) root.style.setProperty('--color-bg', t.background);
    if (t.surface) root.style.setProperty('--color-surface', t.surface);
    if (t.text) root.style.setProperty('--color-text', t.text);
    if (t.textMuted) root.style.setProperty('--color-text-muted', t.textMuted);
    if (t.heroGradient) root.style.setProperty('--hero-gradient', t.heroGradient);
    if (t.fontHeading) root.style.setProperty('--font-heading', t.fontHeading);
    if (t.fontBody) root.style.setProperty('--font-body', t.fontBody);
    if (t.fontScript) root.style.setProperty('--font-script', t.fontScript);

    if (site.animations) {
      root.style.setProperty('--anim-duration', site.animations.defaultDuration || '0.8s');
      if (!site.animations.globalEnabled || site.animations.intensity === 'off') {
        root.style.setProperty('--anim-multiplier', '0');
      } else if (site.animations.intensity === 'subtle') {
        root.style.setProperty('--anim-multiplier', '0.6');
      } else if (site.animations.intensity === 'energetic') {
        root.style.setProperty('--anim-multiplier', '1.3');
      } else {
        root.style.setProperty('--anim-multiplier', '1');
      }
    }

    document.title = `${site.brandName} ${site.brandScript} – ${site.location}`;
  }

  preloadTheme(): Observable<SiteConfig> {
    return this.getSite().pipe(tap((site) => this.applyTheme(site)));
  }
}
