package main

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"CloudBox/internal/config"
	"CloudBox/internal/entity"
	"CloudBox/internal/repository/postgres"
	"CloudBox/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	l, err := logger.New(cfg.Env)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Sync()

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		l.Fatal("db connect", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		l.Fatal("db ping", zap.Error(err))
	}
	l.Info("database connected")

	store := postgres.NewStore(pool)
	userRepo := postgres.NewUserRepo(store)

	_, err = userRepo.GetByEmail(ctx, "nobody@example.com")
	if errors.Is(err, entity.ErrNotFound) {
		l.Info("user repo check ok: ErrNotFound returned as expected")
	} else {
		l.Fatal("user repo check failed", zap.Error(err))
	}
}
