import { AsyncPipe } from '@angular/common';
import { Component, inject } from '@angular/core';
import { combineLatest, map } from 'rxjs';
import { ApiService } from '../../core/api.service';
import { PageHeroComponent } from '../../shared/page-hero/page-hero.component';

@Component({
  selector: 'app-reviews',
  standalone: true,
  imports: [AsyncPipe, PageHeroComponent],
  template: `
    @if (vm$ | async; as vm) {
      <app-page-hero
        eyebrow="◆ Client Testimonials"
        [title]="'Client <span class=\\'accent\\'>Experiences</span>'"
        [subtitle]="'Verified experiences and reviews from clients of ' + vm.site.brandName"
      />
      <section class="page-section review-list">
        @for (r of vm.reviews; track r.id) {
          <article class="review-item">
            <div class="avatar">{{ r.initials }}</div>
            <div class="review-body">
              <div class="author-row">
                <strong>{{ r.name }}</strong>
                <span class="handle">{{ r.handle }}</span>
              </div>
              <p class="stars">★★★★★</p>
              <blockquote>{{ r.quote }}</blockquote>
            </div>
          </article>
        }
      </section>
    }
  `,
  styles: [
    `
      .page-section {
        max-width: 840px;
        margin: 0 auto;
        padding: 2.5rem var(--page-gutter) 4.5rem;
      }
      .review-item {
        display: flex;
        gap: 1.25rem;
        padding: 1.75rem 0;
        border-bottom: 1px solid color-mix(in srgb, var(--color-text) 12%, transparent);
        color: var(--color-text);
      }
      .avatar {
        width: 52px;
        height: 52px;
        border-radius: 50%;
        background: var(--color-accent);
        color: var(--color-primary-dark, #05070a);
        display: flex;
        align-items: center;
        justify-content: center;
        font-weight: 800;
        font-size: 1.1rem;
        flex-shrink: 0;
      }
      .review-body {
        flex: 1;
      }
      .author-row {
        display: flex;
        align-items: baseline;
        gap: 0.5rem;
        flex-wrap: wrap;

        strong {
          font-size: 1.05rem;
          color: var(--color-text);
        }
        .handle {
          color: var(--color-accent);
          font-size: 0.85rem;
          font-weight: 600;
        }
      }
      .stars {
        color: #f59e0b;
        margin: 0.35rem 0 0.5rem;
        letter-spacing: 0.1em;
        font-size: 0.95rem;
      }
      blockquote {
        color: var(--color-text-muted);
        line-height: 1.75;
        margin: 0;
        font-size: 0.95rem;
      }
    `,
  ],
})
export class ReviewsComponent {
  private readonly api = inject(ApiService);
  vm$ = combineLatest([this.api.getSite(), this.api.getReviews()]).pipe(
    map(([site, reviews]) => ({ site, reviews }))
  );
}
