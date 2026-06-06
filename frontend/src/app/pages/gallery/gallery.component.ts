import { AsyncPipe, DecimalPipe } from '@angular/common';
import { Component, HostListener, inject, signal } from '@angular/core';
import { ApiService } from '../../core/api.service';
import { GalleryItem } from '../../core/models';
import { PageHeroComponent } from '../../shared/page-hero/page-hero.component';
import { WebpImagePipe } from '../../shared/webp-image.pipe';

@Component({
  selector: 'app-gallery',
  standalone: true,
  imports: [AsyncPipe, DecimalPipe, PageHeroComponent, WebpImagePipe],
  template: `
    <app-page-hero
      eyebrow="◆ Lookbook"
      title="Shoe <span class='accent'>Gallery</span>"
      subtitle="Official styles styled for the modern professional wardrobe"
    />
    @if (items$ | async; as items) {
      <section class="page-section gallery-grid">
        @for (item of items; track item.id; let i = $index) {
          <figure class="gallery-card">
            <button type="button" class="image-button" (click)="openViewer(items, i)" [attr.aria-label]="'Open ' + item.title">
              <img [src]="item.imageUrl | webpImage:620:70" [alt]="item.alt" loading="lazy" decoding="async" />
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
                [src]="item.imageUrl | webpImage:1300:76"
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
  `,
  styles: [
    `
      .page-section {
        max-width: var(--content-max);
        margin: 0 auto;
        padding: 2rem var(--page-gutter) 4rem;
      }
      .gallery-grid {
        display: grid;
        gap: 1rem;
        grid-template-columns: repeat(auto-fill, minmax(min(260px, 100%), 1fr));
      }
      .gallery-card {
        margin: 0;
        border-radius: 16px;
        overflow: hidden;
        background: var(--color-surface);
        box-shadow: 0 8px 24px rgba(15, 18, 36, 0.08);
      }
      .image-button {
        display: block;
        width: 100%;
        border: 0;
        padding: 0;
        background: transparent;
        cursor: zoom-in;
      }
      img {
        width: 100%;
        aspect-ratio: 4/5;
        object-fit: cover;
        display: block;
      }
      figcaption {
        padding: 0.75rem 1rem;
        font-weight: 600;
        font-size: 0.9rem;
      }
      .lightbox {
        position: fixed;
        inset: 0;
        z-index: 80;
        background: rgba(15, 18, 36, 0.86);
        display: grid;
        place-items: center;
        padding: 1rem;
      }
      .lightbox-panel {
        width: min(1120px, 100%);
        max-height: min(90vh, 820px);
        display: grid;
        grid-template-rows: auto minmax(0, 1fr) auto;
        gap: 0.75rem;
      }
      .lightbox-toolbar {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
      }
      .lightbox-toolbar button {
        border: 0;
        min-width: 2.5rem;
        min-height: 2.5rem;
        border-radius: 999px;
        background: #fff;
        color: var(--color-primary);
        font-weight: 800;
        cursor: pointer;
      }
      .lightbox-toolbar button:disabled {
        opacity: 0.45;
        cursor: not-allowed;
      }
      .zoom-controls {
        display: inline-flex;
        align-items: center;
        gap: 0.35rem;
        border-radius: 999px;
        background: rgba(255, 255, 255, 0.94);
        padding: 0.25rem;
      }
      .zoom-controls span {
        min-width: 3.5rem;
        text-align: center;
        color: var(--color-primary);
        font-size: 0.85rem;
        font-weight: 700;
      }
      .zoom-controls button {
        background: transparent;
      }
      .zoom-controls button:last-child {
        padding-inline: 0.85rem;
      }
      .close {
        margin-left: 0.25rem;
      }
      .zoom-stage {
        min-height: 0;
        height: min(72vh, 680px);
        border-radius: 12px;
        overflow: auto;
        background: #fff;
        display: grid;
        place-items: center;
        overscroll-behavior: contain;
      }
      .zoom-stage img {
        width: 100%;
        height: 100%;
        max-width: 100%;
        max-height: 100%;
        aspect-ratio: auto;
        object-fit: contain;
        transform-origin: center center;
        transition: transform 0.15s ease;
        user-select: none;
      }
      .lightbox-caption {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 1rem;
        color: #fff;
      }
      .lightbox-caption span {
        color: rgba(255, 255, 255, 0.75);
        font-size: 0.9rem;
      }
      @media (min-width: 768px) {
        .gallery-grid {
          grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
        }
      }
      @media (max-width: 640px) {
        .lightbox {
          padding: 0.75rem;
        }
        .lightbox-toolbar {
          justify-content: flex-start;
          overflow-x: auto;
          padding-bottom: 0.25rem;
        }
        .zoom-stage {
          height: 68vh;
        }
      }
    `,
  ],
})
export class GalleryComponent {
  private readonly api = inject(ApiService);
  items$ = this.api.getGallery();

  viewerItems = signal<GalleryItem[]>([]);
  activeIndex = signal(0);
  zoom = signal(1);

  readonly minZoom = 1;
  readonly maxZoom = 4;
  readonly zoomStep = 0.25;

  activeItem = () => this.viewerItems()[this.activeIndex()] ?? null;

  openViewer(items: GalleryItem[], index: number): void {
    this.viewerItems.set(items);
    this.activeIndex.set(index);
    this.resetZoom();
  }

  closeViewer(): void {
    this.viewerItems.set([]);
    this.activeIndex.set(0);
    this.resetZoom();
  }

  next(): void {
    const count = this.viewerItems().length;
    if (!count) return;
    this.activeIndex.update((index) => (index + 1) % count);
    this.resetZoom();
  }

  previous(): void {
    const count = this.viewerItems().length;
    if (!count) return;
    this.activeIndex.update((index) => (index - 1 + count) % count);
    this.resetZoom();
  }

  zoomIn(): void {
    this.zoom.update((value) => Math.min(this.maxZoom, +(value + this.zoomStep).toFixed(2)));
  }

  zoomOut(): void {
    this.zoom.update((value) => Math.max(this.minZoom, +(value - this.zoomStep).toFixed(2)));
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

  @HostListener('document:keydown', ['$event'])
  onKey(event: KeyboardEvent): void {
    if (!this.activeItem()) return;

    if (event.key === 'Escape') this.closeViewer();
    if (event.key === 'ArrowRight') this.next();
    if (event.key === 'ArrowLeft') this.previous();
    if (event.key === '+' || event.key === '=') this.zoomIn();
    if (event.key === '-') this.zoomOut();
    if (event.key === '0') this.resetZoom();
  }
}
