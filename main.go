package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) != 4 || os.Args[1] == "help" || os.Args[1] == "-h" {
		fmt.Println("itinerary usage:\n" +
			"go run . ./input.txt ./output.txt ./airport-lookup.csv")
		return
	}

	inputPath := os.Args[1]
		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			fmt.Println("Input file not found!")
			return
		}

	outputPath := os.Args[2]

	lookupPath := os.Args[3]
		if _, err := os.Stat(lookupPath); os.IsNotExist(err) {
			fmt.Println("Airport lookup file not found!")
			return
		}

	content, err := os.ReadFile(inputPath)
		if err != nil {
				//fmt.Println("Error in Reading Lookup File")
			log.Fatal(err)
		}

	err = os.WriteFile("output.txt", []byte(content), 0644)

		//reg := regexp.MustCompile("#[A-Z]{3}")

		//matches := reg.FindAllString(content, -1)

		//fmt.Println(matches)

		//fmt.Println(string(content))
	*/
}
