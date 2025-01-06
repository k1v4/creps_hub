package CORS

import (
	"github.com/rs/cors"
	"net/http"
)

func CorsSettings() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
		},
		AllowedHeaders: []string{ // получаем с фронта заголовок
			"Refresh-token",
			"Content-Type",
		},
		ExposedHeaders: []string{ // отдаём с фронта заголовок
			"Refresh-token",
		},
		AllowCredentials:    true,
		AllowPrivateNetwork: false,
		OptionsPassthrough:  false,
		Debug:               true,
	})
}
