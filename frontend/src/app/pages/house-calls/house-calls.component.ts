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
        eyebrow="◆ Private Viewing & Delivery"
        [title]="site.houseCalls.title"
        [subtitle]="site.houseCalls.description"
      />
      <section class="page-section">
        <ul class="notes">
          @for (n of site.houseCalls.notes; track n) {
            <li>
              <span class="icon">✦</span>
              <span>{{ n }}</span>
            </li>
          }
        </ul>
        <a class="btn primary" [href]="site.whatsAppUrl" target="_blank" rel="noopener">
          Book Private Viewing & Delivery →
        </a>
      </section>
    }
  `,
  styles: [
    `
      .page-section {
        max-width: 680px;
        margin: 0 auto;
        padding: 2.5rem var(--page-gutter) 4.5rem;
        text-align: center;
      }
      .notes {
        list-style: none;
        padding: 0;
        text-align: left;
        margin: 0 0 2.5rem 0;
        display: grid;
        gap: 1rem;

        li {
          display: flex;
          align-items: baseline;
          gap: 0.75rem;
          color: var(--color-text);
          font-size: 1rem;
          line-height: 1.6;

          .icon {
            color: var(--color-accent);
            font-weight: 800;
            flex-shrink: 0;
          }
        }
      }
      .btn.primary {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        min-height: 48px;
        padding: 0.9rem 2rem;
        border-radius: 999px;
        background: var(--color-accent);
        color: var(--color-primary-dark, #05070a);
        text-decoration: none;
        font-weight: 700;
        font-size: 1rem;
        box-shadow: 0 4px 16px color-mix(in srgb, var(--color-accent) 30%, transparent);
        transition: transform 0.2s, opacity 0.2s;

        &:hover {
          opacity: 0.94;
          transform: translateY(-2px);
        }
      }
    `,
  ],
})
export class HouseCallsComponent {
  private readonly api = inject(ApiService);
  site$ = this.api.getSite();
}
