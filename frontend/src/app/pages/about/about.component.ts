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
        <ul>
          @for (h of site.about.highlights; track h) {
            <li>{{ h }}</li>
          }
        </ul>
        <p class="location">Showroom: {{ site.studioAddress }}, {{ site.location }}</p>
      </section>
    }
  `,
  styles: [
    `
      .page-section {
        max-width: 720px;
        margin: 0 auto;
        padding: 2rem 1.5rem 4rem;
      }
      .intro {
        font-family: var(--font-script);
        font-size: 1.5rem;
        color: var(--color-primary);
        margin-bottom: 1rem;
      }
      .story {
        line-height: 1.8;
        color: var(--color-text-muted);
        margin-bottom: 1.5rem;
      }
      ul {
        padding-left: 1.25rem;
        margin-bottom: 2rem;
        li {
          margin-bottom: 0.5rem;
        }
      }
      .location {
        font-weight: 600;
        color: var(--color-primary);
      }
    `,
  ],
})
export class AboutComponent {
  private readonly api = inject(ApiService);
  site$ = this.api.getSite();
}
