package lib

import (
	"encoding/csv"
	"fmt"
	"io"
	//"log"
	"os"
)

type AirportInfo struct {
	Name string
	City string
}

func LoadAirportData(path string) (map[string]AirportInfo, map[string]AirportInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		//fmt.Println("Airport lookup not found.")
		return nil, nil, err
		//log.Fatalf("Airport Lookup file not found: %w", err)
	}
	defer file.Close()

	iataMap := make(map[string]AirportInfo) // Key: "LAX", Value: {Name:"Los Angeles International Airport", City:"Los Angeles"}
	icaoMap := make(map[string]AirportInfo) // Key: "EGLL", Value: {Name: "London Heathrow Airport", City: "London"}

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		fmt.Println("Airport lookup malformed")
		return nil, nil, fmt.Errorf("Could not read header from airport lookup: %w", err)
	}

	if len(header) != 6 {
		return nil, nil, fmt.Errorf("Airport lookup malformed")
	}

	nameIndex, iataIndex, icaoIndex, cityIndex := -1, -1, -1, -1

	assignedHeaders := make(map[string]bool)

	for i, column := range header {

		if assignedHeaders[column] {
			return nil, nil, fmt.Errorf("Airport lookup malformed")
		}

		switch column {
		case "name":
			nameIndex = i
			assignedHeaders[column] = true
		case "icao_code":
			icaoIndex = i
			assignedHeaders[column] = true
		case "iata_code":
			iataIndex = i
			assignedHeaders[column] = true
		case "municipality":
			cityIndex = i
			assignedHeaders[column] = true
		}
	}

	if nameIndex == -1 || iataIndex == -1 || icaoIndex == -1 || cityIndex == -1 {
    	//fmt.Println("Airport lookup malformed")
		return nil, nil, fmt.Errorf("Airport lookup malformed")
	}


	for {
		data, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, nil, fmt.Errorf("Airport lookup malformed") 
		}

		if len(data) != len(header) {
			return nil, nil, fmt.Errorf("Airport lookup malformed") 
		}

		airportName := data[nameIndex]
		icaoCode := data[icaoIndex]
		iataCode := data[iataIndex]
		cityName := data[cityIndex]

		info := AirportInfo{
			Name: airportName,
			City: cityName,
		}

		for _, col := range data {
			if col == "" {
				return nil, nil, fmt.Errorf("Airport lookup malformed") 
			}
		}

		if icaoCode != "" {
			if len(icaoCode) != 4 {
				return nil, nil, fmt.Errorf("Airport lookup malformed")
			}
			for _, char := range icaoCode {
				if char < 'A' || char > 'Z' {
					return nil, nil, fmt.Errorf("Airport lookup malformed")
				}
			}
		}

		if iataCode != "" {
			if len(iataCode) != 3 {
				return nil, nil, fmt.Errorf("Airport lookup malformed")
			}
			for _, char := range iataCode {
				if char < 'A' || char > 'Z' {
					return nil, nil, fmt.Errorf("Airport lookup malformed")
				}
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
