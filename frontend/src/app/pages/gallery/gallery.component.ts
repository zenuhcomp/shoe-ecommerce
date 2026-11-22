import { AsyncPipe, DecimalPipe } from '@angular/common';
import { Component, HostListener, inject, signal } from '@angular/core';
import { combineLatest, map } from 'rxjs';
import { ApiService } from '../../core/api.service';
import { GalleryItem } from '../../core/models';
import { PageHeroComponent } from '../../shared/page-hero/page-hero.component';
import { WebpImagePipe } from '../../shared/webp-image.pipe';

@Component({
  selector: 'app-gallery',
  standalone: true,
  imports: [AsyncPipe, DecimalPipe, PageHeroComponent, WebpImagePipe],
  template: `
    @if (vm$ | async; as vm) {
      <app-page-hero
        eyebrow="◆ Lookbook & Showcase"
        [title]="(vm.site.brandName || 'Visual') + ' <span class=\\'accent\\'>Showroom</span>'"
        [subtitle]="vm.site.tagline || 'High-resolution showcase and certified showroom lookbook'"
      />
      @if (vm.items; as items) {
        <section class="page-section gallery-grid">
          @for (item of items; track item.id; let i = $index) {
            <figure class="gallery-card">
              <button type="button" class="image-button" (click)="openViewer(items, i)" [attr.aria-label]="'Open ' + item.title">
                <img [src]="item.imageUrl | webpImage:800:75" [alt]="item.alt" loading="lazy" decoding="async" />
              </button>
              <figcaption>{{ item.title }}</figcaption>
            </figure>
          }
        </section>

        @if (activeItem(); as item) {
          <div class="lightbox" role="dialog" aria-modal="true" [attr.aria-label]="item.title" (click)="closeViewer()">
            <div class="lightbox-panel" (click)="$event.stopPropagation()">
              <div class="lightbox-toolbar">
                <button type="button" aria-label="Previous image" (click)="previous()">‹</button>
                <div class="zoom-controls">
                  <button type="button" aria-label="Zoom out" (click)="zoomOut()" [disabled]="zoom() <= minZoom">−</button>
                  <span>{{ zoom() * 100 | number: '1.0-0' }}%</span>
                  <button type="button" aria-label="Zoom in" (click)="zoomIn()" [disabled]="zoom() >= maxZoom">+</button>
                  <button type="button" aria-label="Reset zoom" (click)="resetZoom()">Reset</button>
                </div>
                <button type="button" aria-label="Next image" (click)="next()">›</button>
                <button type="button" class="close" aria-label="Close viewer" (click)="closeViewer()">×</button>
              </div>

              <div class="zoom-stage" (wheel)="onWheel($event)">
                <img
                  [src]="item.imageUrl | webpImage:1400:80"
                  [alt]="item.alt"
                  [style.transform]="'scale(' + zoom() + ')'"
                  draggable="false"
                />
              </div>

              <div class="lightbox-caption">
                <strong>{{ item.title }}</strong>
                <span>{{ activeIndex() + 1 }} / {{ viewerItems().length }}</span>
              </div>
            </div>
          </div>
        }
      }
    }
  `,
  styles: [
    `
      .page-section {
        max-width: var(--content-max);
        margin: 0 auto;
        padding: 2.5rem var(--page-gutter) 4.5rem;
      }
      .gallery-grid {
        display: grid;
        gap: 1.5rem;
        grid-template-columns: repeat(auto-fill, minmax(min(320px, 100%), 1fr));
      }
      .gallery-card {
        margin: 0;
        border-radius: 20px;
        overflow: hidden;
        background: var(--color-surface);
        color: var(--color-text);
        border: 1px solid color-mix(in srgb, var(--color-text) 12%, transparent);
        box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
        transition: transform 0.25s ease, box-shadow 0.25s ease, border-color 0.25s ease;

        &:hover {
          transform: translateY(-4px);
          box-shadow: 0 16px 36px rgba(0, 0, 0, 0.25);
          border-color: color-mix(in srgb, var(--color-accent) 40%, transparent);

          img {
            transform: scale(1.05);
          }
        }
      }
      .image-button {
        display: block;
        width: 100%;
        border: 0;
        padding: 0;
        background: transparent;
        cursor: zoom-in;
        overflow: hidden;
      }
      img {
        width: 100%;
        aspect-ratio: 16 / 10;
        object-fit: cover;
        display: block;
        transition: transform 0.4s cubic-bezier(0.16, 1, 0.3, 1);
      }
      figcaption {
        padding: 1rem 1.25rem;
        font-weight: 700;
        font-size: 1rem;
        color: var(--color-text);
      }
      .lightbox {
        position: fixed;
        inset: 0;
        z-index: 10000;
        background: rgba(5, 7, 10, 0.92);
        backdrop-filter: blur(12px);
        display: grid;
        place-items: center;
        padding: 1.5rem;
      }
      .lightbox-panel {
        width: min(1200px, 100%);
        max-height: min(90vh, 860px);
        display: grid;
        grid-template-rows: auto minmax(0, 1fr) auto;
        gap: 1rem;
      }
      .lightbox-toolbar {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.75rem;
      }
      .lightbox-toolbar button {
        border: 0;
        min-width: 2.75rem;
        min-height: 2.75rem;
        border-radius: 999px;
        background: rgba(255, 255, 255, 0.95);
        color: #0f172a;
        font-weight: 800;
        font-size: 1.1rem;
        cursor: pointer;
        transition: transform 0.15s, background 0.15s;

        &:hover:not(:disabled) {
          background: #ffffff;
          transform: scale(1.08);
        }
      }
      .lightbox-toolbar button:disabled {
        opacity: 0.4;
        cursor: not-allowed;
      }
      .zoom-controls {
        display: inline-flex;
        align-items: center;
        gap: 0.35rem;
        border-radius: 999px;
        background: rgba(255, 255, 255, 0.95);
        padding: 0.3rem 0.5rem;
      }
      .zoom-controls span {
        min-width: 3.5rem;
        text-align: center;
        color: #0f172a;
        font-size: 0.85rem;
        font-weight: 700;
      }
      .zoom-controls button {
        background: transparent;
        color: #0f172a;
        min-width: 2rem;
        min-height: 2rem;
      }
      .zoom-controls button:last-child {
        padding-inline: 0.85rem;
      }
      .close {
        margin-left: 0.5rem;
      }
      .zoom-stage {
        min-height: 0;
        display: grid;
        place-items: center;
        overflow: hidden;

        img {
          max-width: 100%;
          max-height: 70vh;
          object-fit: contain;
          aspect-ratio: auto;
          transition: transform 0.15s ease-out;
        }
      }
      .lightbox-caption {
        display: flex;
        justify-content: space-between;
        align-items: center;
        color: #ffffff;
        padding: 0.5rem 1rem;
        font-size: 1rem;

        strong {
          color: var(--color-accent-light, #bae6fd);
        }

        span {
          opacity: 0.7;
          font-size: 0.88rem;
        }
      }
    `,
  ],
})
export class GalleryComponent {
  private readonly api = inject(ApiService);
  vm$ = combineLatest([this.api.getSite(), this.api.getGallery()]).pipe(
    map(([site, items]) => ({ site, items }))
  );

