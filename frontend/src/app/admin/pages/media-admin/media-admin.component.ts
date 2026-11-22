import { Component, OnInit, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { AdminApiService } from '../../services/admin-api.service';
import { ConfirmDialogService } from '../../../shared/confirm-dialog/confirm-dialog.service';
import { MediaAsset } from '../../../core/models';

@Component({
  selector: 'app-media-admin',
  standalone: true,
  imports: [ReactiveFormsModule],
  templateUrl: './media-admin.component.html',
  styleUrl: './media-admin.component.scss',
})
export class MediaAdminComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly api = inject(AdminApiService);
  private readonly dialog = inject(ConfirmDialogService);

  mediaList = signal<MediaAsset[]>([]);
  activeCategory = signal<string>('all');
  activeType = signal<'all' | 'image' | 'video'>('all');
  loading = signal(false);
  isAdding = signal(false);
  copiedId = signal<string | null>(null);
  message = signal<{ type: 'success' | 'error'; text: string } | null>(null);

  addForm = this.fb.group({
    name: ['', Validators.required],
    category: ['furniture', Validators.required],
    mediaType: ['image', Validators.required],
    url: ['', Validators.required],
    altText: [''],
    tags: [''],
  });

  ngOnInit(): void {
    this.loadMedia();
  }

  loadMedia(): void {
    this.loading.set(true);
    this.api.listMedia(this.activeCategory()).subscribe({
      next: (list) => {
        this.mediaList.set(list);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
        this.message.set({ type: 'error', text: 'Could not load media assets' });
      },
    });
  }

  setCategory(cat: string): void {
    this.activeCategory.set(cat);
    this.loadMedia();
  }

  setTypeFilter(type: 'all' | 'image' | 'video'): void {
    this.activeType.set(type);
  }

  filteredMedia(): MediaAsset[] {
    const type = this.activeType();
    const list = this.mediaList();
    if (type === 'all') return list;
    return list.filter((m) => m.mediaType === type);
  }

  copyUrl(m: MediaAsset): void {
    navigator.clipboard.writeText(m.url).then(() => {
      this.copiedId.set(m.id);
      setTimeout(() => this.copiedId.set(null), 2000);
    });
  }

  startAdd(): void {
    this.isAdding.set(true);
    this.addForm.reset({
      name: '',
      category: this.activeCategory() === 'all' ? 'furniture' : this.activeCategory(),
      mediaType: 'image',
      url: '',
      altText: '',
      tags: '',
    });
  }

  cancelAdd(): void {
    this.isAdding.set(false);
  }

  saveMedia(): void {
    if (this.addForm.invalid) return;
    const val = this.addForm.getRawValue() as Partial<MediaAsset>;
    this.api.createMedia(val).subscribe({
      next: (created) => {
        this.mediaList.update((list) => [created, ...list]);
        this.isAdding.set(false);
        this.message.set({ type: 'success', text: 'Asset added to media library!' });
      },
      error: () => this.message.set({ type: 'error', text: 'Could not save asset' }),
    });
  }

  async deleteMedia(id: string): Promise<void> {
    const confirmed = await this.dialog.confirm({
      title: 'Remove Media Asset?',
      message: 'This asset will be permanently removed from your media library.',
      confirmText: 'Delete Asset',
      cancelText: 'Cancel',
      variant: 'danger',
      icon: '🗑️',
      badge: 'MEDIA LIBRARY',
    });
    if (!confirmed) return;
    this.api.deleteMedia(id).subscribe({
      next: () => {
        this.mediaList.update((list) => list.filter((m) => m.id !== id));
        this.message.set({ type: 'success', text: 'Asset removed.' });
      },
      error: () => this.message.set({ type: 'error', text: 'Could not delete asset' }),
    });
  }
}
