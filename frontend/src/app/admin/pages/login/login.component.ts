import { Component, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { AdminAuthService } from '../../services/admin-auth.service';

@Component({
  selector: 'app-admin-login',
  standalone: true,
  imports: [ReactiveFormsModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss',
})
export class LoginComponent {
  private readonly fb = inject(FormBuilder);
  private readonly auth = inject(AdminAuthService);
  private readonly router = inject(Router);

  error = signal('');
  loading = signal(false);
  showPassword = signal(false);

  form = this.fb.nonNullable.group({
    email: ['', [Validators.required, Validators.email]],
    password: ['', Validators.required],
  });

  toggleShowPassword(): void {
    this.showPassword.update((visible) => !visible);
  }

  submit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.loading.set(true);
    this.error.set('');

    const raw = this.form.getRawValue();
    const email = raw.email.trim();
    const password = raw.password.trim();

    this.auth.login(email, password).subscribe({
      next: () => {
        this.loading.set(false);
        this.router.navigate(['/admin']);
      },
      error: (err) => {
        this.loading.set(false);
        const serverMsg = err?.error?.error;
        if (err.status === 401) {
          this.error.set(serverMsg || 'Invalid email or password. Please verify your credentials.');
        } else if (err.status === 0) {
          this.error.set('Cannot connect to backend server. Please verify backend is running on port 8084.');
        } else {
          this.error.set(serverMsg || `Login failed (${err.status || 'server error'})`);
        }
      },
    });
  }
}
