import { Component, inject, signal } from '@angular/core';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { AdminAuthService } from '../services/admin-auth.service';

export interface NavGroup {
  label: string;
  items: NavItem[];
}

export interface NavItem {
  path: string;
  label: string;
  icon: string;
  exact?: boolean;
}

@Component({
  selector: 'app-admin-layout',
  standalone: true,
  imports: [RouterOutlet, RouterLink, RouterLinkActive],
  templateUrl: './admin-layout.component.html',
  styleUrl: './admin-layout.component.scss',
})
export class AdminLayoutComponent {
  readonly auth = inject(AdminAuthService);
  sidebarOpen = signal(true);

  readonly navGroups: NavGroup[] = [
    {
      label: 'Overview',
      items: [
        { path: '/admin', icon: '◈', label: 'Dashboard', exact: true },
        { path: '/admin/presets', icon: '◇', label: 'Industry Presets' },
      ],
    },
    {
      label: 'Content',
      items: [
        { path: '/admin/sections', icon: '⊞', label: 'Sections Builder' },
        { path: '/admin/site', icon: '◎', label: 'Branding & Theme' },
        { path: '/admin/services', icon: '▦', label: 'Products & Catalog' },
        { path: '/admin/media', icon: '▣', label: 'Media Hub' },
      ],
    },
    {
      label: 'Community',
      items: [
        { path: '/admin/reviews', icon: '★', label: 'Customer Reviews' },
        { path: '/admin/features', icon: '◆', label: 'Features & Guarantees' },
        { path: '/admin/faq', icon: '?', label: 'FAQ Accordion' },
        { path: '/admin/gallery', icon: '▤', label: 'Gallery Lookbook' },
      ],
    },
    {
      label: 'Business',
      items: [
        { path: '/admin/orders', icon: '◉', label: 'Orders & Inquiries' },
        { path: '/admin/payments', icon: '💳', label: 'Payment Gateways' },
        { path: '/admin/press', icon: '◫', label: 'Press Mentions' },
        { path: '/admin/nav', icon: '≡', label: 'Navigation Links' },
      ],
    },
  ];

  toggleSidebar(): void {
    this.sidebarOpen.update((v) => !v);
  }

  logout(): void {
    this.auth.logout();
  }
}
