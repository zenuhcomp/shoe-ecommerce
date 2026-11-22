import { AsyncPipe } from '@angular/common';
import { Component, inject } from '@angular/core';
import { combineLatest, map } from 'rxjs';
import { ApiService } from '../../core/api.service';
import { PageHeroComponent } from '../../shared/page-hero/page-hero.component';
import { BookingFormComponent } from '../../shared/booking-form/booking-form.component';

@Component({
  selector: 'app-contact',
  standalone: true,
  imports: [AsyncPipe, PageHeroComponent, BookingFormComponent],
  template: `
    @if (vm$ | async; as vm) {
      <app-page-hero
        eyebrow="◆ Client Services"
        [title]="'Connect With ' + vm.site.brandName + ' <span class=\\'accent\\'>Concierge</span>'"
        subtitle="Submit an inquiry or order request online, or message our private client team on WhatsApp"
      />
      <section class="page-section contact-grid">
        <div class="info-card">
          <h2>Showroom & Location</h2>
          <p>{{ vm.site.studioAddress }}</p>
          <p class="loc-sub">{{ vm.site.location }}</p>

          <h2>Direct Communications</h2>
          <p><a [href]="telHref(vm.site.phone)">{{ vm.site.phone }}</a></p>
          <p><a [href]="'mailto:' + vm.site.email">{{ vm.site.email }}</a></p>

          <a class="whatsapp" [href]="vm.site.whatsAppUrl" target="_blank" rel="noopener">
            💬 Connect on WhatsApp
          </a>
        </div>
        <div class="form-wrap">
          <h2>Send an inquiry or reservation</h2>
          <app-booking-form [services]="vm.services" />
        </div>
      </section>
    }
  `,
  styles: [
    `
      .page-section {
        max-width: 1100px;
        margin: 0 auto;
        padding: 2.5rem var(--page-gutter) 4.5rem;
        display: grid;
        gap: 2rem;
      }
      @media (min-width: 768px) {
        .contact-grid {
          grid-template-columns: 1fr 1.2fr;
          align-items: start;
        }
      }
      .info-card {
        background: var(--color-surface);
        color: var(--color-text);
        padding: clamp(1.5rem, 4vw, 2rem);
        border-radius: 20px;
        border: 1px solid color-mix(in srgb, var(--color-text) 12%, transparent);
      }
      h2 {
        font-family: var(--font-heading);
        font-size: 1.3rem;
        font-weight: 700;
        color: var(--color-text);
        margin: 1.5rem 0 0.5rem;
        &:first-child {
          margin-top: 0;
        }
      }
      p {
        margin: 0.25rem 0;
        color: var(--color-text-muted);
      }
      .loc-sub {
        color: var(--color-accent);
        font-weight: 600;
      }
      a {
        color: var(--color-accent);
        text-decoration: none;
        font-weight: 600;
        &:hover {
          text-decoration: underline;
        }
      }
      .whatsapp {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        min-height: 48px;
        margin-top: 1.5rem;
        padding: 0.85rem 1.5rem;
        background: #25d366;
        color: #ffffff !important;
        border-radius: 999px;
        text-decoration: none !important;
        font-weight: 700;
        box-shadow: 0 4px 14px rgba(37, 211, 102, 0.3);
        transition: transform 0.2s, opacity 0.2s;

        &:hover {
          opacity: 0.94;
          transform: translateY(-1px);
        }
      }
      .form-wrap {
        background: var(--color-surface);
        color: var(--color-text);
        padding: clamp(1.5rem, 4vw, 2rem);
        border-radius: 20px;
        border: 1px solid color-mix(in srgb, var(--color-text) 12%, transparent);
      }
    `,
  ],
})
export class ContactComponent {
  private readonly api = inject(ApiService);
  vm$ = combineLatest([this.api.getSite(), this.api.getServices()]).pipe(
    map(([site, services]) => ({ site, services }))
  );

  telHref(phone: string): string {
    return 'tel:' + phone.replaceAll(' ', '');
  }
}