  viewerItems = signal<GalleryItem[]>([]);
  activeIndex = signal(0);
  activeItem = signal<GalleryItem | null>(null);
  zoom = signal(1);

  readonly minZoom = 0.5;
  readonly maxZoom = 3;

  openViewer(items: GalleryItem[], index: number): void {
    this.viewerItems.set(items);
    this.activeIndex.set(index);
    this.activeItem.set(items[index] ?? null);
    this.zoom.set(1);
  }

  closeViewer(): void {
    this.activeItem.set(null);
  }

  next(): void {
    const list = this.viewerItems();
    if (!list.length) return;
    const nextIdx = (this.activeIndex() + 1) % list.length;
    this.activeIndex.set(nextIdx);
    this.activeItem.set(list[nextIdx]);
    this.zoom.set(1);
  }

  previous(): void {
    const list = this.viewerItems();
    if (!list.length) return;
    const prevIdx = (this.activeIndex() - 1 + list.length) % list.length;
    this.activeIndex.set(prevIdx);
    this.activeItem.set(list[prevIdx]);
    this.zoom.set(1);
  }

  zoomIn(): void {
    this.zoom.update((z) => Math.min(this.maxZoom, Number((z + 0.25).toFixed(2))));
  }

  zoomOut(): void {
    this.zoom.update((z) => Math.max(this.minZoom, Number((z - 0.25).toFixed(2))));
  }

  resetZoom(): void {
    this.zoom.set(1);
  }

  onWheel(event: WheelEvent): void {
    event.preventDefault();
    if (event.deltaY < 0) {
      this.zoomIn();
    } else {
      this.zoomOut();
    }
  }

  @HostListener('document:keydown.escape')
  onEscape(): void {
    if (this.activeItem()) {
      this.closeViewer();
    }
  }

  @HostListener('document:keydown.arrowleft')
  onArrowLeft(): void {
    if (this.activeItem()) {
      this.previous();
    }
  }

  @HostListener('document:keydown.arrowright')
  onArrowRight(): void {
    if (this.activeItem()) {
      this.next();
    }
  }
}
