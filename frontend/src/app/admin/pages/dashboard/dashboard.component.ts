import { DatePipe } from '@angular/common';
import { Component, OnInit, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';
import { AdminApiService, BookingRecord } from '../../services/admin-api.service';

@Component({
  selector: 'app-admin-dashboard',
  standalone: true,
  imports: [RouterLink, DatePipe],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.scss',
})
export class DashboardComponent implements OnInit {
  private readonly api = inject(AdminApiService);

  counts = signal<Partial<Record<string, number>>>({});
  bookings = signal<BookingRecord[]>([]);

  ngOnInit(): void {
    this.api.dashboard().subscribe((d) => {
      this.counts.set(d.counts);
      this.bookings.set(d.recentBookings ?? []);
    });
  }
}
