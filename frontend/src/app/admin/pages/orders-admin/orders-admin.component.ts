import { DatePipe, DecimalPipe } from '@angular/common';
import { Component, OnInit, inject, signal } from '@angular/core';
import { AdminApiService, BookingRecord } from '../../services/admin-api.service';
import { ConfirmDialogService } from '../../../shared/confirm-dialog/confirm-dialog.service';
import { OrderRecord, SiteConfig } from '../../../core/models';

@Component({
  selector: 'app-orders-admin',
  standalone: true,
  imports: [DatePipe, DecimalPipe],
  templateUrl: './orders-admin.component.html',
  styleUrl: './orders-admin.component.scss',
})
export class OrdersAdminComponent implements OnInit {
  private readonly api = inject(AdminApiService);
  private readonly dialog = inject(ConfirmDialogService);

  orders = signal<OrderRecord[]>([]);
  bookings = signal<BookingRecord[]>([]);
  site = signal<SiteConfig | null>(null);
  activeTab = signal<'orders' | 'bookings'>('orders');
  statusFilter = signal<string>('all');
  loading = signal(false);
  message = signal<{ type: 'success' | 'error'; text: string } | null>(null);

  ngOnInit(): void {
    this.loadData();
  }

  loadData(): void {
    this.loading.set(true);
    this.api.listOrders().subscribe({
      next: (orders) => {
        this.orders.set(orders);
        this.api.listBookings().subscribe({
          next: (bookings) => {
            this.bookings.set(bookings);
            this.api.getSite().subscribe({
              next: (s) => {
                this.site.set(s);
                this.loading.set(false);
              },
              error: () => this.loading.set(false),
            });
          },
          error: () => this.loading.set(false),
        });
      },
      error: () => {
        this.loading.set(false);
        this.message.set({ type: 'error', text: 'Could not load orders' });
      },
    });
  }

  setTab(tab: 'orders' | 'bookings'): void {
    this.activeTab.set(tab);
  }

  setStatusFilter(filter: string): void {
    this.statusFilter.set(filter);
  }

  filteredOrders(): OrderRecord[] {
    const f = this.statusFilter();
    const list = this.orders();
    if (f === 'all') return list;
    return list.filter((o) => o.status === f);
  }

  updateStatus(order: OrderRecord, nextStatus: string): void {
    this.api.updateOrderStatus(order.id, nextStatus).subscribe({
      next: () => {
        this.orders.update((list) =>
          list.map((o) => (o.id === order.id ? { ...o, status: nextStatus } : o))
        );
      },
      error: () => this.message.set({ type: 'error', text: 'Could not update status' }),
    });
  }

  async deleteOrder(id: string): Promise<void> {
    const confirmed = await this.dialog.confirm({
      title: 'Delete Order Record?',
      message: 'This order record will be permanently deleted from the database.',
      confirmText: 'Delete Order',
      cancelText: 'Cancel',
      variant: 'danger',
      icon: '🗑️',
      badge: 'ORDER RECORD',
    });
    if (!confirmed) return;
    this.api.deleteOrder(id).subscribe({
      next: () => {
        this.orders.update((list) => list.filter((o) => o.id !== id));
      },
      error: () => this.message.set({ type: 'error', text: 'Could not delete order' }),
    });
  }

  async deleteBooking(id: string): Promise<void> {
    const confirmed = await this.dialog.confirm({
      title: 'Delete Booking Inquiry?',
      message: 'This booking inquiry will be permanently deleted.',
      confirmText: 'Delete Inquiry',
      cancelText: 'Cancel',
      variant: 'danger',
      icon: '🗑️',
      badge: 'INQUIRY RECORD',
    });
    if (!confirmed) return;
    this.api.deleteBooking(id).subscribe({
      next: () => {
        this.bookings.update((list) => list.filter((b) => b.id !== id));
      },
      error: () => this.message.set({ type: 'error', text: 'Could not delete booking' }),
    });
  }

  getWhatsAppReply(phone: string, name: string, itemName?: string): string {
    const cleanPhone = phone.replace(/[^0-9]/g, '');
    const brand = this.site()?.brandName || 'our team';
    const msg = encodeURIComponent(`Hi ${name}, thank you for contacting ${brand} regarding ${itemName || 'your inquiry'}. How can we assist you today?`);
    return `https://wa.me/${cleanPhone}?text=${msg}`;
  }
}
