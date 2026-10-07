package main

import (
	"context"
	"log"
	"os"
	"time"

	db "CloudBox/db/sqlc"
	"CloudBox/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := database.NewPostgres(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("database connection failed:", err)
	}

	queries := db.New(pool)

	_ = queries

	log.Println("database connected successfully")
	log.Println("sqlc initialized successfully")
}
