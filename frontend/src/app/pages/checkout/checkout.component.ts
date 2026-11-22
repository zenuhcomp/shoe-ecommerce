import { AsyncPipe, DecimalPipe, CommonModule, TitleCasePipe } from '@angular/common';
import { Component, inject, signal, OnInit, OnDestroy } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { combineLatest, map, switchMap, tap } from 'rxjs';
import { ApiService } from '../../core/api.service';
import { PaymentMethodInfo, PaymentInitiateResponse } from '../../core/models';
import { PageHeroComponent } from '../../shared/page-hero/page-hero.component';
import { ProductImageViewerComponent } from '../../shared/product-image-viewer/product-image-viewer.component';

@Component({
  selector: 'app-checkout',
  standalone: true,
  imports: [
    AsyncPipe,
    CommonModule,
    DecimalPipe,
    TitleCasePipe,
    ReactiveFormsModule,
    RouterLink,
    PageHeroComponent,
    ProductImageViewerComponent,
  ],
  templateUrl: './checkout.component.html',
  styleUrl: './checkout.component.scss',
})
export class CheckoutComponent implements OnInit, OnDestroy {
  private readonly api = inject(ApiService);
  private readonly route = inject(ActivatedRoute);
  private readonly fb = inject(FormBuilder);

  submitting = signal(false);
  paymentPending = signal(false);
  successMessage = signal('');
  errorMessage = signal('');
  orderId = signal('');
  stkCheckoutId = signal('');
  verifyReference = signal('');

  paymentMethods = signal<PaymentMethodInfo[]>([]);
  selectedMethodId = signal<string>('');
  offlineChannel = signal<'till' | 'paybill'>('till');

  form = this.fb.nonNullable.group({
    name: ['', Validators.required],
    phone: ['', Validators.required],
    email: [''],
    size: ['', Validators.required],
    quantity: [1, [Validators.required, Validators.min(1)]],
    mpesaPhone: [''],
    paymentReference: [''],
    notes: [''],
  });

  vm$ = this.route.queryParamMap.pipe(
    map((params) => params.get('service') ?? ''),
    switchMap((slug) =>
      combineLatest([
        slug ? this.api.getService(slug) : this.api.getServices().pipe(map((s) => s[0])),
        this.api.getSite(),
      ]).pipe(map(([service, site]) => ({ service, site, slug })))
    )
  );

  ngOnInit(): void {
    this.api.getActivePaymentMethods().subscribe({
      next: ({ methods }) => {
        const enabled = methods.filter((m) => m.enabled);
        this.paymentMethods.set(enabled);
        if (enabled.length > 0) {
          this.selectedMethodId.set(enabled[0].id);
        }
      },
      error: () => {
        // Fallback: show offline M-Pesa as default
        this.paymentMethods.set([]);
      },
    });
  }

  get selectedMethod(): PaymentMethodInfo | undefined {
    return this.paymentMethods().find((m) => m.id === this.selectedMethodId());
  }

  selectMethod(id: string): void {
    this.selectedMethodId.set(id);
    this.errorMessage.set('');
    this.paymentPending.set(false);
    this.stkCheckoutId.set('');
  }

  setOfflineChannel(ch: 'till' | 'paybill'): void {
    this.offlineChannel.set(ch);
  }

  total(price: number): number {
    return price * (this.form.controls.quantity.value || 1);
  }

  methodIcon(type: string): string {
    const icons: Record<string, string> = {
      online: '📱', manual: '🏷️', card: '💳', cash: '💵', bank: '🏦',
    };
    return icons[type] ?? '💰';
  }

  submit(slug: string, totalKes: number): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const methodId = this.selectedMethodId();
    if (!methodId) {
      this.errorMessage.set('Please select a payment method.');
      return;
    }

    const v = this.form.getRawValue();

    // Validate offline reference
    if (methodId === 'mpesa_manual' && !v.paymentReference.trim()) {
      this.errorMessage.set('Enter your M-Pesa confirmation code.');
      return;
    }

    // Validate M-Pesa phone for STK push
    const phone = v.mpesaPhone || v.phone;
    if (methodId === 'mpesa_express' && !phone.trim()) {
      this.errorMessage.set('Enter a valid M-Pesa phone number for STK push.');
      return;
    }

    this.submitting.set(true);
    this.errorMessage.set('');
    this.successMessage.set('');

