package utils

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/JohnNON/foundYourCity/internal/model"
)

const (
	Download = "download"
	Data     = "data"
)

func GetCitiesList(s string) (map[string]model.City, string, []model.CityPos) {
	cities, citiesSearchList, citiesSearchCoord := getCitiesList(s)

	return cities, citiesSearchList, citiesSearchCoord
}

func getCitiesList(s string) (map[string]model.City, string, []model.CityPos) {
	t := time.Now()
	httpClient := http.Client{}

	resp, err := httpClient.Get(
		fmt.Sprintf("https://download.geonames.org/export/dump/%s", s),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", nil
	}

	f, err := os.OpenFile(fmt.Sprintf("%s/%s", Download, s), os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	io.Copy(f, resp.Body)
	err = unzip(fmt.Sprintf("%s/%s", Download, s), Data)
	if err != nil {
		log.Fatal(err)
	}

	cities, citiesSearchList, citiesSearchCoord, err := readFile(fmt.Sprintf("%s/%s", Data, strings.ReplaceAll(s, ".zip", ".txt")))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s, " - ", time.Since(t).Seconds())

	return cities, citiesSearchList, citiesSearchCoord
}

func unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0777); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

func readFile(file string) (map[string]model.City, string, []model.CityPos, error) {
	cities := make(map[string]model.City)
	citiesSearchCoord := make([]model.CityPos, 0, 25000)
	var citiesSearchList strings.Builder
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer f.Close()

	fileScanner := bufio.NewScanner(f)

	for fileScanner.Scan() {
		rows := strings.Split(fileScanner.Text(), "\t")
		if len(rows) >= 15 {
			cities[rows[0]] = model.City{
				Name:       rows[1],
				Latitude:   rows[4],
				Longitude:  rows[5],
				Population: rows[14],
			}
			fmt.Fprintf(&citiesSearchList, "%s|%s|%s\n", rows[0], rows[1], rows[2])

			latitute, err := strconv.ParseFloat(rows[4], 64)
			if err != nil {
				log.Fatalf("Error converting %s: %s", rows[4], err)
			}

			longitute, err := strconv.ParseFloat(rows[5], 64)
			if err != nil {
				log.Fatalf("Error converting %s: %s", rows[5], err)
			}

			citiesSearchCoord = append(citiesSearchCoord, model.CityPos{ID: rows[0], Latitude: latitute, Longitude: longitute})
		}
	}

	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	sort.Slice(citiesSearchCoord, func(i, j int) bool {
		return citiesSearchCoord[i].Latitude < citiesSearchCoord[j].Latitude
	})

	return cities, citiesSearchList.String(), citiesSearchCoord, nil
}
