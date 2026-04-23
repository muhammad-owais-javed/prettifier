# Prettifier
**Prettifier** is a command line tool written in Go that transforms administrator formatted flight itineraries into a clean, customer friendly prettify format.

## Features
* **Smart Code Replacement:** Converts IATA (#LAX), ICAO (##EGLL), and city (*#LHR) codes into human readable names.
* **Dynamic Data Parsing:** Reads airport and city data from a CSV file and automatically adapts to reordered columns by using the header row.
* **Date & Time Formatting:** Parses ISO 8601 timestamps and converts them into clear, customer friendly date and time formats.
* **Whitespace Normalization:** Standardizes all newline characters and collapses multiple blank lines into a single one for clean, readable output.
* **Error Handling:** Provides clear error messages for common issues (e.g., missing files, malformed data) and prevents writing an output file on failure.
* **Formatted Console Output:** Prints a color coded version of the final itinerary to the terminal for review.

## How It Works

The tool operates as a single execution command line application, processing files through a defined pipeline.

* **Argument Parsing:** The main function first validates the command line arguments (input, output, lookup). It also handles the -h/help flag.
* **Data Loading:** It calls the lib.LoadAirportData function, which reads the provided airport CSV. This function maps column names (e.g., iata_code, municipality) to their index, making the column order flexible. The data is loaded into maps for fast lookups.
* **Text Processing:** The core logic resides in the lib.DateTimeParsing function. It performs a series of string replacements in a specific order:
  * First, it replaces city and airport codes (*#LHR, #LAX, etc.) with their corresponding names, wrapping them in ANSI color codes for the console output.
  * Next, it uses a regular expression to find and parse all date/time tags (D(...), T12(...), T24(...)). Each valid tag is replaced with a formatted string.
* **Whitespace Cleanup:** Finally, lib.TrimSpaces is called to normalize all newline characters and remove excessive blank lines.
* **File Output:** The main function writes the processed plain text to the specified output file and prints the color formatted version to the console.

## Setup and Build

### Prerequisites
* Go (version 1.16 or newer)
* Git (for cloning the repository)

### Cloning the Repository
```sh
git clone https://gitea.kood.tech/muhammadowaisjaved/prettifier.git
cd prettifier
```
### Running the Application

Use go run to compile and execute the program in one step. The command requires paths to an input file, an output file, and an airport lookup CSV.
```sh
go run . ./input.txt ./output.txt ./airport-lookup.csv
```

### Building the Executable

To create a standalone binary, use the go build command.
```sh
go build -o prettifier .
```
then run the compiled application directly:
```sh
./prettifier ./input.txt ./output.txt ./airport-lookup.csv
```

## Usage
To process an itinerary, provide the paths to your files as arguments.

### Example input.txt:
```
Your flight from *#LAX to ##EGLL is confirmed.
Departure: D(2024-09-22T18:05:00-07:00 )
Time: T12(2024-09-22T18:05:00-07:00)
```
###  Command:
```sh
go run . ./input.txt ./output.txt ./airport-lookup.csv
```
### Example output.txt:
```
Your flight from Los Angeles to London Heathrow Airport is confirmed.
Departure: 22-Sep-2024
Time: 06:05PM (-07:00)
```
The tool will also print a color highlighted version of this output directly to the terminal.

## Project Structure

```
.
├── main.go                   # Entry point, CLI handling, and workflow orchestration
├── lib/
│   ├── airport_lookup.go     # Loads and parses the airport data CSV
│   ├── date_time_parsing.go  # Core replacement logic for airports, cities, and dates
│   ├── spaces_parsing.go     # Whitespace normalization logic
│   └── formatting.go         # ANSI color code constants
├── go.mod                    # Go module definition
├── README.md                 # This file
├── airport-lookup-sample.csv # Sampple airport-lookup.csv file
└── input-sample.txt          # Sample input.txt file
