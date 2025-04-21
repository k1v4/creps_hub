package config

import (
	"article_service/pkg/DataBase/postgres"
	"article_service/pkg/DataBase/redis"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	postgres.DBConfig
	redis.RedisConfig

	RestServerPort         int `env:"REST_SERVER_PORT" env-description:"rest server port" env-default:"8080"`
	UploaderGRPCServerPort int `env:"UPLOADER_GRPC_SERVER_PORT" env-default:"50053"`
}

func MustLoadConfig() *Config {
	//err := godotenv.Load(".env") // Явно указываем путь
	//if err != nil {
	//	panic(err)
	//}

	cfg := Config{}
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}
