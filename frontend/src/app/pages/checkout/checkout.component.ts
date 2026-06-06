import { AsyncPipe, DecimalPipe } from '@angular/common';
import { Component, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { combineLatest, map, switchMap } from 'rxjs';
import { ApiService } from '../../core/api.service';
import { PaymentMethod } from '../../core/models';
import { PageHeroComponent } from '../../shared/page-hero/page-hero.component';
import { ProductImageViewerComponent } from '../../shared/product-image-viewer/product-image-viewer.component';

@Component({
  selector: 'app-checkout',
  standalone: true,
  imports: [
    AsyncPipe,
    DecimalPipe,
    ReactiveFormsModule,
    RouterLink,
    PageHeroComponent,
    ProductImageViewerComponent,
  ],
  templateUrl: './checkout.component.html',
  styleUrl: './checkout.component.scss',
})
export class CheckoutComponent {
  private readonly api = inject(ApiService);
  private readonly route = inject(ActivatedRoute);
  private readonly fb = inject(FormBuilder);

  submitting = signal(false);
  successMessage = signal('');
  errorMessage = signal('');
  orderId = signal('');

  paymentMethod = signal<PaymentMethod>('mpesa_stk');
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

  setPayment(method: PaymentMethod): void {
    this.paymentMethod.set(method);
    this.errorMessage.set('');
  }

  setOfflineChannel(channel: 'till' | 'paybill'): void {
    this.offlineChannel.set(channel);
    this.paymentMethod.set(channel === 'till' ? 'offline_till' : 'offline_paybill');
  }

  total(price: number): number {
    return price * (this.form.controls.quantity.value || 1);
  }

  submit(slug: string): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const method = this.paymentMethod();
    const v = this.form.getRawValue();

    if ((method === 'offline_till' || method === 'offline_paybill') && !v.paymentReference.trim()) {
      this.errorMessage.set('Enter your M-Pesa confirmation code or reference.');
      return;
    }

    this.submitting.set(true);
    this.errorMessage.set('');
    this.successMessage.set('');

    this.api
      .submitOrder({
        name: v.name,
        phone: v.phone,
        email: v.email || undefined,
        serviceSlug: slug,
        size: v.size,
        quantity: Number(v.quantity),
        paymentMethod: method,
        mpesaPhone: v.mpesaPhone || v.phone,
        offlineChannel: this.offlineChannel(),
        paymentReference: v.paymentReference || undefined,
        notes: v.notes || undefined,
      })
      .subscribe({
        next: (res) => {
          this.submitting.set(false);
          this.orderId.set(res.id);
          this.successMessage.set(res.message);
          if (res.whatsAppUrl) {
            window.open(res.whatsAppUrl, '_blank', 'noopener');
          }
        },
        error: (err) => {
          this.submitting.set(false);
          this.errorMessage.set(err?.error?.error ?? 'Could not place order. Please try again.');
        },
      });
  }
}
