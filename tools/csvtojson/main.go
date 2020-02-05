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
	if len(os.Args) < 3 {
		printHelp(nil)
		return
	}

	sourceFile := os.Args[1]
	resultFile := os.Args[2]

	source, err := os.Open(sourceFile)
	if err != nil {
		log.Println("Failed to open source file")
		log.Fatalln(err)
	}
	defer source.Close()

	result, err := os.Create(resultFile)
	if err != nil {
		log.Println("Failed to create result file")
		log.Fatalln(err)
	}
	defer result.Close()

	csvReader := csv.NewReader(source)

	keys, err := csvReader.Read()
	if err != nil {
		log.Println("Failed to read keys")
		log.Fatalln(err)
	}

	types, err := csvReader.Read()
	if err != nil {
		log.Println("Failed to read types")
		log.Fatalln(err)
	}

	fmt.Fprint(result, "[")
	for rowHasComma := false; true; rowHasComma = true {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println("Failed to read row")
			log.Fatalln(err)
		}
		if rowHasComma {
			fmt.Fprint(result, ",")
		}
		fmt.Fprint(result, "{")

		fieldHasComma := false
		for i, name := range keys {
			if fieldHasComma {
				fmt.Fprint(result, ",")
			} else {
				fieldHasComma = true
			}
			fieldType := types[i]
			value := row[i]
			switch fieldType {
			case "string":
				fmt.Fprintf(result, "%q: %q", name, value)

			case "int":
				fmt.Fprintf(result, "%q: %s", name, value)

			case "bool":
				var boolValue bool
				switch strings.ToLower(value) {
				case "true":
				case "t":
				case "1":
					boolValue = true

				case "false":
				case "f":
				case "0":
					boolValue = false
				}
				fmt.Fprintf(result, "%q: %t", name, boolValue)
			}
		}

		fmt.Fprint(result, "}")
	}
	fmt.Fprint(result, "]")

}

func printHelp(args []string) {
	fmt.Println("Will's CSV to JSON Tool")
	fmt.Println("Usage: csvtojson <in.csv> <out.json>")
	fmt.Println("First row will be used as JSON keys")
	fmt.Println("Second row will be used as data types")
}
