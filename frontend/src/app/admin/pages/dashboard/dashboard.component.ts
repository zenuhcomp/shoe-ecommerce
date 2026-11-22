import { DatePipe, DecimalPipe, UpperCasePipe } from '@angular/common';
import { Component, OnInit, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';
import { AdminApiService, BookingRecord } from '../../services/admin-api.service';
import { OrderRecord, SiteConfig } from '../../../core/models';

@Component({
  selector: 'app-admin-dashboard',
  standalone: true,
  imports: [RouterLink, DatePipe, DecimalPipe, UpperCasePipe],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.scss',
})
export class DashboardComponent implements OnInit {
  private readonly api = inject(AdminApiService);

  counts = signal<Partial<Record<string, number>>>({});
  bookings = signal<BookingRecord[]>([]);
  orders = signal<OrderRecord[]>([]);
  site = signal<SiteConfig | null>(null);

  ngOnInit(): void {
    this.api.dashboard().subscribe((d) => {
      this.counts.set(d.counts);
      this.bookings.set(d.recentBookings ?? []);
      this.orders.set(d.recentOrders ?? []);
    });
    this.api.getSite().subscribe((s) => this.site.set(s));
  }

  getWhatsAppReply(phone: string, name: string): string {
    const cleanPhone = phone.replace(/[^0-9]/g, '');
    const brand = this.site()?.brandName || 'our team';
    const msg = encodeURIComponent(`Hi ${name}, thank you for contacting ${brand}. How can we assist you today?`);
    return `https://wa.me/${cleanPhone}?text=${msg}`;
  }
}
