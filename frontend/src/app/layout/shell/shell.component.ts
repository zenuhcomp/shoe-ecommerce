import { AsyncPipe } from '@angular/common';
import { Component, inject } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { switchMap } from 'rxjs';
import { ApiService } from '../../core/api.service';
import { HeaderComponent } from '../header/header.component';
import { FooterComponent } from '../footer/footer.component';

@Component({
  selector: 'app-shell',
  standalone: true,
  imports: [RouterOutlet, HeaderComponent, FooterComponent, AsyncPipe],
  template: `
    @if (home$ | async; as home) {
      <div class="announcement-bar">
        <span class="announcement-text">{{ home.site.announcement }}</span>
      </div>
      <app-header [site]="home.site" [nav]="home.nav" />
      <main class="main-content">
        <router-outlet />
      </main>
      <app-footer [site]="home.site" [nav]="home.nav" />
      <a class="sticky-book" [href]="home.site.whatsAppUrl" target="_blank" rel="noopener">Shop Now</a>
    }
  `,
  styleUrl: './shell.component.scss',
})
export class ShellComponent {
  private readonly api = inject(ApiService);
  home$ = this.api.preloadTheme().pipe(switchMap(() => this.api.getHome()));
}
