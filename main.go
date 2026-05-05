package main

import (
	"fmt"
	"log"
	"os"
	"prettifier/lib"
)

func main() {

	if len(os.Args) != 4 || os.Args[1] == "help" || os.Args[1] == "-h" {
		fmt.Println("itinerary usage:\n" +
			"go run . ./input.txt ./output.txt ./airport-lookup.csv")
		return
	}

	inputPath, outputPath, lookupPath := os.Args[1], os.Args[2], os.Args[3]

	fileContent, err := os.ReadFile(inputPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Input not found")
		}
		log.Fatal(err)
	}

	iataLookup, icaoLookup, err := lib.LoadAirportData(lookupPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Airport lookup not found")
		}
		log.Fatal(err)
	}

	plainOutput, colorOutput := lib.AirportCodesAndDateTimeParsing(string(fileContent), iataLookup, icaoLookup)

	plainOutput = lib.TrimSpaces(plainOutput)

	err = os.WriteFile(outputPath, []byte(plainOutput), 0644)
	if err != nil {
		log.Fatal(err)
	}

	lib.StdOutput(colorOutput, outputPath)

}
