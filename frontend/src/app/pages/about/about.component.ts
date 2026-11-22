import { AsyncPipe } from '@angular/common';
import { Component, inject } from '@angular/core';
import { ApiService } from '../../core/api.service';
import { PageHeroComponent } from '../../shared/page-hero/page-hero.component';

@Component({
  selector: 'app-about',
  standalone: true,
  imports: [AsyncPipe, PageHeroComponent],
  template: `
    @if (site$ | async; as site) {
      <app-page-hero
        [eyebrow]="'◆ ' + site.about.title"
        [title]="site.about.artistName"
        [subtitle]="site.about.artistTitle"
      />
      <section class="page-section about-content">
        <p class="intro">{{ site.about.intro }}</p>
        <p class="story">{{ site.about.story }}</p>
        <ul class="highlights">
          @for (h of site.about.highlights; track h) {
            <li>
              <span class="bullet">✓</span>
              <span>{{ h }}</span>
            </li>
          }
        </ul>
        <div class="location-card">
          <strong>Showroom & Atelier:</strong>
          <p>{{ site.studioAddress }}, {{ site.location }}</p>
        </div>
      </section>
    }
  `,
  styles: [
    `
      .page-section {
        max-width: 800px;
        margin: 0 auto;
        padding: 2.5rem var(--page-gutter) 4.5rem;
      }
      .intro {
        font-family: var(--font-heading);
        font-size: 1.4rem;
        font-weight: 700;
        color: var(--color-accent);
        margin-bottom: 1.25rem;
        line-height: 1.5;
      }
      .story {
        line-height: 1.8;
        color: var(--color-text);
        margin-bottom: 2rem;
        font-size: 1.05rem;
      }
      .highlights {
        list-style: none;
        padding: 0;
        margin: 0 0 2rem 0;
        display: grid;
        gap: 0.85rem;

        li {
          display: flex;
          align-items: baseline;
          gap: 0.75rem;
          color: var(--color-text);
          font-size: 0.95rem;
          line-height: 1.6;

          .bullet {
            color: var(--color-accent);
            font-weight: 800;
            font-size: 1.1rem;
            flex-shrink: 0;
          }
        }
      }
      .location-card {
        padding: 1.25rem 1.5rem;
        background: var(--color-surface);
        border: 1px solid color-mix(in srgb, var(--color-text) 12%, transparent);
        border-radius: 14px;
        color: var(--color-text);

        strong {
          color: var(--color-accent);
          display: block;
          margin-bottom: 0.35rem;
        }

        p {
          margin: 0;
          color: var(--color-text-muted);
        }
      }
    `,
  ],
})
export class AboutComponent {
  private readonly api = inject(ApiService);
  site$ = this.api.getSite();
}
