package main

import (
	"github.com/PuerkitoBio/goquery"
	"go-clippy/database"
	"go-clippy/database/functions"
	"strings"
)

//type Function struct {
//	Name        string   `json:"name"`
//	Category    string   `json:"category,omitempty"`
//	Syntax      Syntax   `json:"syntax"`
//	Example     string   `json:"example,omitempty"`
//	Description string   `json:"description,omitempty"`
//	URL         string   `json:"url,omitempty"`
//	SeeAlso     string   `json:"see_also,omitempty"`
//	Version     []string `json:"version,omitempty"`
//}
//
//type Syntax struct {
//	Layout string `json:"layout"`
//	Raw    string `json:"raw"`
//	Args   []Args `json:"args,omitempty"`
//}
//
//type Args struct {
//	Name        string `json:"name"`
//	Description string `json:"description"`
//	Type        string `json:"type"` // string, int, boolean, range, array, function (lambda)
//	Variadic    bool   `json:"variadic"`
//	Optional    bool   `json:"optional"`
//}

type SheetsUrl string

func (url *SheetsUrl) Scrape() []functions.Function {
	doc, err := database.UrlToDocument(string(*url))
	if err != nil {
		return nil
	}

	var f []functions.Function

	// Find the table inside the specified section
	targetTable := doc.Find("#hcfe-content > section > div > div.main-content > article > section > div.dyn-table > table")

	// Iterate through table rows
	targetTable.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		// Extract information from each row
		funcName := row.Find("td:nth-child(2)").Text()
		category := row.AttrOr("data-category", "")
		syntaxLayout := row.Find("td:nth-child(3) code").Text()
		descriptionElement := row.Find("td:nth-child(4)")
		description := descriptionElement.Text()
		// Extract the full Google URL and format as [Learn more](url)
		url := descriptionElement.Find("a").AttrOr("href", "")
		if !strings.HasPrefix(url, "http") {
			url = "https://support.google.com" + url
		}
		description = strings.ReplaceAll(description, "[Learn more]", "[Learn more]("+url+")")

		// Create a custom function struct and add it to the list
		function := functions.Function{
			Name:     funcName,
			Category: category,
			Syntax: functions.Syntax{
				Layout: syntaxLayout,
			},
			Description: description,
			URL:         url,
		}
		f = append(f, function)
	})

	return f
}

func (url *SheetsUrl) UpdateUrls(functions []functions.Function) {
	for i, function := range functions {
		if function.Syntax.Raw != "" {
			continue
		}
		url.UpdateSingleUrl(&functions[i])
	}
}

func (url *SheetsUrl) UpdateSingleUrl(function *functions.Function) {
	doc, err := database.UrlToDocument(function.URL)
	if err != nil {
		return
	}

	baseUrl := "https://support.google.com"

	// Find the target section with the provided selector
	targetSection := doc.Find("#hcfe-content > section > div > div.main-content > article > section")

	// Extract the function name from the h1 tag
	function.Name = targetSection.Find("h1").Text()

	// Extract and format sample usage as a single string with newline \n
	sampleUsage := targetSection.Find("h3:contains('Sample Usage') + p code").Map(func(i int, s *goquery.Selection) string {
		return s.Text()
	})
	function.Example = strings.Join(sampleUsage, "\n")

	// Extract the syntax layout
	syntaxLayout := targetSection.Find("h3:contains('Syntax') + p code").Text()
	function.Syntax.Layout = syntaxLayout

	// Extract and format the raw content for Discord
	rawContent := targetSection.Find("h3:contains('Syntax') + ul li").Map(func(i int, s *goquery.Selection) string {
		// Extract the text
		text := s.Text()

		// Format code tags as `code`
		text = strings.ReplaceAll(text, "<code>", "`")
		text = strings.ReplaceAll(text, "</code>", "`")

		// Format strong and bold tags as **text**
		text = strings.ReplaceAll(text, "<strong>", "**")
		text = strings.ReplaceAll(text, "</strong>", "**")
		text = strings.ReplaceAll(text, "<b>", "**")
		text = strings.ReplaceAll(text, "</b>", "**")

		// Format links as [text](link)
		s.Find("a").Each(func(i int, link *goquery.Selection) {
			linkText := link.Text()
			linkHref, _ := link.Attr("href")
			if !strings.HasPrefix(linkHref, "http") {
				linkHref = baseUrl + linkHref
			}
			linkReplacement := "[" + linkText + "](" + linkHref + ")"
			text = strings.ReplaceAll(text, link.Text(), linkReplacement)
		})

		return text
	})

	function.Syntax.Raw = strings.Join(rawContent, "\n")

	// Extract notes and format them
	notes := targetSection.Find("h3:contains('Notes') + ul li").Map(func(i int, s *goquery.Selection) string {
		return s.Text()
	})
	function.Description += "\n\nNotes:\n" + strings.Join(notes, "\n")

	// Extract and format "See Also" links
	seeAlsoLinks := targetSection.Find("h3:contains('See Also') + p a").Map(func(i int, s *goquery.Selection) string {
		linkText := s.Text()
		linkHref, _ := s.Attr("href")
		if !strings.HasPrefix(linkHref, "http") {
			linkHref = baseUrl + linkHref
		}
		return "[" + linkText + "](" + linkHref + ")"
	})
	function.SeeAlso = strings.Join(seeAlsoLinks, "\n")

	// TODO: Parse and format the argument details if needed
	// You can iterate through elements under the syntax section and extract argument details.

	// Example of checking for optional and variadic
	if strings.Contains(function.Syntax.Raw, "[") || strings.Contains(function.Syntax.Raw, "]") {
		function.Syntax.Args[0].Optional = true
	}
	if strings.Contains(function.Syntax.Raw, "...") {
		function.Syntax.Args[1].Variadic = true
	}
}
