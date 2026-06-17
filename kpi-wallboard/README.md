# 🏆 Sales KPI Wallboard

A real-time, competitive **TV wallboard** for a sales team. It combines
per-person **call activity from Aircall** with **sales numbers from an Excel
spreadsheet**, ranks everyone with a configurable scoring model, and renders a
glanceable, leaderboard-style dashboard designed for an always-on office TV.

- **Backend:** Go — scheduled refresh, Aircall client, `.xlsx` parser,
  aggregation/scoring, HTTP API.
- **Frontend:** Svelte + Vite — full-screen 16:9 wallboard, auto-refreshing.
- **Charts:** D3 — animated leaderboard bars, call-volume breakdown, target
  progress.
- **Refresh cadence:** every 5 minutes (configurable), owned entirely by the
  backend. The frontend just consumes dashboard-ready JSON.

It is built to keep working when things go wrong: if Aircall is unreachable or
the spreadsheet is missing, the board keeps showing the last good data with a
clear stale/degraded warning instead of going blank.

---

## Quick start (Docker — recommended)

```bash
cd kpi-wallboard
docker compose up --build
```

Then open **http://localhost:3000** and put the browser into fullscreen on the
TV. With no Aircall credentials configured, the backend serves **deterministic
demo data**, so you get a fully working board immediately.

To go live, copy the env template and fill in real values:

```bash
cp .env.example .env
# edit .env: add AIRCALL_API_ID / AIRCALL_API_TOKEN, point at your spreadsheet, etc.
docker compose up --build
```

| Service  | URL                            | Notes                          |
| -------- | ------------------------------ | ------------------------------ |
| Frontend | http://localhost:3000          | The wallboard (nginx + SPA)    |
| Backend  | http://localhost:8080/api/...  | JSON API (proxied via nginx)   |

---

## How it works

```
Aircall API ─┐
             ├─► Go backend (refresh every 5 min) ─► /api/dashboard ─► Svelte + D3 wallboard
Sales .xlsx ─┘        • fetch + parse                    (JSON)            (auto-polls)
                      • match staff across sources
                      • score + rank
                      • cache last-good data
```

1. Every `REFRESH_INTERVAL_SECONDS` the backend fetches today's calls from
   Aircall and parses the sales spreadsheet.
2. It matches each spreadsheet row and Aircall user to a person in the **staff
   mapping**, computes a weighted **score**, and ranks the team.
3. The result is cached and served as a single JSON document.
4. The frontend polls that document on the same cadence and redraws.

If a source fails, its **last successful** values are reused and the board shows
a warning — it never blanks.

---

## Configuration

All configuration is via environment variables (Compose reads them from `.env`).
See [`.env.example`](./.env.example) for the full annotated list.

### Aircall credentials

Create an API key in Aircall: **Dashboard → Integrations & API → API Keys**.
Aircall uses HTTP Basic auth — the **API ID** is the username and the **API
Token** is the password.

```env
AIRCALL_API_ID=your_api_id
AIRCALL_API_TOKEN=your_api_token
AIRCALL_BASE_URL=https://api.aircall.io/v1
```

Leave both blank to stay in demo mode. The client handles pagination, rate
limits (429 + `Retry-After`), and retries transient 5xx errors with exponential
backoff. The current "daily" view fetches calls from midnight (in
`DASHBOARD_TIMEZONE`) to now; the range helpers (`DayRange`/`WeekRange`/
`MonthRange`) are already in place for future weekly/monthly views.

### Excel spreadsheet

Point `SALES_XLSX_PATH` at your `.xlsx`. **Column names are configurable**, not
hard-coded — set them to match your sheet's header row (matching is case- and
whitespace-insensitive):

```env
SALES_XLSX_PATH=/app/data/sales-sample.xlsx
EXCEL_SHEET=                     # blank = first sheet
EXCEL_COL_NAME=Name              # required
EXCEL_COL_SALES_COUNT=Sales Count
EXCEL_COL_SALES_VALUE=Sales Value
EXCEL_COL_TARGET=Target          # optional quota
```

Only the name column is required. Currency symbols, thousands separators and
stray whitespace in number cells are tolerated (`£18,250` → `18250`). Rows
without a name are skipped.

See [`backend/data/sales-sample.xlsx`](./backend/data) for the expected format
(regenerate it any time with `go run ./cmd/gensample` from `backend/`).

You can also **replace the spreadsheet at runtime** without redeploying:

```bash
curl -F "file=@/path/to/sales.xlsx" http://localhost:8080/api/sales/upload
```

### Staff mapping

This file is the source of truth for who appears on the board and how they map
across sources. Edit [`backend/config/staff-mapping.json`](./backend/config/staff-mapping.example.json):

```json
{
  "staff": [
    {
      "id": "jane",
      "name": "Jane Doe",
      "email": "jane.doe@example.com",
      "aircallUserId": 1001,
      "spreadsheetName": "Jane Doe"
    }
  ]
}
```

- `aircallUserId` — the numeric Aircall user id (matches call data to a person).
- `spreadsheetName` — must equal the `Name` cell in the spreadsheet (defaults to
  `name` if omitted).
