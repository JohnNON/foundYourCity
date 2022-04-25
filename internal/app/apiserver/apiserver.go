package apiserver

import (
	"net/http"

	"github.com/JohnNON/foundYourCity/internal/store"
)

// Start - выполняет запуск сервера
func Start(config *Config, st *store.Store) error {
	srv := newServer(config, st)
	srv.configureRouter()

	return http.ListenAndServe(config.BindAddr, nil)
}
