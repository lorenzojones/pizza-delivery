import { writable } from 'svelte/store';

// Base URL for the API. Empty string => same-origin (e.g. "/api/dashboard").
// In dev, Vite proxies "/api" to the backend; in prod nginx does.
const API_BASE = import.meta.env.VITE_API_BASE ?? '';

const DEFAULT_REFRESH_SECONDS = 300;
const RETRY_AFTER_FAILURE_MS = 30_000;

/**
 * The dashboard store.
 *
 * Shape:
 * {
 *   data:       <last good /api/dashboard payload> | null,
 *   loading:    boolean   // true only before the first successful (or failed) fetch
 *   error:      string|null  // last fetch error message (cleared on success)
 *   lastFetch:  number|null  // ms epoch of the last *successful* fetch
 *   online:     boolean   // was the most recent fetch attempt successful?
 * }
 */
export const dashboard = writable({
  data: null,
  loading: true,
  error: null,
  lastFetch: null,
  online: false,
});

let timer = null;
let started = false;

function apiUrl(path) {
  // Avoid a double slash if API_BASE ends with one.
  const base = API_BASE.replace(/\/$/, '');
  return `${base}${path}`;
}

async function fetchDashboard() {
  const controller = new AbortController();
  const to = setTimeout(() => controller.abort(), 15_000);
  try {
    const res = await fetch(apiUrl('/api/dashboard'), {
      headers: { Accept: 'application/json' },
      signal: controller.signal,
      cache: 'no-store',
    });
    if (!res.ok) {
      throw new Error(`HTTP ${res.status}`);
    }
    const json = await res.json();
    return json;
  } finally {
    clearTimeout(to);
  }
}

function scheduleNext(seconds) {
  if (timer) clearTimeout(timer);
  const ms = Math.max(5, Number(seconds) || DEFAULT_REFRESH_SECONDS) * 1000;
  timer = setTimeout(tick, ms);
}

function scheduleRetry() {
  if (timer) clearTimeout(timer);
  timer = setTimeout(tick, RETRY_AFTER_FAILURE_MS);
}

async function tick() {
  try {
    const json = await fetchDashboard();
    dashboard.update((s) => ({
      ...s,
      data: json,
      loading: false,
      error: null,
      lastFetch: Date.now(),
      online: true,
    }));
    const interval =
      json && Number(json.refreshIntervalSeconds) > 0
        ? json.refreshIntervalSeconds
        : DEFAULT_REFRESH_SECONDS;
    scheduleNext(interval);
  } catch (err) {
    // Keep showing the last good data; just mark us offline/stale.
    dashboard.update((s) => ({
      ...s,
      loading: false,
      error: err && err.message ? err.message : 'fetch failed',
      online: false,
    }));
    // Retry sooner than the normal interval after a failure.
    scheduleRetry();
  }
}

/** Start polling. Safe to call once; subsequent calls are no-ops. */
export function startPolling() {
  if (started) return;
  started = true;
  tick();
}

/** Stop polling (used on teardown). */
export function stopPolling() {
  if (timer) clearTimeout(timer);
  timer = null;
  started = false;
}
