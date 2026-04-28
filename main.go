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

	fileContent, error := os.ReadFile(inputPath)
	if error != nil {
		if os.IsNotExist(error) {
			//log.Fatal("Input not found")
			fmt.Println("Input not found.")
			return
		}
		log.Fatal(error)
	}

	iataLookup, icaoLookup, error := lib.LoadAirportData(lookupPath)
	if error != nil {
		if os.IsNotExist(error) {
			//log.Fatal("Airport lookup not found")
			fmt.Println("Airport lookup not found.")
			return
		}
		log.Fatal(error)
	}

	plainOutput, colorOutput := lib.DateTimeParsing(string(fileContent), iataLookup, icaoLookup)

	plainOutput = lib.TrimSpaces(plainOutput)

	error = os.WriteFile(outputPath, []byte(plainOutput), 0644)
	if error != nil {
		log.Fatal(error)
	}

	lib.StdOutput(colorOutput, outputPath)
	
}
