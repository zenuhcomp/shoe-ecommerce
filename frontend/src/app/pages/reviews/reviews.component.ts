import { AsyncPipe } from '@angular/common';
import { Component, inject } from '@angular/core';
import { ApiService } from '../../core/api.service';
import { PageHeroComponent } from '../../shared/page-hero/page-hero.component';

@Component({
  selector: 'app-reviews',
  standalone: true,
  imports: [AsyncPipe, PageHeroComponent],
  template: `
    <app-page-hero
      eyebrow="◆ Testimonials"
      title="Customer <span class='accent'>Reviews</span>"
      subtitle="What professionals say about stepping out in Luna Collections"
    />
    @if (reviews$ | async; as reviews) {
      <section class="page-section review-list">
        @for (r of reviews; track r.id) {
          <article>
            <div class="avatar">{{ r.initials }}</div>
            <div>
              <strong>{{ r.name }}</strong> <span>{{ r.handle }}</span>
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
        max-width: 800px;
        margin: 0 auto;
        padding: 2rem var(--page-gutter) 4rem;
      }
      article {
        display: flex;
        gap: 1rem;
        padding: 1.5rem 0;
        border-bottom: 1px solid rgba(26, 31, 59, 0.12);
      }
      .avatar {
        width: 52px;
        height: 52px;
        border-radius: 50%;
        background: var(--color-primary);
        color: #fff;
        display: flex;
        align-items: center;
        justify-content: center;
        font-weight: 700;
        flex-shrink: 0;
      }
      span {
        color: var(--color-text-muted);
        font-size: 0.85rem;
        margin-left: 0.35rem;
      }
      .stars {
        color: var(--color-accent);
        margin: 0.35rem 0;
      }
      blockquote {
        color: var(--color-text-muted);
        line-height: 1.7;
        margin: 0;
      }
    `,
  ],
})
export class ReviewsComponent {
  private readonly api = inject(ApiService);
  reviews$ = this.api.getReviews();
}
