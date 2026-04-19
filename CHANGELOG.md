# Changelog
All notable changes to the ACLED Archive Pipeline will be documented here.

---

## [v0.3.0] - 2026-04-19

### Changed
- Migrated database from Railway (cloud) to local Docker PostgreSQL 16
- `DATABASE_URL` now points to `localhost:5432` with `sslmode=disable`
- Trimmed schema from 35 columns down to 24 — dropped `time_precision`, `geo_precision`, `timestamp`, `tags`, `source`, `source_scale`, `admin2`, `admin3`, `location`, `population_1km`, `population_2km`, `population_5km` in favour of a lean research-focused schema
- `inter1` and `inter2` corrected from `INT` to `TEXT` in schema — ACLED academic tier returns string codes not integers
- Added `region` and `sub_event_type` retained as analytically valuable fields
- Kept grey-area fields: `assoc_actor_1`, `assoc_actor_2`, `interaction`, `civilian_targeting`

### Added
- CLI `-start` and `-end` year flags — no longer need to edit source code to change sync range
- Input validation: exits early if `-start` > `-end`
- `DATABASE_URL` empty check added alongside existing credential checks in `main.go`
- Per-year and total event counters in sync output (`yearEvents`, `totalEvents`)
- Per-row batch error draining in `postgres.go` — replaces silent `br.Close()` swallow
- `civilian_targeting` added to `ON CONFLICT DO UPDATE` fields since ACLED revises this
- `last_updated_at` column added to schema (required by `ON CONFLICT` clause)
- Additional indexes: `idx_acled_year`, `idx_acled_event_type`, `idx_acled_actor1`, `idx_acled_region`

### Fixed
- `PopulationBest` changed from `int` to `*int` in `models.go` — prevents `null` being stored as `0`
- Error logs now include year and page number for easier debugging

---

## [v0.2.0] - 2026-04-15

### Fixed
- `Tags` field in `models.go` changed from `string` to `*string` to correctly handle JSON `null` values from ACLED API
- `PopulationBest` field in `models.go` changed from `int` to `*int` to correctly handle JSON `null` values — previously stored as `0`, corrupting population data for ~292k rows
- Database backfill: ran `UPDATE acled_events SET population_best = NULL WHERE population_best = 0` to correct existing corrupted rows (292,361 rows fixed)
- Fixed ignored error return from `http.NewRequest` in `client.go` `FetchPage()` — now properly propagates error

### Notes
- `tags` field confirmed to already be stored as NULL correctly in PostgreSQL — no DB fix needed for that field
- Total dataset as of this release: 394,995 rows (2024, complete) + 118,093 rows (2025, ongoing)

---

## [v0.1.0] - Initial Release

### Added
- OAuth2 authentication with proactive token refresh (5-minute buffer before expiry)
- Page-based sync loop for ACLED API with configurable year targets
- `pgx` batch upserts with `ON CONFLICT (event_id_cnty) DO UPDATE` for idempotent re-runs
- Null-safe handling for ACLED's mixed-type fields (`latitude`, `longitude` as `,string` tagged floats)
- `.env` support via `godotenv` for credentials and database URL
- Initial PostgreSQL schema (`001_init_schema.sql`) with 35 ACLED fields
- Conservative 3-second delay between pages to avoid Cloudflare rate limiting
- Break-on-error strategy (intentional) to protect API account from being flagged
