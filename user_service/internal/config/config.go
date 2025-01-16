package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
	"user_service/pkg/DB/postgres"
)

type Config struct {
	postgres.DBConfig

	GRPCServerPort int           `env:"GRPC_SERVER_PORT" env-description:"grpc server port" env-default:"50052"`
	RestServerPort int           `env:"REST_SERVER_PORT" env-description:"rest server port" env-default:"8081"`
	TokenTTL       time.Duration `env:"TOKEN_TTL" env-default:"1h"`
}

func MustLoadConfig() *Config {
	//err := godotenv.Load(".env") // Явно указываем путь
	//if err != nil {
	//	panic(err)
	//}

	cfg := Config{}
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(nil)
	}

	return &cfg
}
