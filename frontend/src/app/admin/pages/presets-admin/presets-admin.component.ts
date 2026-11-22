import { UpperCasePipe } from '@angular/common';
import { Component, OnInit, inject, signal } from '@angular/core';
import { AdminApiService } from '../../services/admin-api.service';
import { ApiService } from '../../../core/api.service';
import { ConfirmDialogService } from '../../../shared/confirm-dialog/confirm-dialog.service';
import { IndustryPreset, SiteConfig } from '../../../core/models';

@Component({
  selector: 'app-presets-admin',
  standalone: true,
  imports: [UpperCasePipe],
  templateUrl: './presets-admin.component.html',
  styleUrl: './presets-admin.component.scss',
})
export class PresetsAdminComponent implements OnInit {
  private readonly adminApi = inject(AdminApiService);
  private readonly publicApi = inject(ApiService);
  private readonly dialog = inject(ConfirmDialogService);

  presets = signal<IndustryPreset[]>([]);
  currentSite = signal<SiteConfig | null>(null);
  loading = signal(false);
  applyingId = signal<string | null>(null);
  message = signal<{ type: 'success' | 'error'; text: string } | null>(null);

  ngOnInit(): void {
    this.loadData();
  }

  loadData(): void {
    this.loading.set(true);
    this.adminApi.listPresets().subscribe({
      next: (presets) => {
        this.presets.set(presets);
        this.loadSite();
      },
      error: () => {
        this.loading.set(false);
        this.message.set({ type: 'error', text: 'Could not load industry presets.' });
      },
    });
  }

  loadSite(): void {
    this.adminApi.getSite().subscribe({
      next: (site) => {
        this.currentSite.set(site);
        this.loading.set(false);
      },
      error: () => this.loading.set(false),
    });
  }

  isCurrent(preset: IndustryPreset): boolean {
    const site = this.currentSite();
    if (!site) return false;
    return site.industry === preset.industry || site.industry === preset.id;
  }

  async applyPreset(preset: IndustryPreset, seedItems: boolean = true): Promise<void> {
    const confirmed = await this.dialog.confirm({
      title: `Apply ${preset.name} Preset?`,
      message: `This will reconfigure site branding, hero video, theme colors, typography, and section animations.\n\nIt will also replace catalog items with curated ${preset.industry} inventory, features, and showroom guarantees.`,
      confirmText: `Apply ${preset.name}`,
      cancelText: 'Keep Current',
      variant: 'accent',
      icon: preset.icon || '🚀',
      badge: `${preset.industry.toUpperCase()} PRESET`,
    });
    if (!confirmed) return;

    this.applyingId.set(preset.id);
    this.message.set(null);

    this.adminApi.applyPreset(preset.id, seedItems).subscribe({
      next: () => {
        this.applyingId.set(null);
        this.message.set({ type: 'success', text: `✨ Successfully switched to ${preset.name} (${preset.industry})!` });
        this.publicApi.clearCache();
        this.loadSite();
      },
      error: () => {
        this.applyingId.set(null);
        this.message.set({ type: 'error', text: 'Failed to apply preset. Check backend logs.' });
      },
    });
  }
}

