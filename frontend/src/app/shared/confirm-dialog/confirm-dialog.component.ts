import { Component, HostListener, inject } from '@angular/core';
import { ConfirmDialogService } from './confirm-dialog.service';

@Component({
  selector: 'app-confirm-dialog',
  standalone: true,
  template: `
    @if (dialog.current(); as item) {
      <div class="dialog-backdrop" (click)="onBackdropClick($event)">
        <div
          class="dialog-card"
          [class]="'variant-' + (item.options.variant || 'accent')"
          role="dialog"
          aria-modal="true"
          [attr.aria-label]="item.options.title"
          (click)="$event.stopPropagation()"
        >
          <div class="card-top">
            @if (item.options.icon) {
              <div class="dialog-icon-bubble">{{ item.options.icon }}</div>
            }
            @if (item.options.badge) {
              <span class="dialog-badge">{{ item.options.badge }}</span>
            }
          </div>

          <h2 class="dialog-title">{{ item.options.title }}</h2>

          <p class="dialog-message">{{ item.options.message }}</p>

          <div class="dialog-actions">
            <button
              type="button"
              class="btn-cancel"
              (click)="dialog.dismiss()"
            >
              {{ item.options.cancelText || 'Cancel' }}
            </button>
            <button
              type="button"
              class="btn-confirm"
              (click)="dialog.accept()"
            >
              {{ item.options.confirmText || 'Confirm' }}
            </button>
          </div>
        </div>
      </div>
    }
  `,
  styleUrl: './confirm-dialog.component.scss',
})
export class ConfirmDialogComponent {
  readonly dialog = inject(ConfirmDialogService);

  @HostListener('document:keydown.escape')
  onEscape(): void {
    this.dialog.dismiss();
  }

  onBackdropClick(event: MouseEvent): void {
    if (event.target === event.currentTarget) {
      this.dialog.dismiss();
    }
  }
}
