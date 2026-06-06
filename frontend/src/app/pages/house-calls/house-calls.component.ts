import { AsyncPipe } from '@angular/common';
import { Component, inject } from '@angular/core';
import { ApiService } from '../../core/api.service';
import { PageHeroComponent } from '../../shared/page-hero/page-hero.component';

@Component({
  selector: 'app-house-calls',
  standalone: true,
  imports: [AsyncPipe, PageHeroComponent],
  template: `
    @if (site$ | async; as site) {
      <app-page-hero
        eyebrow="◆ Delivery"
        [title]="site.houseCalls.title"
        [subtitle]="site.houseCalls.description"
      />
      <section class="page-section">
        <ul class="notes">
          @for (n of site.houseCalls.notes; track n) {
            <li>{{ n }}</li>
          }
        </ul>
        <a class="btn" [href]="site.whatsAppUrl" target="_blank" rel="noopener">Order for Delivery</a>
      </section>
    }
  `,
  styles: [
    `
      .page-section {
        max-width: 640px;
        margin: 0 auto;
        padding: 2rem var(--page-gutter) 4rem;
        text-align: center;
      }
      .notes {
        text-align: left;
        margin-bottom: 2rem;
        line-height: 1.8;
        padding-left: 1.25rem;
      }
      .btn {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        min-height: 48px;
        padding: 0.9rem 1.75rem;
        border-radius: 999px;
        background: var(--color-primary);
        color: #fff;
        text-decoration: none;
        font-weight: 600;
      }
    `,
  ],
})
export class HouseCallsComponent {
  private readonly api = inject(ApiService);
  site$ = this.api.getSite();
}
