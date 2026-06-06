import { DecimalPipe } from '@angular/common';
import { Component, Input, signal } from '@angular/core';
import { FAQ } from '../../core/models';

@Component({
  selector: 'app-faq-list',
  standalone: true,
  imports: [DecimalPipe],
  templateUrl: './faq-list.component.html',
  styleUrl: './faq-list.component.scss',
})
export class FaqListComponent {
  @Input({ required: true }) items: FAQ[] = [];
  openId = signal<number | null>(null);

  toggle(id: number): void {
    this.openId.update((current) => (current === id ? null : id));
  }
}
