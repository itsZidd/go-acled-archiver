# WorldTension: ACLED Archive Pipeline

A high-performance Go-based data pipeline designed to ingest, transform, and archive global conflict data from the ACLED (Armed Conflict Location & Event Data Project) API into a Neon/PostgreSQL database.

## 🛠 Tech Stack

* Language: Go (Golang)
* Database: Neon (Serverless PostgreSQL)
* Libraries: `pgx/v5` (Batching/Upserting), `godotenv`
* Environment: Arch Linux

## 🚀 Features

* OAuth2 Handshake: Automated token rotation and refresh for 24-hour sessions.
* Stealth Ingestion: Implements randomized jitter (3s-7s) and Chrome-mimicking headers to bypass Cloudflare and rate-limiters.
* JSON-to-Relational Mapping: Handles complex ACLED data types (including the "String-Interaction" mismatch discovered during development).
* Atomic Upserts: Uses PostgreSQL `ON CONFLICT` logic to ensure fatality counts and notes are updated for existing events without creating duplicates.
* Scalable Batching: Processes data in 5,000-row chunks to minimize database round-trips.

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

## 🧩 Project Architecture

* `cmd/sync-archive/`: Entry point. Contains the main loop logic for year/page iteration.
* `internal/acled/`:
  * `client.go`: Handles HTTP requests, OAuth, and JSON decoding.
  * `models.go`: Defines the `Event` and `APIResponse` structs.
* `internal/database/`:
  * `postgres.go`: Manages connection pooling and high-speed batch upserts.

## ⚠️ The "12-Month Wall" (Developer Note)

As of 2026, standard academic API keys are restricted to a 12-month rolling window.

* **Current Range:** April 2025 – April 2026.
* **Error Behavior:** Requesting data outside this range (e.g., year 2010) will return a `302 Redirect` to an HTML landing page, causing JSON parsing errors.
* **Solution:** To unlock 2010–2024, an "Academic Historical Access" request must be submitted to ACLED.

## 🛠 Troubleshooting "Future Me"

* **`invalid character '<' looking for beginning of value`:** This means the API redirected you to an HTML login/ToS page. Check your `.env` credentials or log in via a browser to accept new Terms of Service.
* **`json: cannot unmarshal unquoted value into int`:** ACLED changed a field from a string to a raw number. Check `models.go` and remove the `,string` tag.
* **RAM spikes:** The pipeline handles 5,000 records per batch. If running on a low-resource VPS, lower the `limit` parameter in `client.go`.
