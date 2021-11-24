package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func readCsvFile(filePath string) [][]string {
	// Load a csv file.
	f, _ := os.Open(filePath)

	// Create a new reader.
	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	r.LazyQuotes = true
	r.Comma = ';'

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func main() {
	records := readCsvFile("/Users/ahmet/Desktop/all_data.csv")
	fmt.Println(cap(records), " records")
}
