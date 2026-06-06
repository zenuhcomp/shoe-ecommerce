import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-section-heading',
  standalone: true,
  template: `
    <div class="section-heading" [class.center]="center" [class.light]="light">
      @if (eyebrow) {
        <p class="eyebrow">{{ eyebrow }}</p>
      }
      <h2 [innerHTML]="title"></h2>
      @if (subtitle) {
        <p class="subtitle">{{ subtitle }}</p>
      }
    </div>
  `,
  styleUrl: './section-heading.component.scss',
})
export class SectionHeadingComponent {
  @Input() eyebrow = '';
  @Input({ required: true }) title = '';
  @Input() subtitle = '';
  @Input() center = true;
  @Input() light = false;
}
