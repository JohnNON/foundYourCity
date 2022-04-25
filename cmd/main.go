package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/JohnNON/foundYourCity/internal/app/apiserver"
	"github.com/JohnNON/foundYourCity/internal/store"
	"github.com/JohnNON/foundYourCity/internal/utils"
)

func init() {
	err := os.Mkdir(utils.Download, 0777)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}
	err = os.Mkdir(utils.Data, 0777)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}
}

func main() {
	config := apiserver.NewConfig()

	if os.Getenv("PORT") != "" {
		config.BindAddr = ":" + os.Getenv("PORT")
	}

	if os.Getenv("ENDPOINT") != "" {
		config.EndPoint = os.Getenv("ENDPOINT")
	}

	if os.Getenv("SOURCE") != "" {
		config.Source = os.Getenv("SOURCE")
	}

	cities, citiesSearchList, citiesSearchCoord := utils.GetCitiesList(config.Source)

	st := store.NewStore(cities, citiesSearchList, citiesSearchCoord)

	fmt.Printf("Listen %s\n", config.BindAddr)
	if err := apiserver.Start(config, st); err != nil {
		log.Fatal(err)
	}
}
