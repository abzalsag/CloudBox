package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `env:"APP_ENV" env-default:"dev"`
	HTTPAddr string `env:"HTTP_ADDR" env-default:":8080"`

	DatabaseURL string `env:"DATABASE_URL" env-required:"true"`

	JWTSecret  string        `env:"JWT_SECRET" env-required:"true"`
	AccessTTL  time.Duration `env:"ACCESS_TTL" env-default:"15m"`
	RefreshTTL time.Duration `env:"REFRESH_TTL" env-default:"720h"`

	S3Endpoint  string `env:"S3_ENDPOINT" env-default:"http://localhost:8333"`
	S3AccessKey string `env:"S3_ACCESS_KEY" env-default:"cloudbox"`
	S3SecretKey string `env:"S3_SECRET_KEY" env-default:"cloudbox-secret"`
	S3Bucket    string `env:"S3_BUCKET" env-default:"cloudbox"`
}

func Load() (*Config, error) {
	var cfg Config

	var err error
	if _, statErr := os.Stat(".env"); statErr == nil {
		err = cleanenv.ReadConfig(".env", &cfg)
	} else {
		err = cleanenv.ReadEnv(&cfg)
	}
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	return &cfg, nil
}
