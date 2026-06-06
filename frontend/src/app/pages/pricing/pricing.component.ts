import { AsyncPipe, DecimalPipe } from '@angular/common';
import { Component, inject } from '@angular/core';
import { RouterLink } from '@angular/router';
import { combineLatest, map } from 'rxjs';
import { ApiService } from '../../core/api.service';
import { PageHeroComponent } from '../../shared/page-hero/page-hero.component';
import { ProductImageViewerComponent } from '../../shared/product-image-viewer/product-image-viewer.component';

@Component({
  selector: 'app-pricing',
  standalone: true,
  imports: [AsyncPipe, DecimalPipe, RouterLink, PageHeroComponent, ProductImageViewerComponent],
  templateUrl: './pricing.component.html',
  styleUrl: './pricing.component.scss',
})
export class PricingComponent {
  private readonly api = inject(ApiService);
  vm$ = combineLatest([this.api.getServices(), this.api.getSite()]).pipe(
    map(([services, site]) => ({ services, site }))
  );
}
