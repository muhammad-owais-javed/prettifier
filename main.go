package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
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

	reg := regexp.MustCompile(`(D|T12|T24)\((.*?)\)`)
	matches := reg.FindAllStringSubmatch(textContent, -1)
	//fmt.Println(matches)

	monthMap := map[string]string{
		"01": "Jan", "02": "Feb", "03": "Mar", "04": "Apr", "05": "May", "06": "Jun",
		"07": "Jul", "08": "Aug", "09": "Sep", "10": "Oct", "11": "Nov", "12": "Dec",
	}

	for _, match := range matches {

		defaultDate := match[0]
		dateTag := match[1]
		//fmt.Println("dateTag: " + dateTag)
		isoDate := match[2]

		// To check if the format is not matched, skip it
		if !strings.Contains(isoDate, "T") || !strings.Contains(isoDate, "+") && !strings.Contains(isoDate, "-") && !strings.Contains(isoDate, "Z") {
			continue
		}

		parts := strings.Split(isoDate, "T")
		//fmt.Println(parts)

		if len(parts) != 2 {
			continue
		}

		date := parts[0]
		//fmt.Println("Date: " + date)
		timeWoffset := parts[1]
		//fmt.Println("Time & Offset: " + timeWoffset)

		dateSplit := strings.Split(date, "-")
		year := dateSplit[0]
		month := dateSplit[1]
		day := dateSplit[2]

		//fmt.Println("Year: " + year + " Month: " + month + " Day: " + day)

		var time string
		var offset string

		if strings.HasSuffix(timeWoffset, "Z") {
			time = strings.TrimSuffix(timeWoffset, "Z")
			offset = "+00:00"
			//fmt.Println("Time: " + time + " Offset: " + offset)
		} else if strings.Contains(timeWoffset, "+") {
			toffSplit := strings.Split(timeWoffset, "+")
			time = toffSplit[0]
			offset = "+" + toffSplit[1]
			//fmt.Println("Time: " + time + " Offset: " + offset)
		} else if strings.Contains(timeWoffset, "-") {
			toffSplit := strings.Split(timeWoffset, "-")
			time = toffSplit[0]
			offset = "-" + toffSplit[1]
			//fmt.Println("Time: " + time + " Offset: " + offset)
		} else {
			//fmt.Println("No Offset Found")
			continue
		}

		timeSplit := strings.Split(time, ":")
		hoursStr := timeSplit[0]
		minutes := timeSplit[1]
		//seconds := timeSplit[2]

		//fmt.Println("hoursStr: " + hoursStr + " Minutes: " + minutes + " Seconds: " + seconds)

		var formatResult string

		switch dateTag {
		case "D":
			//fmt.Println("Case D")
			monthName, _ := monthMap[month]
			formatResult = fmt.Sprintf("%s-%s-%s", day, monthName, year)

		case "T12":
			hours, err := strconv.Atoi(hoursStr)
			if err != nil {
				continue
			}
			AMPM := "AM"
			if hours >= 12 {
				AMPM = "PM"
			}
			if hours > 12 {
				hours = hours - 12
			}
			if hours == 0 {
				hours = 12
			}
			formatResult = fmt.Sprintf("%02d:%s:%s (%s)", hours, minutes, AMPM, offset)

		case "T24":
			//fmt.Println("Case T24")
			formatResult = fmt.Sprintf("%s:%s (%s)", hoursStr, minutes, offset)
			//fmt.Println("Case T24 formatResult: " + formatResult)
		}

		textContent = strings.Replace(textContent, defaultDate, formatResult, 1)
		//	fmt.Println("Text Content: " + textContent)
	}

	processedContent := textContent

	processedContent = strings.ReplaceAll(processedContent, "\r\n", "\n") // Handle Windows line endings first
	processedContent = strings.ReplaceAll(processedContent, "\r", "\n")
	processedContent = strings.ReplaceAll(processedContent, "\v", "\n")
	processedContent = strings.ReplaceAll(processedContent, "\f", "\n")

	for strings.Contains(processedContent, "\n\n\n") {
		processedContent = strings.ReplaceAll(processedContent, "\n\n\n", "\n\n")
	}

	processedContent = strings.TrimSpace(processedContent)

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
