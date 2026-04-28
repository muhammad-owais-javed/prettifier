package lib

import (
	"fmt"
	"os"
	"bufio"

)

func StdOutput(colorOutput string, outputPath string) string {

	scan := bufio.NewScanner(os.Stdin)
	fmt.Println("Do you want to print the prettified itinerary to the console? (Enter 'y' for yes)")
	scan.Scan()
	input := scan.Text()

	if input != "y" && input != "Y" {
		//fmt.Println("Output not printed to console.")
		return colorOutput
	}

	colorOutput = TrimSpaces(colorOutput)
	fmt.Println("\n--- Prettified Itinerary ---\n")
	fmt.Println("######\n")
	fmt.Println(colorOutput)
	fmt.Println("\n######")
	fmt.Println("\nSuccessfully wrote output to", outputPath)

	return colorOutput
}