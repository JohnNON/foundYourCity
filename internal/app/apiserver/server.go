package apiserver

import (
	"fmt"
	"net/http"

	"github.com/JohnNON/foundYourCity/internal/store"
)

type server struct {
	endPoint string
	store    *store.Store
}

func newServer(config *Config, st *store.Store) *server {
	s := &server{
		endPoint: config.EndPoint,
		store:    st,
	}

	return s
}

func (s *server) configureRouter() {
	http.HandleFunc(s.endPoint, s.handleSearch())
	http.HandleFunc(fmt.Sprintf("%s/city", s.endPoint), s.handleSearchCity())
}
