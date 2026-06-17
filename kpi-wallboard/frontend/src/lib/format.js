// Shared formatting helpers. GBP currency, compact numbers, time, etc.

const gbp0 = new Intl.NumberFormat('en-GB', {
  style: 'currency',
  currency: 'GBP',
  maximumFractionDigits: 0,
});

const gbpCompact = new Intl.NumberFormat('en-GB', {
  style: 'currency',
  currency: 'GBP',
  notation: 'compact',
  maximumFractionDigits: 1,
});

const num0 = new Intl.NumberFormat('en-GB', { maximumFractionDigits: 0 });

/** Format a currency value in GBP, no decimals. Null/NaN -> "£0". */
export function money(v) {
  const n = Number(v);
  if (!isFinite(n)) return '£0';
  return gbp0.format(n);
}

/** Compact currency for tight spots, e.g. £12k. */
export function moneyCompact(v) {
  const n = Number(v);
  if (!isFinite(n)) return '£0';
  return gbpCompact.format(n);
}

/** Plain integer with thousands separators. */
export function int(v) {
  const n = Number(v);
  if (!isFinite(n)) return '0';
  return num0.format(Math.round(n));
}

/** "—" placeholder for missing data. */
export const DASH = '—';

/** Absolute clock time HH:MM:SS from an ISO string or Date. */
export function clock(d) {
  const dt = d instanceof Date ? d : new Date(d);
  if (isNaN(dt.getTime())) return DASH;
  return dt.toLocaleTimeString('en-GB', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
  });
}

/** Long, readable date e.g. "Wednesday, 17 June 2026". */
export function longDate(d) {
  const dt = d instanceof Date ? d : new Date(d);
  if (isNaN(dt.getTime())) return DASH;
  return dt.toLocaleDateString('en-GB', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  });
}

/** Relative "x ago" given an ISO/Date and a reference now (ms). */
export function relativeAgo(d, nowMs = Date.now()) {
  const dt = d instanceof Date ? d : new Date(d);
  if (isNaN(dt.getTime())) return DASH;
  const diff = Math.max(0, Math.floor((nowMs - dt.getTime()) / 1000));
  if (diff < 5) return 'just now';
  if (diff < 60) return `${diff}s ago`;
  const m = Math.floor(diff / 60);
  if (m < 60) return `${m}m ago`;
  const h = Math.floor(m / 60);
  if (h < 24) return `${h}h ago`;
  const days = Math.floor(h / 24);
  return `${days}d ago`;
}

/** Clamp a 0..1+ progress value to a percentage label. */
export function pct(v) {
  const n = Number(v);
  if (!isFinite(n)) return '0%';
  return `${Math.round(n * 100)}%`;
}

/** Initials for an avatar chip, e.g. "Jane Doe" -> "JD". */
export function initials(name) {
  if (!name) return '?';
  const parts = String(name).trim().split(/\s+/).slice(0, 2);
  return parts.map((p) => p[0] || '').join('').toUpperCase() || '?';
}
