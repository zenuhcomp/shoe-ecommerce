import { AsyncPipe } from '@angular/common';
import { Component, inject } from '@angular/core';
import { ApiService } from '../../core/api.service';
import { PageHeroComponent } from '../../shared/page-hero/page-hero.component';

@Component({
  selector: 'app-aftercare',
  standalone: true,
  imports: [AsyncPipe, PageHeroComponent],
  template: `
    @if (site$ | async; as site) {
      <app-page-hero [title]="site.aftercare.title" [subtitle]="site.aftercare.intro" />
      <section class="page-section tips">
        @for (tip of site.aftercare.tips; track tip.title) {
          <article>
            <h3>{{ tip.title }}</h3>
            <p>{{ tip.description }}</p>
          </article>
        }
      </section>
    }
  `,
  styles: [
    `
      .page-section {
        max-width: 800px;
        margin: 0 auto;
        padding: 2rem 1.5rem 4rem;
        display: grid;
        gap: 1.25rem;
      }
      article {
        padding: 1.25rem 1.5rem;
        background: var(--color-surface);
        border-radius: 16px;
        border-left: 4px solid var(--color-accent);
      }
      h3 {
        font-family: var(--font-heading);
        margin-bottom: 0.5rem;
      }
      p {
        color: var(--color-text-muted);
        line-height: 1.7;
      }
    `,
  ],
})
export class AftercareComponent {
  private readonly api = inject(ApiService);
  site$ = this.api.getSite();
}
