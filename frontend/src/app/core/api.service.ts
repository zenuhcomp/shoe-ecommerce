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
} from './models';

@Injectable({ providedIn: 'root' })
export class ApiService {
  private readonly http = inject(HttpClient);
  private readonly base = '/api';

  private siteCache$?: Observable<SiteConfig>;
  private homeCache$?: Observable<HomePage>;

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

  applyTheme(site: SiteConfig): void {
    const t = site.theme;
    const root = document.documentElement;
    root.style.setProperty('--color-primary', t.primary);
    root.style.setProperty('--color-primary-dark', t.primaryDark);
    root.style.setProperty('--color-accent', t.accent);
    root.style.setProperty('--color-accent-light', t.accentLight);
    root.style.setProperty('--color-bg', t.background);
    root.style.setProperty('--color-surface', t.surface);
    root.style.setProperty('--color-text', t.text);
    root.style.setProperty('--color-text-muted', t.textMuted);
    root.style.setProperty('--hero-gradient', t.heroGradient);
    document.title = `${site.brandName} ${site.brandScript} – ${site.location}`;
  }

  preloadTheme(): Observable<SiteConfig> {
    return this.getSite().pipe(tap((site) => this.applyTheme(site)));
  }
}
