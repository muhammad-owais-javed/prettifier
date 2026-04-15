package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) != 4 || os.Args[1] == "help" {
		fmt.Println("itinerary usage:\n" +
			"go run . ./input.txt ./output.txt ./airport-lookup.csv")
	}

	fmt.Println("Hello World")
}
