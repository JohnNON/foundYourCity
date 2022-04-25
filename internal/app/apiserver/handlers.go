package apiserver

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/JohnNON/foundYourCity/internal/model"
)

func (s *server) response(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) handleSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lat := r.URL.Query().Get("lat")
		long := r.URL.Query().Get("long")
		rad := r.URL.Query().Get("rad")
		if lat == "" || long == "" || rad == "" {
			s.response(w, r, http.StatusNotFound, model.Cities{})
			return
		}

		latitute, err := strconv.ParseFloat(lat, 64)
		if err != nil {
			s.response(w, r, http.StatusNotFound, model.Cities{})
			return
		}

		longitute, err := strconv.ParseFloat(long, 64)
		if err != nil {
			s.response(w, r, http.StatusNotFound, model.Cities{})
			return
		}
		radius, err := strconv.ParseFloat(rad, 64)
		if err != nil {
			s.response(w, r, http.StatusNotFound, model.Cities{})
			return
		}

		k, err := s.store.FindCitiesByCoord(latitute, longitute, radius)
		if err != nil {
			log.Println(err)
			s.response(w, r, http.StatusNotFound, model.Cities{})
			return
		}

		c := model.Cities{
			Cities: s.store.GetCityByKeys(k),
		}

		s.response(w, r, http.StatusOK, &c)
	}
}

func (s *server) handleSearchCity() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			s.response(w, r, http.StatusNotFound, model.Cities{})
			return
		}

		k, err := s.store.FindCitiesByName(name)
		if err != nil {
			log.Println(err)
			s.response(w, r, http.StatusNotFound, model.Cities{})
			return
		}

		c := model.Cities{
			Cities: s.store.GetCityByKeys(k),
		}

		s.response(w, r, http.StatusOK, &c)
	}
}
