package store

import (
	"regexp"
	"sort"
	"sync"

	"github.com/JohnNON/foundYourCity/internal/model"
	"github.com/JohnNON/foundYourCity/internal/utils"
)

type Store struct {
	mu                sync.Mutex
	cities            map[string]model.City
	citiesSearchList  string
	citiesSearchCoord []model.CityPos
}

func (s *Store) GetCityByKey(key string) model.City {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.cities[key]
}

func (s *Store) GetCityByKeys(keys []string) []model.City {
	cities := make([]model.City, 0, len(keys))

	for _, key := range keys {
		cities = append(cities, s.GetCityByKey(key))
	}

	return cities
}

func (s *Store) FindCitiesByName(name string) ([]string, error) {
	re, err := regexp.Compile(`(?i)([0-9]*)\|(.*` + name + `)`)
	if err != nil {
		return nil, err
	}

	match := re.FindAllStringSubmatch(s.citiesSearchList, -1)
	if len(match) == 0 {
		return nil, nil
	}

	keys := make([]string, 0, len(match))
	for _, vals := range match {
		if len(vals) < 2 {
			continue
		}
		keys = append(keys, vals[1])
	}

	return keys, nil
}

func (s *Store) FindCitiesByCoord(lat, long, rad float64) ([]string, error) {
	r := rad / 111000
	posL := sort.Search(len(s.citiesSearchCoord), func(i int) bool { return s.citiesSearchCoord[i].Latitude >= lat-r })
	posR := sort.Search(len(s.citiesSearchCoord), func(i int) bool { return s.citiesSearchCoord[i].Latitude <= lat+r })

	if posL > posR {
		return nil, nil
	}

	coords := s.citiesSearchCoord[posL : posR+1]
	keys := make([]string, 0, len(coords))
	for _, coord := range coords {
		if utils.Distance(lat, long, coord.Latitude, coord.Longitude) > r {
			continue
		}
		keys = append(keys, coord.ID)
	}

	return keys, nil
}

func NewStore(cities map[string]model.City, citiesSearchList string, citiesSearchCoord []model.CityPos) *Store {
	return &Store{
		cities:            cities,
		citiesSearchList:  citiesSearchList,
		citiesSearchCoord: citiesSearchCoord,
	}
}
