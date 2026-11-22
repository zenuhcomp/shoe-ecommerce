import { Component, OnInit, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { DecimalPipe, KeyValuePipe } from '@angular/common';
import { LashService } from '../../../core/models';
import { AdminApiService } from '../../services/admin-api.service';
import { ConfirmDialogService } from '../../../shared/confirm-dialog/confirm-dialog.service';

interface SpecRow {
  key: string;
  value: string;
}

@Component({
  selector: 'app-services-admin',
  standalone: true,
  imports: [ReactiveFormsModule, DecimalPipe, KeyValuePipe],
  templateUrl: './services-admin.component.html',
  styleUrl: './services-admin.component.scss',
})
export class ServicesAdminComponent implements OnInit {
  private readonly api = inject(AdminApiService);
  private readonly fb = inject(FormBuilder);
  private readonly dialog = inject(ConfirmDialogService);

  items = signal<LashService[]>([]);
  editing = signal<LashService | null>(null);
  showForm = signal(false);
  message = signal('');
  specsList = signal<SpecRow[]>([]);

  form = this.fb.nonNullable.group({
    name: ['', Validators.required],
    slug: [''],
    category: [''],
    badge: [''],
    emoji: ['✨'],
    style: [''],
    length: [''],
    priceKes: [0, Validators.required],
    originalPriceKes: [0],
    currency: ['KES'],
    featured: [false],
    inStock: [true],
    description: [''],
    imagesText: [''],
    view360Text: [''],
    videoUrl: [''],
  });

  ngOnInit(): void {
    this.load();
  }

  load(): void {
    this.api.listServices().subscribe((list) => this.items.set(list));
  }

  startCreate(): void {
    this.editing.set(null);
    this.showForm.set(true);
    this.specsList.set([
      { key: 'Material', value: '' },
      { key: 'Dimensions', value: '' },
    ]);
    this.form.reset({
      name: '',
      slug: '',
      category: '',
      badge: '',
      emoji: '✨',
      style: '',
      length: '',
      priceKes: 5000,
      originalPriceKes: 0,
      currency: 'KES',
      featured: false,
      inStock: true,
      description: '',
      imagesText: '',
      view360Text: '',
      videoUrl: '',
    });
  }

  startEdit(s: LashService): void {
    this.editing.set(s);
    this.showForm.set(true);
    const specsArr: SpecRow[] = [];
    if (s.specs) {
      for (const [k, v] of Object.entries(s.specs)) {
        specsArr.push({ key: k, value: v });
      }
    }
    this.specsList.set(specsArr);

    this.form.patchValue({
      name: s.name,
      slug: s.slug,
      category: s.category || s.style,
      badge: s.badge || '',
      emoji: s.emoji || '✨',
      style: s.style || '',
      length: s.length || '',
      priceKes: s.priceKes,
      originalPriceKes: s.originalPriceKes || 0,
      currency: s.currency || 'KES',
      featured: s.featured ?? false,
      inStock: s.inStock ?? true,
      description: s.description || '',
      imagesText: (s.images || []).join('\n'),
      view360Text: (s.view360 || []).join('\n'),
      videoUrl: s.videoUrl || '',
    });
  }

  addSpec(): void {
    this.specsList.update((list) => [...list, { key: '', value: '' }]);
  }

  removeSpec(index: number): void {
    this.specsList.update((list) => list.filter((_, i) => i !== index));
  }

  updateSpecKey(index: number, key: string): void {
    this.specsList.update((list) =>
      list.map((item, i) => (i === index ? { ...item, key } : item))
    );
  }

  updateSpecValue(index: number, value: string): void {
    this.specsList.update((list) =>
      list.map((item, i) => (i === index ? { ...item, value } : item))
    );
  }

  save(): void {
    if (this.form.invalid) return;
    const v = this.form.getRawValue();

    const specsMap: Record<string, string> = {};
    for (const item of this.specsList()) {
      const k = item.key.trim();
      const val = item.value.trim();
      if (k) specsMap[k] = val;
    }

    const images = v.imagesText
      .split('\n')
      .map((s) => s.trim())
      .filter(Boolean);

    const view360 = v.view360Text
      .split('\n')
      .map((s) => s.trim())
      .filter(Boolean);

    const payload: LashService = {
      id: this.editing()?.id ?? '',
      name: v.name,
      slug: v.slug || v.name.toLowerCase().replace(/\s+/g, '-').replace(/[^a-z0-9-]/g, ''),
      category: v.category || v.style,
      badge: v.badge,
      emoji: v.emoji,
      style: v.style || v.category,
      length: v.length,
      priceKes: Number(v.priceKes),
      originalPriceKes: Number(v.originalPriceKes),
      currency: v.currency || 'KES',
      specs: specsMap,
      featured: v.featured,
      inStock: v.inStock,
      description: v.description,
      images: images.length ? images : (this.editing()?.images ?? []),
      view360: view360.length ? view360 : (this.editing()?.view360 ?? []),
      videoUrl: v.videoUrl,
    };

    const req = this.editing()
      ? this.api.updateService(this.editing()!.id, payload)
      : this.api.createService(payload);

    req.subscribe({
      next: () => {
        this.message.set('Saved product successfully.');
        this.editing.set(null);
        this.showForm.set(false);
        this.load();
      },
      error: () => this.message.set('Save failed. Please check form fields.'),
    });
  }

  async remove(s: LashService): Promise<void> {
    const confirmed = await this.dialog.confirm({
      title: `Delete "${s.name}"?`,
      message: 'This product/service item will be permanently removed from the catalog.',
      confirmText: 'Delete Item',
      cancelText: 'Cancel',
      variant: 'danger',
      icon: '🗑️',
      badge: 'CATALOG ITEM',
    });
    if (!confirmed) return;
    this.api.deleteService(s.id).subscribe(() => this.load());
  }

  cancel(): void {
    this.editing.set(null);
    this.showForm.set(false);
  }
}
