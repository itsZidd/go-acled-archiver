# ACLED Archive Pipeline

A high-performance Go-based data pipeline designed to ingest, transform, and archive global conflict data from the ACLED (Armed Conflict Location & Event Data Project) API into a Neon/PostgreSQL database.

## 🛠 Tech Stack

- **Language:** Go (Golang)
- **Database:** Neon (Serverless PostgreSQL)
- **Libraries:** `pgx/v5` (Batching/Upserting), `godotenv`
- **Environment:** Arch Linux

## 🚀 Features

- **OAuth2 Handshake:** Automated token rotation and refresh with proactive 5-minute buffer before expiry.
- **Stealth Ingestion:** Chrome-mimicking User-Agent headers and conservative 3s delay between pages to avoid Cloudflare flagging.
- **Null-Safe JSON Mapping:** Correctly handles ACLED's mixed-type fields including nullable `tags` (`*string`) and `population_best` (`*int`), preventing silent data corruption.
- **Atomic Upserts:** Uses PostgreSQL `ON CONFLICT` logic to ensure fatality counts and population data are updated for existing events without creating duplicates.
- **Scalable Batching:** Processes data in 5,000-row chunks using `pgx` batch sends to minimize database round-trips.

## ⚙️ Setup & Installation

### 1. Environment Configuration

Create a `.env` file in the root directory:

```bash
ACLED_EMAIL=your_email@uni.edu
ACLED_PASSWORD=your_password
DATABASE_URL=postgres://user:pass@ep-host.neon.tech/neondb?sslmode=require
```

### 2. Database Migration

Execute the schema found in `migrations/001_init_schema.sql`.

**Note:** Interaction fields (`inter1`, `inter2`, `interaction`) are stored as TEXT because standard academic tiers return strings (e.g., "Protesters") rather than integer codes.

### 3. Build and Run

```bash
go mod tidy
go run cmd/sync-archive/main.go
```

### 4. Changing Target Years

Open `cmd/sync-archive/main.go` and edit the `years` slice:

```go
// Sync a single year
years := []int{2023}

// Sync multiple years
years := []int{2023, 2024, 2025}
```

Then run:

```bash
go run cmd/sync-archive/main.go
```

> ⚠️ The pipeline will break on first error per year to protect your ACLED account from Cloudflare flagging. If it stops mid-way, just re-run — `ON CONFLICT` upserts make it safe to resume from scratch.

## 🧩 Project Architecture

- `cmd/sync-archive/` — Entry point. Contains the main loop logic for year/page iteration.
- `internal/acled/`
  - `client.go` — Handles HTTP requests, OAuth2, and JSON decoding.
  - `models.go` — Defines the `Event` and `APIResponse` structs with null-safe pointer types.
- `internal/database/`
  - `postgres.go` — Manages high-speed batch upserts via `pgx`.
- `internal/migrations/`
  - `001_init_schema.sql` — Initial table schema.

## ⚠️ The "12-Month Wall" (Developer Note)

As of 2026, standard academic API keys are restricted to a 12-month rolling window.

- **Current Range:** April 2025 – April 2026.
- **Error Behavior:** Requesting data outside this range will return a `302 Redirect` to an HTML landing page, causing JSON parsing errors.
- **Solution:** To unlock historical data (2010–2024), submit an "Academic Historical Access" request to ACLED.

## 🛠 Troubleshooting "Future Me"

- **`invalid character '<' looking for beginning of value`** — The API redirected to an HTML login/ToS page. Check `.env` credentials or log in via browser to accept new Terms of Service.
- **`json: cannot unmarshal unquoted value into int`** — ACLED changed a field from string to raw number. Check `models.go` and remove the `,string` tag from the relevant field.
- **`UPDATE 0` on null-fix SQL** — Field was already stored correctly as NULL (e.g. `tags`). No action needed.
- **RAM spikes** — The pipeline handles 5,000 records per batch. If running on a low-resource VPS, lower the `limit` parameter in `client.go`.
- **Account banned / Cloudflare block** — Do not retry aggressively on errors. The pipeline intentionally breaks on first failure per year to protect the account. Re-run manually after a cooldown.
