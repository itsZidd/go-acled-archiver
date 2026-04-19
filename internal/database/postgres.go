package database

import (
	"acled-archiver/internal/acled"
	"context"

	"github.com/jackc/pgx/v5"
)

func UpsertEvents(ctx context.Context, conn *pgx.Conn, events []acled.Event) error {
	batch := &pgx.Batch{}

	sql := `
	INSERT INTO acled_events (
		event_id_cnty, event_date, year, disorder_type,
		event_type, sub_event_type,
		actor1, assoc_actor_1, inter1,
		actor2, assoc_actor_2, inter2,
		interaction, civilian_targeting,
		iso, region, country, admin1,
		latitude, longitude,
		fatalities, notes, population_best
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
		$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
		$21, $22, $23
	)
	ON CONFLICT (event_id_cnty) DO UPDATE SET
		fatalities       = EXCLUDED.fatalities,
		population_best  = EXCLUDED.population_best,
		civilian_targeting = EXCLUDED.civilian_targeting,
		last_updated_at  = CURRENT_TIMESTAMP;`

	for _, e := range events {
		batch.Queue(sql,
			e.EventID, e.EventDate, e.Year, e.DisorderType,
			e.EventType, e.SubEventType,
			e.Actor1, e.AssocActor1, e.Inter1,
			e.Actor2, e.AssocActor2, e.Inter2,
			e.Interaction, e.CivilianTargeting,
			e.ISO, e.Region, e.Country, e.Admin1,
			e.Latitude, e.Longitude,
			e.Fatalities, e.Notes, e.PopulationBest,
		)
	}

	br := conn.SendBatch(ctx, batch)
	defer br.Close()

	// Drain results to catch per-row errors
	for range events {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}

	return nil
}
