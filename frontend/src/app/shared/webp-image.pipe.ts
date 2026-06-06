import { Pipe, PipeTransform } from '@angular/core';

const UNSPLASH_HOSTS = new Set(['images.unsplash.com', 'plus.unsplash.com']);

function isRelativeUrl(url: string): boolean {
  return url.startsWith('/') || url.startsWith('./') || url.startsWith('../');
}

function replaceLocalImageExtension(path: string): string {
  return path.replace(/\.(?:jpe?g|png)(?=([?#]|$))/i, '.webp');
}

export function toWebpImageUrl(src: string | null | undefined, width = 800, quality = 72): string {
  const raw = (src ?? '').trim();
  if (!raw) return '';
  if (raw.startsWith('data:') || raw.startsWith('blob:')) return raw;

  if (isRelativeUrl(raw)) {
    return replaceLocalImageExtension(raw);
  }

  try {
    const url = new URL(raw);

    if (UNSPLASH_HOSTS.has(url.hostname)) {
      url.searchParams.set('auto', 'format');
      url.searchParams.set('fit', 'crop');
      url.searchParams.set('fm', 'webp');
      url.searchParams.set('q', String(quality));
      url.searchParams.set('w', String(width));
      return url.toString();
    }

    return raw;
  } catch {
    return replaceLocalImageExtension(raw);
  }
}

@Pipe({
  name: 'webpImage',
  standalone: true,
  pure: true,
})
export class WebpImagePipe implements PipeTransform {
  transform(src: string | null | undefined, width = 800, quality = 72): string {
    return toWebpImageUrl(src, width, quality);
  }
}
