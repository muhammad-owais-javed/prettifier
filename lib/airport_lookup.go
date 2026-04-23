package lib

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

type AirportInfo struct {
	Name string
	City string
}

func LoadAirportData(path string) (map[string]AirportInfo, map[string]AirportInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Airport Lookup file not found: %w", err)
	}
	defer file.Close()

	iataMap := make(map[string]AirportInfo) // Key: "LAX", Value: {Name:"Los Angeles International Airport", City:"Los Angeles"}
	icaoMap := make(map[string]AirportInfo) // Key: "EGLL", Value: {Name: "London Heathrow Airport", City: "London"}

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("Could not read header from airport lookup: %w", err)
	}

	var nameIndex, iataIndex, icaoIndex, cityIndex int

	for i, column := range header {
		switch column {
		case "name":
			nameIndex = i
		case "iata_code":
			iataIndex = i
		case "icao_code":
			icaoIndex = i
		case "municipality":
			cityIndex = i
		}
	}

	for {
		data, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, nil, fmt.Errorf("Error reading airport-lookup record: %w", err)
		}

		if len(data) != len(header) {
			return nil, nil, fmt.Errorf("airport lookup malformed: incorrect number of columns")
		}

		airportName := data[nameIndex]
		iataCode := data[iataIndex]
		icaoCode := data[icaoIndex]
		cityName := data[cityIndex]

		info := AirportInfo{
			Name: airportName,
			City: cityName,
		}

		for i, col := range data {
			if col == "" {
				return nil, nil, fmt.Errorf("airport lookup malformed: blank data in column %d", i+1)
			}
		}

		if icaoCode != "" {
			icaoMap[icaoCode] = info
		}

		if iataCode != "" {
			iataMap[iataCode] = info
		}

	}

	return iataMap, icaoMap, nil
}
