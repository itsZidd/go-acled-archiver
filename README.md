# ACLED Archive Pipeline

A high-performance Go-based data pipeline designed to ingest, transform, and archive global conflict data from the ACLED (Armed Conflict Location & Event Data Project) API into a local PostgreSQL database.

## 🛠 Tech Stack

- **Language:** Go (Golang)
- **Database:** PostgreSQL 16 (Docker)
- **Libraries:** `pgx/v5` (Batching/Upserting), `godotenv`
- **Environment:** Windows (Docker Desktop)

## 🚀 Features

- **OAuth2 Handshake:** Automated token rotation and refresh with proactive 5-minute buffer before expiry.
- **Stealth Ingestion:** Chrome-mimicking User-Agent headers and conservative 3s delay between pages to avoid Cloudflare flagging.
- **Null-Safe JSON Mapping:** Correctly handles ACLED's mixed-type fields including nullable `tags` (`*string`) and `population_best` (`*int`), preventing silent data corruption.
- **Atomic Upserts:** Uses PostgreSQL `ON CONFLICT` logic to ensure fatality counts and population data are updated for existing events without creating duplicates.
- **Scalable Batching:** Processes data in 5,000-row chunks using `pgx` batch sends to minimize database round-trips.
- **Per-Row Error Draining:** Batch results are drained individually so a single bad row surfaces an error immediately instead of silently failing.
- **CLI Year Range:** Sync any year range via `-start` and `-end` flags without touching the source code.

## ⚙️ Setup & Installation

### 1. Environment Configuration

Create a `.env` file in the root directory:

```bash
ACLED_EMAIL=your_email@uni.edu
ACLED_PASSWORD=your_password
DATABASE_URL=postgres://user:pass@localhost:5432/acled_archive?sslmode=disable
```

### 2. Start the Database

```bash
docker compose up -d
```

The database schema is applied automatically via `init.sql` on first container creation.

### 3. Build and Run

```bash
go mod tidy

# Sync a single year
go run ./cmd/sync-archive/main.go -start 2024 -end 2024

# Sync a range of years
go run ./cmd/sync-archive/main.go -start 2020 -end 2024
```

> ℹ️ Defaults to `-start 2023 -end 2023` if no flags are provided.

## 🧩 Project Architecture

```
go-acled-archiver/
├── cmd/sync-archive/
│   └── main.go              # Entry point — CLI flags, year/page loop, DB connection
├── internal/
│   ├── acled/
│   │   ├── client.go        # HTTP requests, OAuth2, JSON decoding
│   │   └── models.go        # Event & APIResponse structs with null-safe pointer types
│   └── database/
│       └── postgres.go      # High-speed batch upserts via pgx
├── migrations/
│   └── 001_init_schema.sql  # Reference schema (full 35-column version)
├── .env
├── go.mod
└── go.sum
```

## 📊 Lean Schema Design

The pipeline uses a trimmed schema optimised for **research and analytical use**, dropping fields with low analytical value (source, source_scale, admin2, admin3, location, geo_precision, timestamp, tags, population_1km, population_2km, population_5km) while retaining:

- Temporal: `event_date`, `year`
- Classification: `disorder_type`, `event_type`, `sub_event_type`
- Actor: `actor1`, `assoc_actor_1`, `inter1`, `actor2`, `assoc_actor_2`, `inter2`, `interaction`, `civilian_targeting`
- Geographic: `iso`, `region`, `country`, `admin1`, `latitude`, `longitude`
- Impact: `fatalities`, `notes`, `population_best`

## ⚠️ The "12-Month Wall" (Developer Note)

As of 2026, standard academic API keys are restricted to a 12-month rolling window.

- **Current Range:** April 2025 – April 2026.
- **Error Behavior:** Requesting data outside this range returns a `302 Redirect` to an HTML landing page, causing JSON parsing errors.
- **Solution:** Submit an "Academic Historical Access" request to ACLED to unlock historical data (2010–2024).

## 🛠 Troubleshooting

| Error | Cause | Fix |
|---|---|---|
| `invalid character '<' looking for beginning of value` | API redirected to HTML login/ToS page | Check `.env` credentials or accept new ToS via browser |
| `json: cannot unmarshal unquoted value into int` | ACLED changed a field from string to raw number | Remove `,string` tag from the relevant field in `models.go` |
| `UPDATE 0` on null-fix SQL | Field already stored correctly as NULL | No action needed |
| RAM spikes | 5,000 records per batch | Lower the `limit` parameter in `client.go` |
| Account banned / Cloudflare block | Aggressive retries | Pipeline intentionally breaks on first failure. Re-run manually after cooldown |
| `sslmode` connection error | Local Docker doesn't require SSL | Ensure `?sslmode=disable` in `DATABASE_URL` |
