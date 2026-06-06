import { AsyncPipe, DecimalPipe } from '@angular/common';
import { Component, inject } from '@angular/core';
import { RouterLink } from '@angular/router';
import { combineLatest, map } from 'rxjs';
import { ApiService } from '../../core/api.service';
import { SectionHeadingComponent } from '../../shared/section-heading/section-heading.component';
import { ServiceGridComponent } from '../../shared/service-grid/service-grid.component';
import { FaqListComponent } from '../../shared/faq-list/faq-list.component';
import { WebpImagePipe } from '../../shared/webp-image.pipe';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [
    AsyncPipe,
    DecimalPipe,
    RouterLink,
    SectionHeadingComponent,
    ServiceGridComponent,
    FaqListComponent,
    WebpImagePipe,
  ],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss',
})
export class HomeComponent {
  private readonly api = inject(ApiService);
  home$ = combineLatest([this.api.getHome(), this.api.getGallery()]).pipe(
    map(([home, gallery]) => ({ ...home, gallery: gallery.slice(0, 6) }))
  );
}
