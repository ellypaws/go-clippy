package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	discord "go-clippy/bot"
	"go-clippy/database"
	"go-clippy/database/functions"
	"log"
	"os"
	"slices"
)

var sliceFuncs []functions.Function

func main() {
	choose()
}

func choose() {
	fmt.Println("Choose an option:")
	fmt.Println("1. Run Discord bot")
	fmt.Println("2. Record Excel functions to database")
	fmt.Println("2.1. Update Excel function URLs")
	fmt.Println("3. Record Sheets functions to database")
	fmt.Println("4. Test random url")
	fmt.Println("5. Write to database")
	fmt.Println("6. Test syntax")
	fmt.Print("Option: ")
	var choice string
	excelUrl := UrlToScrape("https://support.microsoft.com/en-us/office/excel-functions-alphabetical-b3944572-255d-4efb-bb96-c6d90033e188")
	fmt.Scanln(&choice)
	switch choice {
	case "1":
		fmt.Println("Running Discord bot...")
		discord.Run()
	case "2":
		fmt.Println("Recording Excel functions...")
		sliceFuncs = excelUrl.Scrape()
		for _, function := range sliceFuncs[:min(100, len(sliceFuncs))] {
			log.Printf("%+v\n", function.Name)
		}
	case "2.1":
		fmt.Println("Updating Excel function URLs...")
		if sliceFuncs == nil {
			sliceFuncs = excelUrl.Scrape()
		}
		excelUrl.UpdateUrls(sliceFuncs[:5])
		fmt.Printf("%+v\n", sliceFuncs[:1])

		// save to json
		indent, _ := json.MarshalIndent(sliceFuncs[:5], "", "    ")
		os.WriteFile("excel.json", indent, 0644)
	case "3":
		fmt.Println("Recording Sheets functions...")
	case "4":
		fmt.Println("Testing random url...")
		testFunction := functions.Function{
			Name:        "",
			Category:    "",
			Syntax:      functions.Syntax{},
			Example:     "",
			Description: "",
			URL:         "https://support.microsoft.com/en-us/office/sum-function-043e1c7d-7726-4e80-8f32-07b23e057f89",
			SeeAlso:     "",
			Version:     nil,
		}
		excelUrl.UpdateSingleUrl(&testFunction)
		fmt.Printf("%+v\n", testFunction)

		// save to json
		indent, _ := json.MarshalIndent(testFunction, "", "    ")
		os.WriteFile("sum.json", indent, 0644)
	case "5":
		fmt.Println("Writing to database...")

		// reverse slicefuncs
		slices.Reverse(sliceFuncs)

		// sort the slice by alphabetical order in a.Name
		slices.SortStableFunc(sliceFuncs, func(a, b functions.Function) int {
			if n := cmp.Compare(a.Name, b.Name); n != 0 {
				return n
			}
			return cmp.Compare(a.Syntax.Layout, a.Syntax.Layout)
		})

		log.Printf("Is this sorted?")
		for _, function := range sliceFuncs[:min(25, len(sliceFuncs))] {
			log.Printf("%+v\n", function.Name)
		}

		for _, function := range sliceFuncs {
			function.Record(functions.GetCollection("excel"))
		}
	case "close":
		fmt.Println("Closing driver...")
		database.Close()
	default:
		SyntaxTest()
	}
	fmt.Println("Done!")
	choose()
}
