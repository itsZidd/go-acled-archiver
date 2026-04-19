package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"acled-archiver/internal/acled"
	"acled-archiver/internal/database"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	startYear := flag.Int("start", 2023, "Start year to sync")
	endYear := flag.Int("end", 2023, "End year to sync")
	flag.Parse()

	if *startYear > *endYear {
		log.Fatal("❌ ERROR: -start year cannot be greater than -end year")
	}

	_ = godotenv.Load()

	email := os.Getenv("ACLED_EMAIL")
	password := os.Getenv("ACLED_PASSWORD")
	dbURL := os.Getenv("DATABASE_URL")

	if email == "" || password == "" {
		log.Fatal("❌ ERROR: ACLED_EMAIL or ACLED_PASSWORD missing in .env")
	}
	if dbURL == "" {
		log.Fatal("❌ ERROR: DATABASE_URL missing in .env")
	}

	acledClient := &acled.Client{
		Email:    email,
		Password: password,
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatal("❌ Database connection failed:", err)
	}
	defer conn.Close(ctx)

	fmt.Printf("🚀 Starting sync from %d to %d\n", *startYear, *endYear)

	totalEvents := 0

	for year := *startYear; year <= *endYear; year++ {
		page := 1
		yearEvents := 0

		for {
			fmt.Printf("🔍 [YEAR %d] Page %d: Fetching...\n", year, page)

			events, err := acledClient.FetchPage(year, page)
			if err != nil {
				log.Printf("❌ ACLED Error (year %d, page %d): %v", year, page, err)
				break
			}

			if len(events) == 0 {
				fmt.Printf("✅ Year %d complete — %d events synced.\n", year, yearEvents)
				break
			}

			fmt.Printf("📦 Received %d events. Writing to database...\n", len(events))
			startSave := time.Now()

			if err := database.UpsertEvents(ctx, conn, events); err != nil {
				log.Printf("❌ Database Error (year %d, page %d): %v", year, page, err)
				break
			}

			yearEvents += len(events)
			totalEvents += len(events)

			fmt.Printf("💾 Saved in %v. Waiting 3s...\n", time.Since(startSave).Truncate(time.Millisecond))
			time.Sleep(3 * time.Second)
			page++
		}
	}

	fmt.Printf("🏁 SYNC COMPLETE — %d total events archived.\n", totalEvents)
}
