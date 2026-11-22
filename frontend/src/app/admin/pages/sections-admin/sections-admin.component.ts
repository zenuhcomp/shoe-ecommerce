import { Component, OnInit, inject, signal } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { AdminApiService } from '../../services/admin-api.service';
import { ConfirmDialogService } from '../../../shared/confirm-dialog/confirm-dialog.service';
import { Section } from '../../../core/models';

@Component({
  selector: 'app-sections-admin',
  standalone: true,
  imports: [ReactiveFormsModule],
  templateUrl: './sections-admin.component.html',
  styleUrl: './sections-admin.component.scss',
})
export class SectionsAdminComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly api = inject(AdminApiService);
  private readonly dialog = inject(ConfirmDialogService);

  sections = signal<Section[]>([]);
  editingSection = signal<Section | null>(null);
  isCreating = signal(false);
  loading = signal(false);
  saving = signal(false);
  message = signal<{ type: 'success' | 'error'; text: string } | null>(null);

  editForm = this.fb.group({
    id: ['', Validators.required],
    type: ['custom_block', Validators.required],
    enabled: [true],
    title: [''],
    subtitle: [''],
    badge: [''],
    content: [''],
    backgroundType: ['default'],
    backgroundColor: ['#ffffff'],
    backgroundImage: [''],
    backgroundVideo: [''],
    overlayOpacity: [40],
    animation: ['slide-up'],
    animationDuration: ['0.8s'],
    animationDelay: ['0.1s'],
    hoverEffect: ['lift'],
    primaryCtaText: [''],
    primaryCtaUrl: [''],
    secondaryCtaText: [''],
    secondaryCtaUrl: [''],
  });

  ngOnInit(): void {
    this.loadSections();
  }

  loadSections(): void {
    this.loading.set(true);
    this.api.listSections().subscribe({
      next: (list) => {
        this.sections.set(list);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
        this.message.set({ type: 'error', text: 'Could not load sections' });
      },
    });
  }

  toggleEnabled(sec: Section): void {
    const updated = { ...sec, enabled: !sec.enabled };
    this.api.updateSection(sec.id, updated).subscribe({
      next: () => {
        this.sections.update((list) =>
          list.map((item) => (item.id === sec.id ? updated : item))
        );
      },
      error: () => this.message.set({ type: 'error', text: 'Could not toggle section' }),
    });
  }

  moveUp(index: number): void {
    if (index <= 0) return;
    const list = [...this.sections()];
    const temp = list[index - 1];
    list[index - 1] = list[index];
    list[index] = temp;
    this.sections.set(list);
    this.api.reorderSections(list.map((s) => s.id)).subscribe();
  }

  moveDown(index: number): void {
    const list = [...this.sections()];
    if (index >= list.length - 1) return;
    const temp = list[index + 1];
    list[index + 1] = list[index];
    list[index] = temp;
    this.sections.set(list);
    this.api.reorderSections(list.map((s) => s.id)).subscribe();
  }

  startEdit(sec: Section): void {
    this.isCreating.set(false);
    this.editingSection.set(sec);
    this.editForm.patchValue({
      id: sec.id,
      type: sec.type,
      enabled: sec.enabled,
      title: sec.title,
      subtitle: sec.subtitle,
      badge: sec.badge,
      content: sec.content,
      backgroundType: sec.backgroundType,
      backgroundColor: sec.backgroundColor || '#ffffff',
      backgroundImage: sec.backgroundImage,
      backgroundVideo: sec.backgroundVideo,
      overlayOpacity: sec.overlayOpacity || 40,
      animation: sec.animation || 'slide-up',
      animationDuration: sec.animationDuration || '0.8s',
      animationDelay: sec.animationDelay || '0.1s',
      hoverEffect: sec.hoverEffect || 'lift',
      primaryCtaText: sec.primaryCtaText,
      primaryCtaUrl: sec.primaryCtaUrl,
      secondaryCtaText: sec.secondaryCtaText,
      secondaryCtaUrl: sec.secondaryCtaUrl,
    });
  }

  startCreate(): void {
    this.isCreating.set(true);
    this.editingSection.set(null);
    this.editForm.reset({
      id: 'sec-' + Math.random().toString(36).substring(2, 8),
      type: 'custom_block',
      enabled: true,
      title: 'New Section Title',
      subtitle: 'Optional subtitle describing this section',
      badge: '~ Feature ~',
      content: '',
      backgroundType: 'default',
      backgroundColor: '#ffffff',
      backgroundImage: '',
      backgroundVideo: '',
      overlayOpacity: 40,
      animation: 'slide-up',
      animationDuration: '0.8s',
      animationDelay: '0.1s',
      hoverEffect: 'lift',
      primaryCtaText: 'Learn More',
      primaryCtaUrl: '/collections',
      secondaryCtaText: '',
      secondaryCtaUrl: '',
    });
  }

  cancelEdit(): void {
    this.editingSection.set(null);
    this.isCreating.set(false);
  }

  saveSection(): void {
    if (this.editForm.invalid) return;
    this.saving.set(true);
    const formVal = this.editForm.getRawValue() as Section;

    if (this.isCreating()) {
      this.api.createSection(formVal).subscribe({
        next: (created) => {
          this.sections.update((list) => [...list, created]);
          this.saving.set(false);
          this.cancelEdit();
          this.message.set({ type: 'success', text: 'Section created successfully!' });
        },
        error: () => {
          this.saving.set(false);
          this.message.set({ type: 'error', text: 'Could not create section.' });
        },
      });
    } else {
      this.api.updateSection(formVal.id, formVal).subscribe({
        next: (updated) => {
          this.sections.update((list) =>
            list.map((item) => (item.id === updated.id ? updated : item))
          );
          this.saving.set(false);
          this.cancelEdit();
          this.message.set({ type: 'success', text: 'Section updated successfully!' });
        },
        error: () => {
          this.saving.set(false);
          this.message.set({ type: 'error', text: 'Could not update section.' });
        },
      });
    }
  }

  async deleteSection(id: string): Promise<void> {
    const confirmed = await this.dialog.confirm({
      title: 'Remove Section?',
      message: 'This section will be permanently deleted from the homepage layout.',
      confirmText: 'Delete Section',
      cancelText: 'Keep Section',
      variant: 'danger',
      icon: '🗑️',
      badge: 'HOMEPAGE BUILDER',
    });
    if (!confirmed) return;
    this.api.deleteSection(id).subscribe({
      next: () => {
        this.sections.update((list) => list.filter((s) => s.id !== id));
        this.message.set({ type: 'success', text: 'Section removed.' });
      },
      error: () => this.message.set({ type: 'error', text: 'Could not delete section.' }),
    });
  }
}
