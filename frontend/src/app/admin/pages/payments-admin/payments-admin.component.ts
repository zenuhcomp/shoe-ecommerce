import { Component, inject, signal, OnInit } from '@angular/core';
import { FormBuilder, ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { AdminApiService } from '../../services/admin-api.service';
import { ConfirmDialogService } from '../../../shared/confirm-dialog/confirm-dialog.service';
import { ConfirmDialogComponent } from '../../../shared/confirm-dialog/confirm-dialog.component';
import { PaymentConfig, PaymentMethodInfo } from '../../../core/models';

interface PaymentSection {
  id: keyof PaymentConfig;
  label: string;
  icon: string;
  description: string;
  variant: 'primary' | 'accent' | 'danger' | 'warning';
}

@Component({
  selector: 'app-payments-admin',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, ConfirmDialogComponent],
  templateUrl: './payments-admin.component.html',
  styleUrl: './payments-admin.component.scss',
})
export class PaymentsAdminComponent implements OnInit {
  private readonly adminApi = inject(AdminApiService);
  private readonly fb = inject(FormBuilder);
  readonly confirmDialog = inject(ConfirmDialogService);

  loading = signal(true);
  saving = signal(false);
  toast = signal<{ type: 'success' | 'error'; message: string } | null>(null);
  methods = signal<PaymentMethodInfo[]>([]);

  form = this.fb.nonNullable.group({
    // M-Pesa Express (STK Push)
    mpesaExpressEnabled: [false],
    consumerKey: [''],
    consumerSecret: [''],
    passKey: [''],
    shortCode: [''],
    partyB: [''],
    accountReference: [''],
    environment: ['sandbox'],
    callbackUrl: [''],
    // M-Pesa Offline
    mpesaManualEnabled: [false],
    tillNumber: [''],
    paybillNumber: [''],
    paybillAccount: [''],
    mpesaShortCode: [''],
    mpesaInstructions: [''],
    // Card
    cardEnabled: [false],
    cardGateway: [''],
    cardPublicKey: [''],
    cardSecretKey: [''],
    cardInstructions: [''],
    // Cash
    cashEnabled: [false],
    cashInstructions: [''],
    // Bank
    bankEnabled: [false],
    bankName: [''],
    bankAccountName: [''],
    bankAccountNo: [''],
    bankBranch: [''],
    bankSwiftCode: [''],
    bankInstructions: [''],
  });

  // Password field visibility
  showConsumerSecret = signal(false);
  showPassKey = signal(false);
  showCardSecret = signal(false);

  ngOnInit(): void {
    this.loadConfig();
  }

  loadConfig(): void {
    this.loading.set(true);
    this.adminApi.getPaymentConfig().subscribe({
      next: ({ config, methods }) => {
        this.form.patchValue(config as any);
        this.methods.set(methods);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
        this.showToast('error', 'Failed to load payment configuration.');
      },
    });
  }

  async save(): Promise<void> {
    const confirmed = await this.confirmDialog.confirm({
      title: 'Save Payment Configuration',
      message: 'This will immediately update which payment methods are available to customers at checkout.',
      confirmText: 'Save Changes',
      cancelText: 'Cancel',
      variant: 'accent',
      icon: '💳',
      badge: 'Payment Gateway',
    });

    if (!confirmed) return;

    this.saving.set(true);
    const cfg = this.form.getRawValue() as unknown as PaymentConfig;

    this.adminApi.updatePaymentConfig(cfg).subscribe({
      next: ({ message, methods }) => {
        this.saving.set(false);
        this.methods.set(methods);
        this.showToast('success', message || 'Payment settings saved successfully.');
      },
      error: (err) => {
        this.saving.set(false);
        this.showToast('error', err?.error?.error ?? 'Failed to save payment settings.');
      },
    });
  }

  private showToast(type: 'success' | 'error', message: string): void {
    this.toast.set({ type, message });
    setTimeout(() => this.toast.set(null), 4500);
  }

  isMethodEnabled(id: string): boolean {
    return this.methods().find(m => m.id === id)?.enabled ?? false;
  }

  toggleConsumerSecret(): void { this.showConsumerSecret.update(v => !v); }
  togglePassKey(): void { this.showPassKey.update(v => !v); }
  toggleCardSecret(): void { this.showCardSecret.update(v => !v); }

  simOrderId = signal('');
  simulating = signal(false);
  simResult = signal('');

  simulateCallback(): void {
    const oid = this.simOrderId().trim();
    if (!oid) {
      this.showToast('error', 'Enter an Order ID to simulate reconciliation.');
      return;
    }
    this.simulating.set(true);
    this.simResult.set('');
    this.adminApi.simulateCallback({ orderId: oid }).subscribe({
      next: (res) => {
        this.simulating.set(false);
        const receipt = res?.order?.paymentReference || 'CONFIRMED';
        this.simResult.set(`✓ Order ${oid} successfully reconciled! Status: PAID (Receipt: ${receipt})`);
        this.showToast('success', `Order ${oid} reconciled to PAID with receipt ${receipt}`);
      },
      error: (err) => {
        this.simulating.set(false);
        this.showToast('error', err?.error?.error ?? 'Simulation failed. Verify Order ID exists.');
      },
    });
  }
}
