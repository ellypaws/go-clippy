package bingo

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

//type Function struct {
//	Name        string   `json:"name,omitempty"`
//	Category    string   `json:"category,omitempty"`
//	Args        []string `json:"args,omitempty"`
//	Example     string   `json:"example,omitempty"`
//	Description string   `json:"description,omitempty"`
//	URL         string   `json:"url,omitempty"`
//	SeeAlso     string   `json:"seealso,omitempty"`
//}

func RecordExcel() []Function {
	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://support.microsoft.com/en-us/office/excel-functions-alphabetical-b3944572-255d-4efb-bb96-c6d90033e188", nil)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to request the webpage")
		return nil
	}
	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println("Failed to parse the webpage")
		return nil
	}

	var functions []Function
	doc.Find("#supArticleContent > article > section.ocpIntroduction > table").Each(func(i int, table *goquery.Selection) {
		table.Find("tbody tr").Each(func(j int, tr *goquery.Selection) {
			funcName := tr.Find("td:first-child p a")
			categoryDesc := strings.Split(tr.Find("td:nth-child(2) p").Text(), ":")
			url := funcName.AttrOr("href", "")
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
	})

	return functions
}
