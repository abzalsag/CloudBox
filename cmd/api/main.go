package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"CloudBox/internal/config"
	httpdelivery "CloudBox/internal/delivery/http"
	"CloudBox/internal/delivery/http/handler"
	"CloudBox/internal/repository/postgres"
	"CloudBox/internal/repository/s3"
	"CloudBox/internal/usecase/auth"
	"CloudBox/internal/usecase/file"
	"CloudBox/internal/usecase/folder"
	"CloudBox/pkg/hasher"
	"CloudBox/pkg/logger"
	"CloudBox/pkg/token"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		l.Fatal("db connect", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		l.Fatal("db ping", zap.Error(err))
	}
	l.Info("database connected")

	storage, err := s3.New(ctx, s3.Config{
		Endpoint:  cfg.S3Endpoint,
		AccessKey: cfg.S3AccessKey,
		SecretKey: cfg.S3SecretKey,
		Bucket:    cfg.S3Bucket,
	})
	if err != nil {
		l.Fatal("s3 init", zap.Error(err))
	}
	l.Info("s3 connected", zap.String("bucket", storage.Bucket()))

	store := postgres.NewStore(pool)
	userRepo := postgres.NewUserRepo(store)
	refreshRepo := postgres.NewRefreshTokenRepo(store)
	folderRepo := postgres.NewFolderRepo(store)
	fileRepo := postgres.NewFileRepo(store)

	tokenMgr := token.NewManager(cfg.JWTSecret, cfg.AccessTTL)
	authUC := auth.NewUsecase(userRepo, refreshRepo, store, hasher.Bcrypt{}, tokenMgr, cfg.AccessTTL, cfg.RefreshTTL)
	folderUC := folder.NewUsecase(folderRepo)
	fileUC := file.NewUsecase(fileRepo, folderRepo, userRepo, storage, store)

	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := httpdelivery.NewRouter(
		l,
		handler.NewAuthHandler(authUC),
		handler.NewFolderHandler(folderUC),
		handler.NewFileHandler(fileUC),
		tokenMgr,
	)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		l.Info("http server started", zap.String("addr", cfg.HTTPAddr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Fatal("http server", zap.Error(err))
		}
	}()

	<-ctx.Done()
	l.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		l.Error("shutdown", zap.Error(err))
	}
}
