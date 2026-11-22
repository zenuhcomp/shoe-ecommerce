import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import {
  FAQ,
  Feature,
  GalleryItem,
  LashService,
  NavLink,
  PressMention,
  Review,
  SiteConfig,
  Section,
  IndustryPreset,
  MediaAsset,
  OrderRecord,
  PaymentConfig,
  PaymentMethodInfo,
} from '../../core/models';
import { AdminAuthService } from './admin-auth.service';

export interface DashboardResponse {
  counts: Record<string, number>;
  recentBookings: BookingRecord[];
  recentOrders?: OrderRecord[];
}

@Injectable({ providedIn: 'root' })
export class AdminApiService {
  private readonly http = inject(HttpClient);
  private readonly auth = inject(AdminAuthService);
  private readonly base = '/api/admin';

  private headers(): HttpHeaders {
    const t = this.auth.token();
    return new HttpHeaders(t ? { Authorization: `Bearer ${t}` } : {});
  }

  dashboard(): Observable<DashboardResponse> {
    return this.http.get<DashboardResponse>(
      `${this.base}/dashboard`,
      { headers: this.headers() }
    );
  }

  // Presets
  listPresets(): Observable<IndustryPreset[]> {
    return this.http.get<IndustryPreset[]>(`${this.base}/presets`, { headers: this.headers() });
  }

  applyPreset(presetId: string, seedItems: boolean = true): Observable<{ success: boolean; presetId: string; message: string }> {
    return this.http.post<{ success: boolean; presetId: string; message: string }>(
      `${this.base}/presets/apply`,
      { presetId, seedItems },
      { headers: this.headers() }
    );
  }

  // Sections
  listSections(): Observable<Section[]> {
    return this.http.get<Section[]>(`${this.base}/sections`, { headers: this.headers() });
  }

  createSection(s: Partial<Section>): Observable<Section> {
    return this.http.post<Section>(`${this.base}/sections`, s, { headers: this.headers() });
  }

  updateSection(id: string, s: Section): Observable<Section> {
    return this.http.put<Section>(`${this.base}/sections/${id}`, s, { headers: this.headers() });
  }

  deleteSection(id: string): Observable<void> {
    return this.http.delete<void>(`${this.base}/sections/${id}`, { headers: this.headers() });
  }

  reorderSections(ids: string[]): Observable<void> {
    return this.http.put<void>(`${this.base}/sections/reorder`, { ids }, { headers: this.headers() });
  }

  // Media
  listMedia(category?: string): Observable<MediaAsset[]> {
    const url = category ? `${this.base}/media?category=${encodeURIComponent(category)}` : `${this.base}/media`;
    return this.http.get<MediaAsset[]>(url, { headers: this.headers() });
  }

  createMedia(m: Partial<MediaAsset>): Observable<MediaAsset> {
    return this.http.post<MediaAsset>(`${this.base}/media`, m, { headers: this.headers() });
  }

  deleteMedia(id: string): Observable<void> {
    return this.http.delete<void>(`${this.base}/media/${id}`, { headers: this.headers() });
  }

  // Orders
  listOrders(): Observable<OrderRecord[]> {
    return this.http.get<OrderRecord[]>(`${this.base}/orders`, { headers: this.headers() });
  }

  updateOrderStatus(id: string, status: string): Observable<{ id: string; status: string }> {
    return this.http.patch<{ id: string; status: string }>(
      `${this.base}/orders/${id}/status`,
      { status },
      { headers: this.headers() }
    );
  }

  deleteOrder(id: string): Observable<void> {
    return this.http.delete<void>(`${this.base}/orders/${id}`, { headers: this.headers() });
  }

  // Site
  getSite(): Observable<SiteConfig> {
    return this.http.get<SiteConfig>(`${this.base}/site`, { headers: this.headers() });
  }

  putSite(site: SiteConfig): Observable<SiteConfig> {
    return this.http.put<SiteConfig>(`${this.base}/site`, site, { headers: this.headers() });
  }

  // Services / Products
  listServices(): Observable<LashService[]> {
    return this.http.get<LashService[]>(`${this.base}/services`, { headers: this.headers() });
  }
  createService(s: LashService): Observable<LashService> {
    return this.http.post<LashService>(`${this.base}/services`, s, { headers: this.headers() });
  }
  updateService(id: string, s: LashService): Observable<LashService> {
    return this.http.put<LashService>(`${this.base}/services/${id}`, s, { headers: this.headers() });
  }
  deleteService(id: string): Observable<void> {
    return this.http.delete<void>(`${this.base}/services/${id}`, { headers: this.headers() });
  }

  // Features
  listFeatures(): Observable<Feature[]> {
    return this.http.get<Feature[]>(`${this.base}/features`, { headers: this.headers() });
  }
  createFeature(f: Feature): Observable<Feature> {
    return this.http.post<Feature>(`${this.base}/features`, f, { headers: this.headers() });
  }
  updateFeature(id: string, f: Feature): Observable<Feature> {
    return this.http.put<Feature>(`${this.base}/features/${id}`, f, { headers: this.headers() });
  }
  deleteFeature(id: string): Observable<void> {
    return this.http.delete<void>(`${this.base}/features/${id}`, { headers: this.headers() });
  }

