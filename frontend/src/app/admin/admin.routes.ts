import { Routes } from '@angular/router';
import { adminAuthGuard, adminGuestGuard } from './guards/admin-auth.guard';
import { AdminLayoutComponent } from './layout/admin-layout.component';

export const adminRoutes: Routes = [
  {
    path: 'login',
    canActivate: [adminGuestGuard],
    loadComponent: () => import('./pages/login/login.component').then((m) => m.LoginComponent),
  },
  {
    path: '',
    canActivate: [adminAuthGuard],
    component: AdminLayoutComponent,
    children: [
      { path: '', loadComponent: () => import('./pages/dashboard/dashboard.component').then((m) => m.DashboardComponent) },
      { path: 'presets', loadComponent: () => import('./pages/presets-admin/presets-admin.component').then((m) => m.PresetsAdminComponent) },
      { path: 'sections', loadComponent: () => import('./pages/sections-admin/sections-admin.component').then((m) => m.SectionsAdminComponent) },
      { path: 'media', loadComponent: () => import('./pages/media-admin/media-admin.component').then((m) => m.MediaAdminComponent) },
      { path: 'orders', loadComponent: () => import('./pages/orders-admin/orders-admin.component').then((m) => m.OrdersAdminComponent) },
      { path: 'site', loadComponent: () => import('./pages/site-editor/site-editor.component').then((m) => m.SiteEditorComponent) },
      { path: 'services', loadComponent: () => import('./pages/services-admin/services-admin.component').then((m) => m.ServicesAdminComponent) },
      { path: 'features', loadComponent: () => import('./pages/content-list/content-list.component').then((m) => m.ContentListComponent), data: { entity: 'features' } },
      { path: 'press', loadComponent: () => import('./pages/content-list/content-list.component').then((m) => m.ContentListComponent), data: { entity: 'press' } },
      { path: 'reviews', loadComponent: () => import('./pages/content-list/content-list.component').then((m) => m.ContentListComponent), data: { entity: 'reviews' } },
      { path: 'faq', loadComponent: () => import('./pages/content-list/content-list.component').then((m) => m.ContentListComponent), data: { entity: 'faq' } },
      { path: 'gallery', loadComponent: () => import('./pages/content-list/content-list.component').then((m) => m.ContentListComponent), data: { entity: 'gallery' } },
      { path: 'nav', loadComponent: () => import('./pages/content-list/content-list.component').then((m) => m.ContentListComponent), data: { entity: 'nav' } },
      { path: 'bookings', loadComponent: () => import('./pages/orders-admin/orders-admin.component').then((m) => m.OrdersAdminComponent) },
      { path: 'payments', loadComponent: () => import('./pages/payments-admin/payments-admin.component').then((m) => m.PaymentsAdminComponent) },
    ],
  },
];
