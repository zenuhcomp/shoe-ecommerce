import { AsyncPipe } from '@angular/common';
import { Component, inject } from '@angular/core';
import { combineLatest, map } from 'rxjs';
import { ApiService } from '../../core/api.service';
import { PageHeroComponent } from '../../shared/page-hero/page-hero.component';
import { ServiceGridComponent } from '../../shared/service-grid/service-grid.component';

@Component({
  selector: 'app-collections',
  standalone: true,
  imports: [AsyncPipe, PageHeroComponent, ServiceGridComponent],
  template: `
    @if (vm$ | async; as vm) {
      <app-page-hero
        [eyebrow]="'◆ ' + (vm.site.brandName || 'Store') + ' Catalog'"
        [title]="(vm.site.brandName || 'Curated') + ' <span class=\\'accent\\'>Collection</span>'"
        [subtitle]="vm.site.tagline || 'Explore verified specifications, 360° interactive views, and immediate ordering'"
      />
      <section class="page-section">
        <app-service-grid [services]="vm.services" [site]="vm.site" />
      </section>
    }
  `,
  styles: [
    `
      .page-section {
        max-width: 1200px;
        margin: 0 auto;
        padding: 2rem clamp(1rem, 4vw, 1.5rem) 4rem;
      }
    `,
  ],
})
export class CollectionsComponent {
  private readonly api = inject(ApiService);
  vm$ = combineLatest([this.api.getServices(), this.api.getSite()]).pipe(
    map(([services, site]) => ({ services, site }))
  );
}
