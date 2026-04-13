package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"

	"acled-sync/internal/acled"
	"acled-sync/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())

	client := &acled.Client{
		Email:    os.Getenv("ACLED_EMAIL"),
		Password: os.Getenv("ACLED_PASSWORD"),
	}

	conn, err := database.Connect(ctx)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer conn.Close(ctx)

	// STARTING FROM 2024 because of the 12-month account limit
	for year := 2024; year <= 2026; year++ {
		page := 1
		for {
			log.Printf("Checking Year %d | Page %d", year, page)
			events, err := client.FetchPage(year, page)
			if err != nil {
				log.Printf("!!! ERROR: %v", err)
				break
			}

			if len(events) == 0 {
				log.Printf("No data for %d.", year)
				break
			}

			log.Printf("Saving %d events to Neon...", len(events))
			if err := database.UpsertEvents(ctx, conn, events); err != nil {
				log.Printf("DB Error: %v", err)
			}

			if len(events) < 5000 {
				break
			}

			jitter := time.Duration(3000+rand.Intn(4000)) * time.Millisecond
			time.Sleep(jitter)
			page++
		}
	}
}
