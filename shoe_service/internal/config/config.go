package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"shoe_service/pkg/DB/postgres"
	"shoe_service/pkg/DB/redis"
	"time"
)

type Config struct {
	postgres.DBConfig
	redis.RedisConfig

	GRPCServerPort         int           `env:"GRPC_SERVER_PORT" env-description:"grpc server port" env-default:"50052"`
	RestServerPort         int           `env:"REST_SERVER_PORT" env-description:"rest server port" env-default:"8081"`
	TokenTTL               time.Duration `env:"TOKEN_TTL" env-default:"1h"`
	RefreshTokenTTL        time.Duration `env:"REFRESH_TOKEN_TTL" env-default:"24h"`
	UploaderGRPCServerPort int           `env:"UPLOADER_GRPC_SERVER_PORT" env-default:"50053"`
}

func MustLoadConfig() *Config {
	err := godotenv.Load(".env") // Явно указываем путь
	if err != nil {
		panic(err)
	}

	cfg := Config{}
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}
