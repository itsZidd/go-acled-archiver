# Changelog

All notable changes to the WorldTension ACLED Archive Pipeline will be documented here.

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
