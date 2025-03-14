package CORS

import (
	"github.com/rs/cors"
	"net/http"
)

// Settings add cors settings
func Settings() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:63342",
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		},
		AllowedHeaders: []string{ // получаем с фронта заголовок
			"Refresh-token",
			"Content-Type",
			"Authorization",
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
