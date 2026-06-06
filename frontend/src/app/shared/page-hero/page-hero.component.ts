import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-page-hero',
  standalone: true,
  template: `
    <header class="page-hero" id="page-top">
      @if (eyebrow) {
        <p class="eyebrow">{{ eyebrow }}</p>
      }
      <h1 [innerHTML]="title"></h1>
      @if (subtitle) {
        <p class="subtitle">{{ subtitle }}</p>
      }
    </header>
  `,
  styleUrl: './page-hero.component.scss',
})
export class PageHeroComponent {
  @Input() eyebrow = '';
  @Input({ required: true }) title = '';
  @Input() subtitle = '';
}
