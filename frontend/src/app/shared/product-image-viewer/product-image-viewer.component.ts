import { DecimalPipe } from '@angular/common';
import { Component, HostListener, Input, computed, signal } from '@angular/core';
import { WebpImagePipe } from '../webp-image.pipe';

type ViewerMode = 'gallery' | '3d' | 'video';

@Component({
  selector: 'app-product-image-viewer',
  standalone: true,
  imports: [DecimalPipe, WebpImagePipe],
  templateUrl: './product-image-viewer.component.html',
  styleUrl: './product-image-viewer.component.scss',
})
export class ProductImageViewerComponent {
  @Input({ required: true }) alt = 'Product image';
  @Input() emoji = '✨';
  @Input() videoUrl?: string;
  @Input() set images(value: string[] | undefined) {
    const list = (value ?? []).filter(Boolean);
    this.imageList.set(list.length ? list : []);
    this.activeIndex.set(0);
    this.zoom.set(1);
    this.frameIndex.set(0);
  }
  @Input() set view360(value: string[] | undefined) {
    const list = (value ?? []).filter(Boolean);
    this.frames360.set(list.length ? list : []);
    this.frameIndex.set(0);
  }

  imageList = signal<string[]>([]);
  frames360 = signal<string[]>([]);
  activeIndex = signal(0);
  frameIndex = signal(0);
  zoom = signal(1);
  mode = signal<ViewerMode>('gallery');
  dragging3d = signal(false);
  dragStartX = 0;
  dragStartFrame = 0;

  readonly minZoom = 1;
  readonly maxZoom = 3;
  readonly zoomStep = 0.25;

  hasImages = computed(() => this.imageList().length > 0);
  has3d = computed(() => this.frames360().length > 1);
  hasVideo = computed(() => !!this.videoUrl);

  currentGallerySrc = computed(() => {
    const imgs = this.imageList();
    if (!imgs.length) return '';
    return imgs[this.activeIndex()] ?? imgs[0];
  });

  current3dSrc = computed(() => {
    const frames = this.frames360();
    if (!frames.length) return this.currentGallerySrc();
    const idx = Math.min(this.frameIndex(), frames.length - 1);
    return frames[idx] ?? frames[0];
  });

  displaySrc = computed(() => (this.mode() === '3d' ? this.current3dSrc() : this.currentGallerySrc()));

  setMode(next: ViewerMode): void {
    this.mode.set(next);
    this.zoom.set(1);
  }

  selectImage(index: number): void {
    this.activeIndex.set(index);
    this.zoom.set(1);
  }

  zoomIn(): void {
    this.zoom.update((z) => Math.min(this.maxZoom, +(z + this.zoomStep).toFixed(2)));
  }

  zoomOut(): void {
    this.zoom.update((z) => Math.max(this.minZoom, +(z - this.zoomStep).toFixed(2)));
  }

  onWheel(event: WheelEvent): void {
    if (this.mode() === 'video') return;
    event.preventDefault();
    if (event.deltaY < 0) {
      this.zoomIn();
    } else {
      this.zoomOut();
    }
  }

  on3dPointerDown(event: PointerEvent): void {
    if (this.mode() !== '3d' || !this.has3d()) return;
    this.dragging3d.set(true);
    this.dragStartX = event.clientX;
    this.dragStartFrame = this.frameIndex();
    (event.target as HTMLElement).setPointerCapture?.(event.pointerId);
  }

  on3dPointerMove(event: PointerEvent): void {
    if (!this.dragging3d()) return;
    const frames = this.frames360();
    if (frames.length < 2) return;
    const delta = event.clientX - this.dragStartX;
    const step = Math.round(delta / 28);
    const next = (this.dragStartFrame + step) % frames.length;
    this.frameIndex.set(next < 0 ? frames.length + next : next);
  }

  on3dPointerUp(): void {
    this.dragging3d.set(false);
  }

  @HostListener('document:keydown', ['$event'])
  onKey(event: KeyboardEvent): void {
    if (!this.hasImages() || this.mode() === 'video') return;
    if (event.key === '+' || event.key === '=') this.zoomIn();
    if (event.key === '-') this.zoomOut();
  }
}
