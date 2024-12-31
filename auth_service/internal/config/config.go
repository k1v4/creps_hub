package config

import "auth_service/pkg/DB/postgres"

type Config struct {
	postgres.Config

	GRPCServerPort int `env:"GRPC_SERVER_PORT"`
	RestServerPort int `env:"REST_SERVER_PORT"`
}
