package config

import (
	"auth_service/pkg/DataBase/postgres"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	postgres.DBConfig

	GRPCServerPort int `env:"GRPC_SERVER_PORT" env-description:"grpc server port" env-default:"50051"`
	RestServerPort int `env:"REST_SERVER_PORT" env-description:"rest server port" env-default:"8080"`
}

func MustLoadConfig() *Config {
	//err := godotenv.Load(".env") // Явно указываем путь
	//if err != nil {
	//	panic(err)
	//}

	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}
