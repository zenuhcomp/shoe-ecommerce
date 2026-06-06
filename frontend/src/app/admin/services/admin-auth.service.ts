import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { tap } from 'rxjs';

export interface AdminUser {
  id: string;
  email: string;
  displayName: string;
  role: string;
}

interface LoginResponse {
  token: string;
  user: AdminUser;
}

@Injectable({ providedIn: 'root' })
export class AdminAuthService {
  private readonly http = inject(HttpClient);
  private readonly router = inject(Router);
  private readonly tokenKey = 'luna_cms_token';
  private readonly userKey = 'luna_cms_user';

  readonly user = signal<AdminUser | null>(this.loadUser());

  isLoggedIn(): boolean {
    return !!sessionStorage.getItem(this.tokenKey);
  }

  token(): string | null {
    return sessionStorage.getItem(this.tokenKey);
  }

  login(email: string, password: string) {
    return this.http.post<LoginResponse>('/api/admin/auth/login', { email, password }).pipe(
      tap((res) => {
        sessionStorage.setItem(this.tokenKey, res.token);
        sessionStorage.setItem(this.userKey, JSON.stringify(res.user));
        this.user.set(res.user);
      })
    );
  }

  logout(): void {
    sessionStorage.removeItem(this.tokenKey);
    sessionStorage.removeItem(this.userKey);
    this.user.set(null);
    this.router.navigate(['/admin/login']);
  }

  fetchMe() {
    return this.http.get<AdminUser>('/api/admin/auth/me').pipe(
      tap((u) => {
        sessionStorage.setItem(this.userKey, JSON.stringify(u));
        this.user.set(u);
      })
    );
  }

  private loadUser(): AdminUser | null {
    const raw = sessionStorage.getItem(this.userKey);
    if (!raw) return null;
    try {
      return JSON.parse(raw) as AdminUser;
    } catch {
      return null;
    }
  }
}
