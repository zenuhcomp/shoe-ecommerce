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
    <app-page-hero
      eyebrow="◆ Contact"
      title="Order Your <span class='accent'>Perfect Pair</span>"
      subtitle="Request a style online or message us on WhatsApp for sizing help"
    />
    @if (vm$ | async; as vm) {
      <section class="page-section contact-grid">
        <div class="info">
          <h2>Visit the showroom</h2>
          <p>{{ vm.site.studioAddress }}</p>
          <p>{{ vm.site.location }}</p>
          <h2>Reach us</h2>
          <p><a [href]="telHref(vm.site.phone)">{{ vm.site.phone }}</a></p>
          <p><a [href]="'mailto:' + vm.site.email">{{ vm.site.email }}</a></p>
          <a class="whatsapp" [href]="vm.site.whatsAppUrl" target="_blank" rel="noopener">Chat on WhatsApp</a>
        </div>
        <div class="form-wrap">
          <h2>Send an order inquiry</h2>
          <app-booking-form [services]="vm.services" />
        </div>
      </section>
    }
  `,
  styles: [
    `
      .page-section {
        max-width: 1000px;
        margin: 0 auto;
        padding: 2rem var(--page-gutter) 4rem;
        display: grid;
        gap: 2rem;
      }
      @media (min-width: 768px) {
        .contact-grid {
          grid-template-columns: 1fr 1fr;
        }
      }
      h2 {
        font-family: var(--font-heading);
        font-size: 1.25rem;
        margin: 1.25rem 0 0.5rem;
        &:first-child {
          margin-top: 0;
        }
      }
      a {
        color: var(--color-primary);
      }
      .whatsapp {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        min-height: 44px;
        margin-top: 1rem;
        padding: 0.75rem 1.25rem;
        background: #25d366;
        color: #fff !important;
        border-radius: 999px;
        text-decoration: none;
        font-weight: 600;
      }
      .form-wrap {
        background: var(--color-surface);
        padding: clamp(1.25rem, 4vw, 1.5rem);
        border-radius: 20px;
        border: 1px solid rgba(26, 31, 59, 0.1);
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
