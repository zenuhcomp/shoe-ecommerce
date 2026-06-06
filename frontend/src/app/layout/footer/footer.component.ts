import { Component, Input } from '@angular/core';
import { RouterLink } from '@angular/router';
import { NavLink, SiteConfig } from '../../core/models';

@Component({
  selector: 'app-footer',
  standalone: true,
  imports: [RouterLink],
  templateUrl: './footer.component.html',
  styleUrl: './footer.component.scss',
})
export class FooterComponent {
  @Input({ required: true }) site!: SiteConfig;
  @Input() nav: NavLink[] = [];

  telHref(phone: string): string {
    return 'tel:' + phone.replaceAll(' ', '');
  }
}
