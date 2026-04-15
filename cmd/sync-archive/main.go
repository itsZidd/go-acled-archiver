package main

import (
	"context"
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
	_ = godotenv.Load()
	email := os.Getenv("ACLED_EMAIL")
	password := os.Getenv("ACLED_PASSWORD")
	dbURL := os.Getenv("DATABASE_URL")

	if email == "" || password == "" {
		log.Fatal("❌ ERROR: ACLED_EMAIL or PASSWORD missing in .env")
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

	years := []int{2023}

	for _, year := range years {
		page := 1
		for {
			fmt.Printf("🔍 [YEAR %d] Page %d: Authenticating & Fetching...\n", year, page)

			events, err := acledClient.FetchPage(year, page)
			if err != nil {
				log.Printf("❌ ACLED Error: %v", err)
				break
			}

			if len(events) == 0 {
				fmt.Printf("✅ Year %d sync complete.\n", year)
				break
			}

			fmt.Printf("📦 Received %d events. Writing to Railway...\n", len(events))

			startSave := time.Now()
			err = database.UpsertEvents(ctx, conn, events)
			if err != nil {
				log.Printf("❌ Database Error: %v", err)
				break
			}

			fmt.Printf("💾 Save successful (took %v). Waiting 3s...\n", time.Since(startSave).Truncate(time.Millisecond))
			time.Sleep(3 * time.Second)
			page++
		}
	}
	fmt.Println("🏁 THE ARCHIVE IS FULL.")
}
