import { Component, Input, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { ApiService } from '../../core/api.service';
import { LashService } from '../../core/models';

@Component({
  selector: 'app-booking-form',
  standalone: true,
  imports: [ReactiveFormsModule],
  templateUrl: './booking-form.component.html',
  styleUrl: './booking-form.component.scss',
})
export class BookingFormComponent {
  @Input() services: LashService[] = [];

  private readonly fb = inject(FormBuilder);
  private readonly api = inject(ApiService);

  submitted = signal(false);
  successMessage = signal('');
  errorMessage = signal('');

  form = this.fb.nonNullable.group({
    name: ['', Validators.required],
    phone: ['', Validators.required],
    email: [''],
    serviceSlug: [''],
    preferredAt: [''],
    message: [''],
  });

  onSubmit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }
    this.errorMessage.set('');
    this.api.submitBooking(this.form.getRawValue()).subscribe({
      next: (res) => {
        this.submitted.set(true);
        this.successMessage.set(res.message);
        this.form.reset();
      },
      error: () => this.errorMessage.set('Something went wrong. Please try WhatsApp instead.'),
    });
  }
}
