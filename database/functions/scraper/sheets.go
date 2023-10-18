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
		description = strings.ReplaceAll(description, "Learn more", "[Learn more]("+url+")")

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
	//function.Name = targetSection.Find("h1").Text()

	// Extract and format sample usage as a single string with newline \n
	sampleUsage := targetSection.Find("h3:contains('Sample') + p code").Map(func(i int, s *goquery.Selection) string {
		return s.Text()
	})
	function.Example = strings.Join(sampleUsage, "\n")

	// Extract the syntax layout
	syntaxLayout := targetSection.Find("h3:contains('Syntax') + p code").Text()
	function.Syntax.Layout = syntaxLayout

	// Extract and format the raw content for Discord
	var rawContents []string
	targetSection.Find("h3:contains('Syntax')").NextUntil("h3").Each(func(i int, s *goquery.Selection) {
		if s.Is("p") {
			// Handle <code>, <strong> and <a> tags specially
			content := ""
			s.Contents().Each(func(i int, contentNode *goquery.Selection) {
				switch goquery.NodeName(contentNode) {
				case "code":
					content += "`" + contentNode.Text() + "`"
				case "a":
					linkHref, _ := contentNode.Attr("href")
					if !strings.HasPrefix(linkHref, "http") {
						linkHref = baseUrl + linkHref
					}
					content += "[" + contentNode.Text() + "](" + linkHref + ")"
				case "strong", "bold":
					content += "**" + contentNode.Text() + "**"
				default:
					content += contentNode.Text()
				}
			})
			rawContents = append(rawContents, content)
		} else if s.Is("ul") {
			s.Find("li p").Each(func(j int, listItem *goquery.Selection) {
				// Handle <code>, <strong> and <a> tags specially within list items
				content := ""
				listItem.Contents().Each(func(k int, contentNode *goquery.Selection) {
					switch goquery.NodeName(contentNode) {
					case "code":
						content += "`" + contentNode.Text() + "`"
					case "a":
						linkHref, _ := contentNode.Attr("href")
						if !strings.HasPrefix(linkHref, "http") {
							linkHref = baseUrl + linkHref
						}
						content += "[" + contentNode.Text() + "](" + linkHref + ")"
					case "strong", "bold":
						content += "**" + contentNode.Text() + "**"
					default:
						content += contentNode.Text()
					}
				})
				rawContents = append(rawContents, content)
			})
		}
	})

	function.Syntax.Raw = strings.Join(rawContents, "\n")

	// Extract and format "See Also" links
	seeAlsoLinks := targetSection.Find("h3:contains('See Also') ~ p").Map(func(i int, s *goquery.Selection) string {
		var description, linkText, linkHref string
		s.Contents().Each(func(_ int, content *goquery.Selection) {
			if content.Is("a") {
				linkText = content.Text()
				linkHref, _ = content.Attr("href")
				if !strings.HasPrefix(linkHref, "http") {
					linkHref = baseUrl + linkHref
				}
			} else {
				description += content.Text() + " "
			}
		})
		return "[" + linkText + "](" + linkHref + "): " + strings.TrimSpace(description)
	})
	function.SeeAlso = strings.Join(seeAlsoLinks, "\n")

	// TODO: Parse and format the argument details if needed
	// You can iterate through elements under the syntax section and extract argument details.

}
