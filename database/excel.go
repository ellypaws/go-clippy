package bingo

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//type Function struct {
//	Name        string   `json:"name,omitempty"`
//	Category    string   `json:"category,omitempty"`
//	Args        []string `json:"args,omitempty"`
//	Example     string   `json:"example,omitempty"`
//	Description string   `json:"description,omitempty"`
//	URL         string   `json:"url,omitempty"`
//}

func RecordExcel() []Function {
	// Make HTTP request
	response, err := http.Get("https://support.microsoft.com/en-us/office/excel-functions-alphabetical-b3944572-255d-4efb-bb96-c6d90033e188")
	if err != nil {
		fmt.Println("Failed to request the webpage")
		return nil
	}
	defer response.Body.Close()

	// Create a goquery document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println("Failed to parse the webpage")
		return nil
	}

	var functions []Function

	doc.Find("table.ms-rteTable-0").Each(func(i int, s *goquery.Selection) {
		s.Find("tbody tr").Each(func(j int, contentSelection *goquery.Selection) {
			function := Function{}

			contentSelection.Find("td").Each(func(k int, td *goquery.Selection) {
				if k == 0 {
					function.Name = strings.Split(td.Text(), " ")[0]
					function.URL, _ = td.Find("a").Attr("href")
				} else if k == 1 {
					function.Category = strings.Split(td.Text(), ":")[0]
				}
			})

			functions = append(functions, function)
		})
	})

	return functions
}
