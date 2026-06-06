import { DecimalPipe } from '@angular/common';
import { Component, Input } from '@angular/core';
import { RouterLink } from '@angular/router';
import { LashService, SiteConfig } from '../../core/models';
import { WebpImagePipe } from '../webp-image.pipe';

@Component({
  selector: 'app-service-grid',
  standalone: true,
  imports: [RouterLink, DecimalPipe, WebpImagePipe],
  templateUrl: './service-grid.component.html',
  styleUrl: './service-grid.component.scss',
})
export class ServiceGridComponent {
  @Input({ required: true }) services: LashService[] = [];
  @Input() site?: SiteConfig;
  @Input() compact = false;
}
