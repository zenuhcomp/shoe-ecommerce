import { DatePipe } from '@angular/common';
import { Component, OnInit, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { AdminApiService } from '../../services/admin-api.service';
import { ConfirmDialogService } from '../../../shared/confirm-dialog/confirm-dialog.service';

type EntityType = 'features' | 'press' | 'reviews' | 'faq' | 'gallery' | 'nav' | 'bookings';

@Component({
  selector: 'app-content-list',
  standalone: true,
  imports: [ReactiveFormsModule, DatePipe],
  templateUrl: './content-list.component.html',
  styleUrl: './content-list.component.scss',
})
export class ContentListComponent implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly api = inject(AdminApiService);
  private readonly fb = inject(FormBuilder);
  private readonly dialog = inject(ConfirmDialogService);

  type = signal<EntityType>('features');
  title = signal('Content');
  items = signal<Record<string, unknown>[]>([]);
  editingId = signal<string | number | null>(null);
  showForm = signal(false);
  message = signal('');

  form = this.fb.group({
    name: [''],
    icon: [''],
    title: [''],
    description: [''],
    note: [''],
    handle: [''],
    initials: [''],
    rating: [5],
    quote: [''],
    published: [true],
    question: [''],
    answer: [''],
    category: [''],
    imageUrl: [''],
    alt: [''],
    label: [''],
    path: [''],
    visible: [true],
    priceKes: [0],
    slug: [''],
    style: [''],
    length: [''],
    featured: [false],
    emoji: [''],
  });

  ngOnInit(): void {
    const t = this.route.snapshot.data['entity'] as EntityType;
    this.type.set(t);
    this.title.set(this.labels[t] ?? 'Content');
    this.load();
  }

  private readonly labels: Record<EntityType, string> = {
    features: 'Features',
    press: 'Press mentions',
    reviews: 'Reviews',
    faq: 'FAQ',
    gallery: 'Gallery',
    nav: 'Navigation links',
    bookings: 'Booking inquiries',
  };

  load(): void {
    const t = this.type();
    const handlers: Record<EntityType, () => void> = {
      features: () => this.api.listFeatures().subscribe((d) => this.items.set(d as never)),
      press: () => this.api.listPress().subscribe((d) => this.items.set(d as never)),
      reviews: () => this.api.listReviews().subscribe((d) => this.items.set(d as never)),
      faq: () => this.api.listFaq().subscribe((d) => this.items.set(d as never)),
      gallery: () => this.api.listGallery().subscribe((d) => this.items.set(d as never)),
      nav: () => this.api.listNav().subscribe((d) => this.items.set(d as never)),
      bookings: () => this.api.listBookings().subscribe((d) => this.items.set(d as never)),
    };
    handlers[t]();
  }

  startCreate(): void {
    if (this.type() === 'bookings') return;
    this.editingId.set(null);
    this.showForm.set(true);
    this.form.reset();
  }

  startEdit(row: Record<string, unknown>): void {
    if (this.type() === 'bookings') return;
    this.editingId.set(row['id'] as string | number);
    this.showForm.set(true);
    this.form.patchValue(row as never);
  }

  save(): void {
    const t = this.type();
    const v = this.form.getRawValue();
    const id = this.editingId();

    const done = () => {
      this.message.set('Saved.');
      this.showForm.set(false);
      this.load();
    };
    const fail = () => this.message.set('Save failed.');

    if (t === 'features') {
      const body = { id: String(id ?? ''), icon: v.icon ?? '', title: v.title ?? '', description: v.description ?? '' };
      if (!body.title) return;
      (id ? this.api.updateFeature(String(id), body) : this.api.createFeature(body)).subscribe({ next: done, error: fail });
    } else if (t === 'press') {
      const body = { id: String(id ?? ''), name: v.name ?? '', note: v.note ?? '' };
      if (!body.name) return;
      (id ? this.api.updatePress(String(id), body) : this.api.createPress(body)).subscribe({ next: done, error: fail });
    } else if (t === 'reviews') {
      const body = {
        id: String(id ?? ''),
        name: v.name ?? '',
        handle: v.handle ?? '',
        initials: v.initials ?? '',
        rating: Number(v.rating) || 5,
        quote: v.quote ?? '',
        published: !!v.published,
      };
      (id ? this.api.updateReview(String(id), body) : this.api.createReview(body)).subscribe({ next: done, error: fail });
    } else if (t === 'faq') {
      const body = { id: Number(id) || 0, question: v.question ?? '', answer: v.answer ?? '' };
      if (!body.question) return;
      (id ? this.api.updateFaq(Number(id), body) : this.api.createFaq(body)).subscribe({ next: done, error: fail });
    } else if (t === 'gallery') {
      const body = {
        id: String(id ?? ''),
        title: v.title ?? '',
        category: v.category ?? '',
        imageUrl: v.imageUrl ?? '',
        alt: v.alt ?? '',
      };
      (id ? this.api.updateGallery(String(id), body) : this.api.createGallery(body)).subscribe({ next: done, error: fail });
    } else if (t === 'nav') {
      const body = { id: String(id ?? ''), label: v.label ?? '', path: v.path ?? '', visible: !!v.visible };
      (id ? this.api.updateNav(String(id), body) : this.api.createNav(body)).subscribe({ next: done, error: fail });
    }
  }

  async remove(row: Record<string, unknown>): Promise<void> {
    const t = this.type();
    const id = row['id'];
    const confirmed = await this.dialog.confirm({
      title: t === 'bookings' ? 'Delete Booking Record?' : 'Delete Item?',
      message: 'This record will be permanently deleted from the database.',
      confirmText: 'Delete',
      cancelText: 'Cancel',
      variant: 'danger',
      icon: '🗑️',
      badge: 'DELETE',
    });
    if (!confirmed) return;

    if (t === 'bookings') {
      this.api.deleteBooking(String(id)).subscribe(() => this.load());
      return;
    }
    const map: Record<string, () => void> = {
      features: () => this.api.deleteFeature(String(id)).subscribe(() => this.load()),
      press: () => this.api.deletePress(String(id)).subscribe(() => this.load()),
      reviews: () => this.api.deleteReview(String(id)).subscribe(() => this.load()),
      faq: () => this.api.deleteFaq(Number(id)).subscribe(() => this.load()),
      gallery: () => this.api.deleteGallery(String(id)).subscribe(() => this.load()),
      nav: () => this.api.deleteNav(String(id)).subscribe(() => this.load()),
    };
    map[t]?.();
  }

  cancel(): void {
    this.showForm.set(false);
  }

  rowLabel(row: Record<string, unknown>): string {
    const t = this.type();
    if (t === 'bookings') return String(row['name']);
    return String(row['title'] ?? row['name'] ?? row['question'] ?? row['label'] ?? row['id']);
  }
}