    // First create the order
    this.api.submitOrder({
      name: v.name,
      phone: v.phone,
      email: v.email || undefined,
      serviceSlug: slug,
      size: v.size,
      quantity: Number(v.quantity),
      paymentMethod: methodId as any,
      mpesaPhone: phone,
      paymentReference: v.paymentReference || undefined,
      notes: v.notes || undefined,
    }).subscribe({
      next: (orderRes) => {
        this.orderId.set(orderRes.id);

        // Then initiate payment strategy
        this.api.initiatePayment({
          orderId: orderRes.id,
          serviceSlug: slug,
          customerName: v.name,
          customerPhone: v.phone,
          customerEmail: v.email || undefined,
          amountKes: totalKes,
          method: methodId,
          mpesaPhone: phone,
          notes: v.notes || undefined,
        }).subscribe({
          next: (payRes: PaymentInitiateResponse) => {
            this.submitting.set(false);

            if (payRes.promptRequired) {
              // STK Push — wait for PIN
              this.paymentPending.set(true);
              this.stkCheckoutId.set(payRes.transactionId ?? '');
              this.successMessage.set(payRes.message);
              this.startPolling(orderRes.id);
            } else {
              this.successMessage.set(payRes.message || 'Order placed successfully!');
              if (orderRes.whatsAppUrl) {
                window.open(orderRes.whatsAppUrl, '_blank', 'noopener');
              }
            }
          },
          error: (err) => {
            this.submitting.set(false);
            this.errorMessage.set(err?.error?.error ?? 'Payment initiation failed. Your order was saved — contact us to complete payment.');
          },
        });
      },
      error: (err) => {
        this.submitting.set(false);
        this.errorMessage.set(err?.error?.error ?? 'Could not place order. Please try again.');
      },
    });
  }

  verifyPayment(): void {
    const ref = this.verifyReference();
    if (!ref.trim()) return;

    this.submitting.set(true);
    this.api.verifyPayment({
      orderId: this.orderId(),
      method: this.selectedMethodId(),
      reference: ref,
    }).subscribe({
      next: (res) => {
        this.submitting.set(false);
        this.paymentPending.set(false);
        this.stopPolling();
        this.successMessage.set(res.message || 'Payment verified successfully!');
      },
      error: (err) => {
        this.submitting.set(false);
        this.errorMessage.set(err?.error?.error ?? 'Verification failed. Please try again.');
      },
    });
  }

  private pollInterval?: any;

  private startPolling(orderId: string): void {
    this.stopPolling();
    let attempts = 0;
    this.pollInterval = setInterval(() => {
      attempts++;
      if (attempts > 30) {
        this.stopPolling();
        return;
      }
      this.api.getOrderStatus(orderId).subscribe({
        next: (statusRes) => {
          if (statusRes.status === 'paid') {
            this.stopPolling();
            this.paymentPending.set(false);
            const ref = statusRes.paymentReference ? ` (M-Pesa Receipt: ${statusRes.paymentReference})` : '';
            this.successMessage.set(`✓ Payment Confirmed${ref}! Your order is processed.`);
          } else if (statusRes.status === 'payment_failed') {
            this.stopPolling();
            this.paymentPending.set(false);
            this.errorMessage.set('Payment was cancelled or failed. Please try again.');
          }
        },
        error: () => {},
      });
    }, 2000);
  }

  private stopPolling(): void {
    if (this.pollInterval) {
      clearInterval(this.pollInterval);
      this.pollInterval = undefined;
    }
  }

  simulateMpesaSuccess(totalKes: number): void {
    const oid = this.orderId();
    if (!oid) return;
    this.submitting.set(true);
    this.api.simulateMpesaCallback({
      orderId: oid,
      checkoutRequestId: this.stkCheckoutId() || oid,
      resultCode: 0,
      resultDesc: 'The service request is processed successfully.',
      amount: totalKes,
      phone: this.form.getRawValue().phone,
    }).subscribe({
      next: (res) => {
        this.submitting.set(false);
        this.stopPolling();
        this.paymentPending.set(false);
        const receipt = res?.order?.paymentReference || 'MPE' + Date.now().toString().slice(-6);
        this.successMessage.set(`✓ M-Pesa Callback Reconciled (Receipt: ${receipt})! Order ${oid} marked as PAID.`);
      },
      error: (err) => {
        this.submitting.set(false);
        this.errorMessage.set(err?.error?.error ?? 'Simulation failed.');
      },
    });
  }

  ngOnDestroy(): void {
    this.stopPolling();
  }
}
