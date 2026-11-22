import { DecimalPipe } from '@angular/common';
import { Component, Input, computed, signal } from '@angular/core';
import { RouterLink } from '@angular/router';
import { LashService, SiteConfig } from '../../core/models';
import { WebpImagePipe } from '../webp-image.pipe';
import { ProductImageViewerComponent } from '../product-image-viewer/product-image-viewer.component';

@Component({
  selector: 'app-service-grid',
  standalone: true,
  imports: [RouterLink, DecimalPipe, WebpImagePipe, ProductImageViewerComponent],
  templateUrl: './service-grid.component.html',
  styleUrl: './service-grid.component.scss',
})
export class ServiceGridComponent {
  @Input({ required: true }) set services(list: LashService[]) {
    this._services.set(list || []);
  }
  @Input() site?: SiteConfig;
  @Input() compact = false;
  @Input() showCategoryFilter = true;

  readonly _services = signal<LashService[]>([]);
  activeCategory = signal<string>('all');
  selectedItemForViewer = signal<LashService | null>(null);

  categories = computed(() => {
    const set = new Set<string>();
    for (const s of this._services()) {
      const cat = s.category || s.style;
      if (cat) set.add(cat);
    }
    return Array.from(set);
  });

  filteredServices = computed(() => {
    const cat = this.activeCategory();
    const list = this._services();
    if (cat === 'all') return list;
    return list.filter((s) => (s.category || s.style) === cat);
  });

  setCategory(cat: string): void {
    this.activeCategory.set(cat);
  }

  openViewer(svc: LashService): void {
    this.selectedItemForViewer.set(svc);
  }

  closeViewer(): void {
    this.selectedItemForViewer.set(null);
  }

  getSpecsArray(specs?: Record<string, string>): { key: string; value: string }[] {
    if (!specs) return [];
    return Object.entries(specs).map(([key, value]) => ({ key, value }));
  }

  getWhatsAppLink(svc: LashService): string {
    if (!this.site?.whatsAppUrl) return '';
    const text = encodeURIComponent(`Hi ${this.site.brandName}, I'm interested in "${svc.name}" (${svc.currency || 'KES'} ${svc.priceKes}). Could you please share more details?`);
    const sep = this.site.whatsAppUrl.includes('?') ? '&' : '?';
    return `${this.site.whatsAppUrl}${sep}text=${text}`;
  }
}
