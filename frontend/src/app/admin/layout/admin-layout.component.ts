import { Component, inject } from '@angular/core';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { AdminAuthService } from '../services/admin-auth.service';

@Component({
  selector: 'app-admin-layout',
  standalone: true,
  imports: [RouterOutlet, RouterLink, RouterLinkActive],
  templateUrl: './admin-layout.component.html',
  styleUrl: './admin-layout.component.scss',
})
export class AdminLayoutComponent {
  readonly auth = inject(AdminAuthService);

  readonly nav = [
    { path: '/admin', label: 'Dashboard', exact: true },
    { path: '/admin/site', label: 'Site settings' },
    { path: '/admin/services', label: 'Shoe styles' },
    { path: '/admin/features', label: 'Features' },
    { path: '/admin/press', label: 'Press' },
    { path: '/admin/reviews', label: 'Reviews' },
    { path: '/admin/faq', label: 'FAQ' },
    { path: '/admin/gallery', label: 'Gallery' },
    { path: '/admin/nav', label: 'Navigation' },
    { path: '/admin/bookings', label: 'Bookings' },
  ];

  logout(): void {
    this.auth.logout();
  }
}
