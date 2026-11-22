import { UpperCasePipe } from '@angular/common';
import { Component, OnInit, inject, signal } from '@angular/core';
import { FormArray, FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { AdminApiService } from '../../services/admin-api.service';

@Component({
  selector: 'app-site-editor',
  standalone: true,
  imports: [ReactiveFormsModule, UpperCasePipe],
  templateUrl: './site-editor.component.html',
  styleUrl: './site-editor.component.scss',
})
export class SiteEditorComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly api = inject(AdminApiService);

  activeTab = signal('brand');
  message = signal('');
  saving = signal(false);

  form = this.fb.group({
    industry: [''],
    brandName: ['', Validators.required],
    brandScript: [''],
    tagline: [''],
    heroTitle: [''],
    heroSubtitle: [''],
    heroBadge: [''],
    logoUrl: [''],
    heroVideoUrl: [''],
    heroBackgroundType: ['video'],
    heroOverlayOpacity: [40],
    announcement: [''],
    location: [''],
    studioAddress: [''],
    phone: [''],
    email: [''],
    whatsAppUrl: [''],
    instagramUrl: [''],
    tiktokUrl: [''],
    priceRange: [''],
    copyright: [''],
    metaDescription: [''],
    theme: this.fb.group({
      primary: [''],
      primaryDark: [''],
      accent: [''],
      accentLight: [''],
      background: [''],
      surface: [''],
      text: [''],
      textMuted: [''],
      heroGradient: [''],
      fontHeading: [''],
      fontBody: [''],
      fontScript: [''],
    }),
    animations: this.fb.group({
      globalEnabled: [true],
      intensity: ['balanced'],
      scrollReveal: [true],
      defaultDuration: ['0.8s'],
      cardHoverEffect: ['lift'],
    }),
    about: this.fb.group({
      title: [''],
      intro: [''],
      story: [''],
      artistName: [''],
      artistTitle: [''],
      highlights: this.fb.array<string>([]),
    }),
    houseCalls: this.fb.group({
      title: [''],
      description: [''],
      notes: this.fb.array<string>([]),
    }),
    aftercare: this.fb.group({
      title: [''],
      intro: [''],
      tips: this.fb.array<ReturnType<typeof this.fb.group>>([]),
    }),
    pricingHighlight: this.fb.group({
      title: [''],
      subtitle: [''],
      priceFrom: [0],
      priceTo: [0],
      currency: ['KES'],
      badge: [''],
      includes: [''],
      studioNote: [''],
    }),
  });

  ngOnInit(): void {
    this.api.getSite().subscribe((site) => {
      this.form.patchValue({
        ...site,
        theme: { ...site.theme },
        animations: site.animations ? { ...site.animations } : {
          globalEnabled: true,
          intensity: 'balanced',
          scrollReveal: true,
          defaultDuration: '0.8s',
          cardHoverEffect: 'lift',
        },
        about: { ...site.about, highlights: undefined },
        houseCalls: { ...site.houseCalls, notes: undefined },
        aftercare: { title: site.aftercare.title, intro: site.aftercare.intro },
        pricingHighlight: { ...site.pricingHighlight },
      });
      this.setStringArray('about', 'highlights', site.about?.highlights || []);
      this.setStringArray('houseCalls', 'notes', site.houseCalls?.notes || []);
      this.setTips(site.aftercare?.tips || []);
    });
  }

  highlights(): FormArray {
    return this.form.get('about.highlights') as FormArray;
  }

  houseNotes(): FormArray {
    return this.form.get('houseCalls.notes') as FormArray;
  }

  tips(): FormArray {
    return this.form.get('aftercare.tips') as FormArray;
  }

  addHighlight(): void {
    this.highlights().push(this.fb.control(''));
  }

  addHouseNote(): void {
    this.houseNotes().push(this.fb.control(''));
  }

  addTip(): void {
    this.tips().push(this.fb.group({ title: [''], description: [''] }));
  }

  save(): void {
    if (this.form.invalid) return;
    this.saving.set(true);
    const raw = this.form.getRawValue();
    const payload = {
      ...raw,
      about: {
        ...raw.about,
        highlights: raw.about?.highlights?.filter((h) => h?.trim()) ?? [],
      },
      houseCalls: {
        ...raw.houseCalls,
        notes: raw.houseCalls?.notes?.filter((n) => n?.trim()) ?? [],
      },
      aftercare: {
        title: raw.aftercare?.title ?? '',
        intro: raw.aftercare?.intro ?? '',
        tips: raw.aftercare?.tips ?? [],
      },
    };
    this.api.putSite(payload as never).subscribe({
      next: () => {
        this.saving.set(false);
        this.message.set('Site settings saved successfully.');
      },
      error: () => {
        this.saving.set(false);
        this.message.set('Failed to save. Please try again.');
      },
    });
  }

  private setStringArray(group: string, field: string, values: string[]): void {
    const arr = this.form.get(`${group}.${field}`) as FormArray | null;
    if (!arr) return;
    arr.clear();
    values.forEach((v) => arr.push(this.fb.control(v)));
  }

  private setTips(tips: { title: string; description: string }[]): void {
    const arr = this.tips();
    arr.clear();
    tips.forEach((t) => arr.push(this.fb.group({ title: [t.title], description: [t.description] })));
  }
}
