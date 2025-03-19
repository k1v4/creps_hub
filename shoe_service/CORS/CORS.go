package CORS

import (
	"github.com/rs/cors"
	"net/http"
)

// Settings add cors settings
func Settings() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:63342",
			"http://localhost:63342/creps_hub/index.html?_ijt=ar5hg8rk35foclss6l26fbfcrg&_ij_reload=RELOAD_ON_SAVE",
			"http://localhost:3000",
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
