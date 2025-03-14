package config

import (
	"article_service/pkg/DataBase/postgres"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	postgres.DBConfig

	RestServerPort int `env:"REST_SERVER_PORT" env-description:"rest server port" env-default:"8080"`
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