  // Press
  listPress(): Observable<PressMention[]> {
    return this.http.get<PressMention[]>(`${this.base}/press`, { headers: this.headers() });
  }
  createPress(p: PressMention): Observable<PressMention> {
    return this.http.post<PressMention>(`${this.base}/press`, p, { headers: this.headers() });
  }
  updatePress(id: string, p: PressMention): Observable<PressMention> {
    return this.http.put<PressMention>(`${this.base}/press/${id}`, p, { headers: this.headers() });
  }
  deletePress(id: string): Observable<void> {
    return this.http.delete<void>(`${this.base}/press/${id}`, { headers: this.headers() });
  }

  // Reviews
  listReviews(): Observable<Review[]> {
    return this.http.get<Review[]>(`${this.base}/reviews`, { headers: this.headers() });
  }
  createReview(r: Review): Observable<Review> {
    return this.http.post<Review>(`${this.base}/reviews`, r, { headers: this.headers() });
  }
  updateReview(id: string, r: Review): Observable<Review> {
    return this.http.put<Review>(`${this.base}/reviews/${id}`, r, { headers: this.headers() });
  }
  deleteReview(id: string): Observable<void> {
    return this.http.delete<void>(`${this.base}/reviews/${id}`, { headers: this.headers() });
  }

  // FAQ
  listFaq(): Observable<FAQ[]> {
    return this.http.get<FAQ[]>(`${this.base}/faq`, { headers: this.headers() });
  }
  createFaq(faq: FAQ): Observable<FAQ> {
    return this.http.post<FAQ>(`${this.base}/faq`, faq, { headers: this.headers() });
  }
  updateFaq(id: number, faq: FAQ): Observable<FAQ> {
    return this.http.put<FAQ>(`${this.base}/faq/${id}`, faq, { headers: this.headers() });
  }
  deleteFaq(id: number): Observable<void> {
    return this.http.delete<void>(`${this.base}/faq/${id}`, { headers: this.headers() });
  }

  // Gallery
  listGallery(): Observable<GalleryItem[]> {
    return this.http.get<GalleryItem[]>(`${this.base}/gallery`, { headers: this.headers() });
  }
  createGallery(g: GalleryItem): Observable<GalleryItem> {
    return this.http.post<GalleryItem>(`${this.base}/gallery`, g, { headers: this.headers() });
  }
  updateGallery(id: string, g: GalleryItem): Observable<GalleryItem> {
    return this.http.put<GalleryItem>(`${this.base}/gallery/${id}`, g, { headers: this.headers() });
  }
  deleteGallery(id: string): Observable<void> {
    return this.http.delete<void>(`${this.base}/gallery/${id}`, { headers: this.headers() });
  }

  // Nav
  listNav(): Observable<NavLink[]> {
    return this.http.get<NavLink[]>(`${this.base}/nav`, { headers: this.headers() });
  }
  createNav(n: NavLink): Observable<NavLink> {
    return this.http.post<NavLink>(`${this.base}/nav`, n, { headers: this.headers() });
  }
  updateNav(id: string, n: NavLink): Observable<NavLink> {
    return this.http.put<NavLink>(`${this.base}/nav/${id}`, n, { headers: this.headers() });
  }
  deleteNav(id: string): Observable<void> {
    return this.http.delete<void>(`${this.base}/nav/${id}`, { headers: this.headers() });
  }

  // Bookings / Inquiries
  listBookings(): Observable<BookingRecord[]> {
    return this.http.get<BookingRecord[]>(`${this.base}/bookings`, { headers: this.headers() });
  }
  deleteBooking(id: string): Observable<void> {
    return this.http.delete<void>(`${this.base}/bookings/${id}`, { headers: this.headers() });
  }

  // ── Payment Gateway Configuration ─────────────────────────────────────────
  getPaymentConfig(): Observable<{ config: PaymentConfig; methods: PaymentMethodInfo[] }> {
    return this.http.get<{ config: PaymentConfig; methods: PaymentMethodInfo[] }>(
      `${this.base}/payments`,
      { headers: this.headers() }
    );
  }

  updatePaymentConfig(cfg: PaymentConfig): Observable<{ success: boolean; message: string; config: PaymentConfig; methods: PaymentMethodInfo[] }> {
    return this.http.put<{ success: boolean; message: string; config: PaymentConfig; methods: PaymentMethodInfo[] }>(
      `${this.base}/payments`,
      cfg,
      { headers: this.headers() }
    );
  }

  simulateCallback(payload: { orderId: string; checkoutRequestId?: string; resultCode?: number; receiptNumber?: string; amount?: number; phone?: string }): Observable<any> {
    return this.http.post<any>(`${this.base}/payments/simulate-callback`, payload, { headers: this.headers() });
  }
}

export interface BookingRecord {
  id: string;
  name: string;
  phone: string;
  email?: string;
  serviceSlug?: string;
  preferredAt?: string;
  message?: string;
  createdAt: string;
}
