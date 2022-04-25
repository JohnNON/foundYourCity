package model

type City struct {
	Name       string `json:"name"`
	Latitude   string `json:"lat"`
	Longitude  string `json:"long"`
	Population string `json:"population"`
}

type Cities struct {
	Cities []City `json:"cities"`
}

type CityPos struct {
	ID        string
	Latitude  float64
	Longitude float64
}
