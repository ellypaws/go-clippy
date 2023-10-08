package bingo

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	"strings"
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
	//response, err := http.Get("https://support.microsoft.com/en-us/office/excel-functions-alphabetical-b3944572-255d-4efb-bb96-c6d90033e188")
	//if err != nil {
	//	fmt.Println("Failed to request the webpage")
	//	return nil
	//}
	//defer response.Body.Close()
	//
	//doc, err := goquery.NewDocumentFromReader(response.Body)
	file, _ := os.Open("database/excel.html")
	defer file.Close()
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		fmt.Println("Failed to parse the webpage")
		return nil
	}

	var functions []Function

	doc.Find("table.banded.flipColors tbody tr").Each(func(i int, tr *goquery.Selection) {
		funcName := tr.Find("td:first-child p a")
		categoryDesc := strings.Split(tr.Find("td:nth-child(2) p").Text(), ":")
		url := funcName.AttrOr("href", "")
		// append the base url if the url is relative
		if !strings.HasPrefix(url, "http") {
			url = "https://support.microsoft.com" + url
		}

		function := Function{
			Name:        strings.TrimSpace(strings.Split(funcName.Text(), " ")[0]),
			URL:         url,
			Category:    strings.TrimSpace(categoryDesc[0]),
			Description: strings.TrimSpace(categoryDesc[1]),
		}

		functions = append(functions, function)
	})

	return functions
}