- A person with no matching Aircall user or no spreadsheet row still appears on
  the board with a clear **missing-data** indicator. Spreadsheet rows that match
  nobody are logged as unmatched.

### Scoring weights

Ranking uses a simple, configurable weighted model. Each metric is normalised
against the team's best performer in that metric, then combined and scaled to
0–100. Weights need not sum to 1 (they are normalised):

```env
WEIGHT_SALES=0.50
WEIGHT_CALLS_OUT=0.25
WEIGHT_CALLS_IN=0.10
WEIGHT_MINUTES=0.15
```

Raw metrics are always included in the API response, so the wallboard shows
actual numbers, not just scores.

---

## API

| Method & path             | Description                                                        |
| ------------------------- | ------------------------------------------------------------------ |
| `GET /api/dashboard`      | Full dashboard payload (people, totals, ranks, status, scoring).   |
| `GET /api/health`         | Service health (`ok` / `degraded` / `starting`), uptime, version.  |
| `GET /api/status`         | Per-source status, last successful refresh, errors, interval.      |
| `POST /api/sales/upload`  | Replace the sales spreadsheet (`multipart/form-data`, field `file`).|

The `/api/dashboard` response includes, per person: name, calls in/out, total
calls, minutes on phone, sales count/value, target & progress, score, rank, and
`hasAircallData` / `hasSalesData` flags — plus team totals, top performer,
biggest mover, source status, last-updated timestamp, and the active scoring
weights.

---

## Local development (without Docker)

**Prerequisites:** Go 1.24+, Node 20+.

**Backend:**

```bash
cd backend
go run ./cmd/gensample          # (re)generate the sample spreadsheet
go run ./cmd/server             # serves on :8080 (demo mode if no Aircall creds)
```

**Frontend:**

```bash
cd frontend
npm install
npm run dev                     # serves on :5173, proxies /api to :8080
```

Open http://localhost:5173.

---

## Testing

```bash
cd backend
go test ./...
```

Covered: Excel parsing (configurable columns, messy currency values, missing
columns/names), aggregation & ranking, team totals, missing-data flags, biggest
mover, and the Aircall client (live pagination/aggregation via a stub server,
auth-error handling, and deterministic demo data).

---

## Deploying to an office mini-PC / TV

A small always-on machine (Intel NUC, mini-PC, even a Raspberry Pi 4 with a
64-bit OS) connected to the TV is ideal.

1. **Install Docker** and clone this repo onto the machine.
2. `cd kpi-wallboard && cp .env.example .env`, then edit `.env` with real
   Aircall credentials and your spreadsheet/column settings. Drop your `.xlsx`
   into `backend/data/` (or upload it via the API) and edit
   `backend/config/staff-mapping.json`.
3. `docker compose up -d --build` — `restart: unless-stopped` brings both
   services back automatically after a reboot or power cut.
4. **Auto-launch the board in fullscreen.** Configure the OS to log in
   automatically and open a kiosk browser at boot, e.g.:
   ```bash
   chromium-browser --kiosk --noerrdialogs --disable-infobars \
     --incognito http://localhost:3000
   ```
   (On a Pi, add that line to autostart; on Windows, use Chrome/Edge
   `--kiosk` in a Startup shortcut.)
5. **Stop the screen sleeping** — disable display power management / screen
   blanking in the OS so the TV stays on.
6. The board auto-refreshes every 5 minutes and never needs manual interaction.
   Updating the spreadsheet (file or upload) or staff mapping is picked up on
   the next refresh — no rebuild required (config & data are volume-mounted).

**Tips**
- Hide the mouse cursor when idle (e.g. `unclutter` on Linux).
- For 4K TVs the layout uses relative units and scales cleanly; just run the
  browser at native resolution.
- Health-check the backend at `http://localhost:8080/api/health` if you add
  external monitoring.

---

## Project structure

```
kpi-wallboard/
├── docker-compose.yml          # full stack
├── .env.example                # annotated configuration
├── backend/                    # Go service
│   ├── cmd/server/             # main entrypoint
│   ├── cmd/gensample/          # sample-spreadsheet generator
│   ├── internal/
│   │   ├── config/             # env + staff-mapping loader
│   │   ├── models/             # shared data models / JSON contract
│   │   ├── aircall/            # Aircall client (live + demo)
│   │   ├── excel/              # configurable .xlsx parser
│   │   ├── scoring/            # weighted scoring
│   │   ├── aggregate/          # combine sources, rank, totals, movers
│   │   ├── store/              # thread-safe last-good cache + staleness
│   │   ├── refresh/            # scheduled pipeline
│   │   └── api/                # HTTP handlers
│   ├── config/staff-mapping.json
│   └── data/sales-sample.xlsx
└── frontend/                   # Svelte + D3 wallboard
    ├── src/lib/components/     # Header, Leaderboard, TopPerformer, Ticker, ...
    └── src/lib/charts/         # D3: LeaderboardBars, CallBreakdown, TargetProgress
```
