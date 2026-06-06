import { Routes } from '@angular/router';
import { ShellComponent } from './layout/shell/shell.component';
import { HomeComponent } from './pages/home/home.component';
import { CollectionsComponent } from './pages/collections/collections.component';
import { PricingComponent } from './pages/pricing/pricing.component';
import { GalleryComponent } from './pages/gallery/gallery.component';
import { AboutComponent } from './pages/about/about.component';
import { HouseCallsComponent } from './pages/house-calls/house-calls.component';
import { ReviewsComponent } from './pages/reviews/reviews.component';
import { ContactComponent } from './pages/contact/contact.component';
import { CheckoutComponent } from './pages/checkout/checkout.component';
import { AftercareComponent } from './pages/aftercare/aftercare.component';

export const routes: Routes = [
  {
    path: 'admin',
    loadChildren: () => import('./admin/admin.routes').then((m) => m.adminRoutes),
  },
  {
    path: '',
    component: ShellComponent,
    children: [
      { path: '', component: HomeComponent },
      { path: 'collections', component: CollectionsComponent },
      { path: 'lash-styles', redirectTo: 'collections', pathMatch: 'full' },
      { path: 'pricing', component: PricingComponent },
      { path: 'checkout', component: CheckoutComponent },
      { path: 'gallery', component: GalleryComponent },
      { path: 'about', component: AboutComponent },
      { path: 'delivery', component: HouseCallsComponent },
      { path: 'house-calls', redirectTo: 'delivery', pathMatch: 'full' },
      { path: 'reviews', component: ReviewsComponent },
      { path: 'contact', component: ContactComponent },
      { path: 'shoe-care', component: AftercareComponent },
      { path: 'aftercare', redirectTo: 'shoe-care', pathMatch: 'full' },
      { path: '**', redirectTo: '' },
    ],
  },
];
