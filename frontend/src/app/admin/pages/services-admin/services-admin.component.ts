import { Component, OnInit, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { LashService } from '../../../core/models';
import { AdminApiService } from '../../services/admin-api.service';

@Component({
  selector: 'app-services-admin',
  standalone: true,
  imports: [ReactiveFormsModule],
  templateUrl: './services-admin.component.html',
  styleUrl: './services-admin.component.scss',
})
export class ServicesAdminComponent implements OnInit {
  private readonly api = inject(AdminApiService);
  private readonly fb = inject(FormBuilder);

  items = signal<LashService[]>([]);
  editing = signal<LashService | null>(null);
  showForm = signal(false);
  message = signal('');

  form = this.fb.nonNullable.group({
    name: ['', Validators.required],
    slug: [''],
    emoji: [''],
    style: [''],
    length: [''],
    description: [''],
    priceKes: [0, Validators.required],
    featured: [false],
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
    this.form.reset({ name: '', slug: '', emoji: '✨', style: '', length: '', description: '', priceKes: 450, featured: false });
  }

  startEdit(s: LashService): void {
    this.editing.set(s);
    this.showForm.set(true);
    this.form.patchValue(s);
  }

  save(): void {
    if (this.form.invalid) return;
    const v = this.form.getRawValue();
    const payload: LashService = {
      id: this.editing()?.id ?? '',
      name: v.name,
      slug: v.slug || v.name.toLowerCase().replace(/\s+/g, '-'),
      emoji: v.emoji,
      style: v.style,
      length: v.length,
      description: v.description,
      priceKes: Number(v.priceKes),
      featured: v.featured,
    };
    const req = this.editing()
      ? this.api.updateService(this.editing()!.id, payload)
      : this.api.createService(payload);
    req.subscribe({
      next: () => {
        this.message.set('Saved.');
        this.editing.set(null);
        this.showForm.set(false);
        this.load();
      },
      error: () => this.message.set('Save failed.'),
    });
  }

  remove(s: LashService): void {
    if (!confirm(`Delete "${s.name}"?`)) return;
    this.api.deleteService(s.id).subscribe(() => this.load());
  }

  cancel(): void {
    this.editing.set(null);
    this.showForm.set(false);
  }
}
