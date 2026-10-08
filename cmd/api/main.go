package main

import (
	"context"
	"log"
	"os"
	"time"

	db "CloudBox/db/sqlc"
	"CloudBox/internal/database"
	"CloudBox/internal/user"

	"github.com/gin-gonic/gin"
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

	userRepository := user.NewRepository(queries)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"app":    "cloudbox",
			"status": "ok",
		})
	})

	api := router.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
	}

	log.Println("database connected successfully")
	log.Println("sqlc initialized successfully")
	log.Println("server running on :8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
