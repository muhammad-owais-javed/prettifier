package main

import (
	"fmt"
	"log"
	"os"
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

	//lookupPath := os.Args[3]

	/*
		if _, err := os.Stat(lookupPath); os.IsNotExist(err) {
			fmt.Println("Airport lookup file not found!")
			return
		}
	*/

	content, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatal(err)
	}

	processedContent := content

	err = os.WriteFile(outputPath, []byte(processedContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Process Completed...!")
}
