import { Injectable, signal } from '@angular/core';

export interface ConfirmDialogOptions {
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  variant?: 'primary' | 'accent' | 'danger' | 'warning';
  icon?: string;
  badge?: string;
}

export interface DialogInstance {
  options: ConfirmDialogOptions;
  resolve: (value: boolean) => void;
}

@Injectable({ providedIn: 'root' })
export class ConfirmDialogService {
  readonly current = signal<DialogInstance | null>(null);

  confirm(options: ConfirmDialogOptions): Promise<boolean> {
    return new Promise<boolean>((resolve) => {
      this.current.set({ options, resolve });
    });
  }

  accept(): void {
    const dialog = this.current();
    if (dialog) {
      dialog.resolve(true);
      this.current.set(null);
    }
  }

  dismiss(): void {
    const dialog = this.current();
    if (dialog) {
      dialog.resolve(false);
      this.current.set(null);
    }
  }
}
