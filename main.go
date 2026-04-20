package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {

	if len(os.Args) != 4 || os.Args[1] == "help" || os.Args[1] == "-h" {
		fmt.Println("itinerary usage:\n" +
			"go run . ./input.txt ./output.txt ./airport-lookup.csv")
		return
	}

	fmt.Println("Process Initializing...")

	inputPath := os.Args[1]

	// Read directly later to avoid TOCTOU race conditions and redundant Stat calls.
	/*
		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			fmt.Println("Input file not found!")
			return
		}
	*/
	outputPath := os.Args[2]
	lookupPath := os.Args[3]
	/*
		if _, err := os.Stat(lookupPath); os.IsNotExist(err) {
			fmt.Println("Airport lookup file not found!")
			return
		}
	*/

	iataMap, icaoMap, err := loadAirportData(lookupPath)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("IATA Map")
	// fmt.Println(iataMap)
	// fmt.Println("ICAO Map")
	// fmt.Print(icaoMap)

	content, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatal(err)
	}

	textContent := string(content)

	for code, airportName := range iataMap {
		match := "#" + code
		textContent = strings.ReplaceAll(textContent, match, airportName)
	}

	for code, airportName := range icaoMap {
		match := "##" + code
		textContent = strings.ReplaceAll(textContent, match, airportName)
	}

	processedContent := textContent

	err = os.WriteFile(outputPath, []byte(processedContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Process Completed...!")
}

func loadAirportData(path string) (map[string]string, map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Airport Lookup file not found: %w", err)
	}
	defer file.Close()

	iataMap := make(map[string]string) // Key: "LAX", Value: "Los Angeles International Airport"
	icaoMap := make(map[string]string) // Key: "EGLL", Value: "London Heathrow Airport"

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("Could not read header from airport lookup: %w", err)
	}

	var nameIndex int
	var iataIndex int
	var icaoIndex int

	for i, column := range header {
		switch column {
		case "name":
			nameIndex = i
		case "iata_code":
			iataIndex = i
		case "icao_code":
			icaoIndex = i
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

		for i, col := range data {
			if col == "" {
				return nil, nil, fmt.Errorf("airport lookup malformed: blank data in column %d", i+1)
			}
		}

		if icaoCode != "" {
			icaoMap[icaoCode] = airportName
		}

		if iataCode != "" {
			iataMap[iataCode] = airportName
		}

	}

	return iataMap, icaoMap, nil
}
