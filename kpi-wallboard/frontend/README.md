# Sales Wallboard — Frontend

A real-time, always-on **TV wallboard** for a sales team. Landscape 16:9,
glanceable in under 5 seconds: a competitive leaderboard, a top-performer hero
card, team totals, a call-volume donut, and a scrolling stats ticker. Built with
**Svelte 4 + Vite** and **D3** for all charts. Dark, high-contrast theme tuned
for 1080p and 4K displays.

## Develop

```bash
npm install
npm run dev
```

The app polls `GET /api/dashboard` and re-polls every `refreshIntervalSeconds`
(default 300s). In dev, Vite proxies `/api` to the Go backend at
`http://localhost:8080` (the backend lives in `../backend`). Override the proxy
target with `VITE_DEV_API_TARGET`.

To point the build at a different API origin, set `VITE_API_BASE` (default `""`,
i.e. same-origin):

```bash
VITE_API_BASE=https://api.example.com npm run build
```

## Build

```bash
npm run build     # outputs to dist/
npm run preview   # serve the production build locally
```

## Docker

Multi-stage build (`node:22-alpine` → `nginx:alpine`). nginx serves the SPA on
port 80 and proxies `/api/` to the `backend` service (docker-compose).

```bash
docker build -t wallboard-frontend .
```

## Backend

This is the frontend only. It expects the Go backend (`../backend`) to expose
`GET /api/dashboard` returning the dashboard payload (title, ranked people, team
totals, top performer, biggest mover, per-source status, scoring weights). If the
backend is down or stale, the UI keeps showing the last good data and surfaces a
warning banner instead of blanking the screen.

## Structure

- `src/lib/api.js` — fetch + polling, exposes a Svelte store.
- `src/lib/charts/` — reusable D3 components (`LeaderboardBars`, `CallBreakdown`,
  `TargetProgress`), each ResizeObserver-driven and animated.
- `src/lib/components/` — UI (`Header`, `SourceBadge`, `Leaderboard`,
  `TopPerformerCard`, `TeamTotals`, `Ticker`, `StaleBanner`).
