package database

import (
	"acled-sync/internal/acled"
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context) (*pgx.Conn, error) {
	return pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
}

func UpsertEvents(ctx context.Context, conn *pgx.Conn, events []acled.Event) error {
	batch := &pgx.Batch{}
	sql := `
		INSERT INTO acled_events (
			event_id_cnty, event_date, year, disorder_type, event_type,
			sub_event_type, actor1, actor2, inter1, inter2,
			interaction, country, iso, fatalities, notes, timestamp
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		ON CONFLICT (event_id_cnty) DO UPDATE SET
			fatalities = EXCLUDED.fatalities,
			notes = EXCLUDED.notes,
			timestamp = EXCLUDED.timestamp,
			last_updated_at = CURRENT_TIMESTAMP;`

	for _, e := range events {
		batch.Queue(sql,
			e.EventID, e.EventDate, e.Year, e.DisorderType, e.EventType,
			e.SubEventType, e.Actor1, e.Actor2, e.Inter1, e.Inter2,
			e.Interaction, e.Country, e.ISO, e.Fatalities, e.Notes, e.Timestamp,
		)
	}

	br := conn.SendBatch(ctx, batch)
	return br.Close()
}
