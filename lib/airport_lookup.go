package lib

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type AirportInfo struct {
	Name string
	City string
}

func LoadAirportData(path string) (map[string]AirportInfo, map[string]AirportInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to open airport data file at %s: %w", path, err)
	}
	defer file.Close()

	iataMap := make(map[string]AirportInfo) // Key: "LAX", Value: {Name:"Los Angeles International Airport", City:"Los Angeles"}
	icaoMap := make(map[string]AirportInfo) // Key: "EGLL", Value: {Name: "London Heathrow Airport", City: "London"}

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	line := 1

	header, err := reader.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("Could not read header from airport data file: %w", err)
	}
	line++

	if len(header) != 6 {
		fmt.Println("Airport lookup file malformed")
		return nil, nil, fmt.Errorf("Number of Headers in Columns are not following Standard")
	}

	nameIndex, iataIndex, icaoIndex, cityIndex := -1, -1, -1, -1

	assignedHeaders := make(map[string]bool)

	for i, column := range header {

		if assignedHeaders[column] {
			fmt.Println("Airport lookup file malformed")
			return nil, nil, fmt.Errorf("Invalid header: Duplicate column name '%s' found", column)
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
		fmt.Println("Airport lookup malformed")
		return nil, nil, fmt.Errorf("Missing required Column in the header")
	}

	for {
		data, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, nil, fmt.Errorf("Error parsing airport data at line %d: %w", line, err)
		}

		if len(data) != len(header) {
			return nil, nil, fmt.Errorf("Column count mismatch at line %d: expected %d, got %d", line, len(header), len(data))
		}

		airportName := data[nameIndex]

		if airportName != "" {
			for _, char := range airportName {
				if char < 32 || char > 126 {
					fmt.Println("Airport lookup file malformed")
					return nil, nil, fmt.Errorf("Invalid character found in airport name '%s' at line %d", airportName, line)
				}
			}
		}

		icaoCode := data[icaoIndex]
		iataCode := data[iataIndex]
		cityName := data[cityIndex]

		if airportName == "" || icaoCode == "" || iataCode == "" || cityName == "" {
			return nil, nil, fmt.Errorf("Empty value for a required field at line %d", line)
		}

		info := AirportInfo{
			Name: airportName,
			City: cityName,
		}

		for _, col := range data {
			if col == "" {
				return nil, nil, fmt.Errorf("Airport lookup malformed")
			}
		}

		// Uncomment following lines to enable Hard Fail Test for ICAO and IATA code
		/*
			if icaoCode != "" {
				if len(icaoCode) != 4 {
					return nil, nil, fmt.Errorf("Invalid ICAO Code Length")
				}
				for _, char := range icaoCode {
					if char < 'A' || char > 'Z' {
						return nil, nil, fmt.Errorf("Invalid ICAO Code Character")
					}
				}
			}

			if iataCode != "" {
				if len(iataCode) != 3 {
					return nil, nil, fmt.Errorf("Invalid IATA Code Length")
				}
				for _, char := range iataCode {
					if char < 'A' || char > 'Z' {
						return nil, nil, fmt.Errorf("Invalid IATA Character")
					}
				}
			}
		*/

		if icaoCode != "" {
			icaoMap[icaoCode] = info
		}

		if iataCode != "" {
			iataMap[iataCode] = info
		}

		line++

	}

	return iataMap, icaoMap, nil
}
