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
		event_id_cnty, event_date, year, time_precision, disorder_type,
		event_type, sub_event_type, actor1, assoc_actor_1, inter1,
		actor2, assoc_actor_2, inter2, interaction, civilian_targeting,
		iso, region, country, admin1, admin2,
		admin3, location, geo_precision, latitude, longitude,
		source, source_scale, notes, fatalities, timestamp,
		tags, population_1km, population_2km, population_5km, population_best
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
		$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
		$21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
		$31, $32, $33, $34, $35
	)
	ON CONFLICT (event_id_cnty) DO UPDATE SET
		fatalities = EXCLUDED.fatalities,
		population_best = EXCLUDED.population_best,
		last_updated_at = CURRENT_TIMESTAMP;`

	for _, e := range events {
		batch.Queue(sql,
			e.EventID, e.EventDate, e.Year, e.TimePrecision, e.DisorderType,
			e.EventType, e.SubEventType, e.Actor1, e.AssocActor1, e.Inter1,
			e.Actor2, e.AssocActor2, e.Inter2, e.Interaction, e.CivilianTargeting,
			e.ISO, e.Region, e.Country, e.Admin1, e.Admin2,
			e.Admin3, e.Location, e.GeoPrecision, e.Latitude, e.Longitude,
			e.Source, e.SourceScale, e.Notes, e.Fatalities, e.Timestamp,
			e.Tags,
			e.Population1km, e.Population2km, e.Population5km, e.PopulationBest,
		)
	}

	br := conn.SendBatch(ctx, batch)
	return br.Close()
}
