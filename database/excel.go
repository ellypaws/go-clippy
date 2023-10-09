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
//	Syntax      Syntax   `json:"args,omitempty"`
//	Example     string   `json:"example,omitempty"`
//	Description string   `json:"description,omitempty"`
//	URL         string   `json:"url,omitempty"`
//	SeeAlso     string   `json:"seealso,omitempty"`
//	Version     []string `json:"version,omitempty"`
//}
//
//type Syntax struct {
//	Layout string `json:"layout,omitempty"`
//	Raw    string `json:"raw,omitempty"`
//	Args   []Args `json:"args,omitempty"`
//}
//
//type Args struct {
//	Description string `json:"description,omitempty"`
//	Type        string `json:"type,omitempty"` // string, int, boolean, range, array, function (lambda)
//	Variadic    bool   `json:"variadic,omitempty"`
//	Optional    bool   `json:"optional,omitempty"`
//}

func RecordExcel() []Function {
	doc, err := urlToDocument("https://support.microsoft.com/en-us/office/excel-functions-alphabetical-b3944572-255d-4efb-bb96-c6d90033e188")
	if err != nil {
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

func urlToDocument(url string) (*goquery.Document, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to request the webpage")
		return nil, err
	}
	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println("Failed to parse the webpage")
		return nil, err
	}
	return doc, nil
}

func UpdateURLs(functions []Function) {
	for i, function := range functions {
		if function.Syntax.Raw != "" {
			continue
		}
		UpdateSingleURL(&functions[i])
	}
}

func UpdateSingleURL(function *Function) {
	doc, err := urlToDocument(function.URL)
	if err != nil {
		return
	}

	// Parse the article content
	doc.Find("#supArticleContent").Each(func(i int, s *goquery.Selection) {

		// Function syntax
		parseSyntaxSection(s, function)

		// Function example
		example := s.Find("section .ocpSection h2:contains('Example')").First().Next()
		function.Example = strings.TrimSpace(example.Text())

		// Function description
		function.Description = strings.TrimSpace(s.Find("section .ocpSection h2:contains('Description')").First().Next().Text())

		// See also links
		s.Find("section .ocpSection h2:contains('See Also')").Each(func(i int, seeAlsoSection *goquery.Selection) {
			seeAlsoLinks := []string{}

			// Start with the next sibling and continue until another h2 or no more siblings
			for node := seeAlsoSection.Next(); node.Size() > 0 && node.Is("p"); node = node.Next() {
				link := node.Find(".ocpArticleLink").First()
				url := link.AttrOr("href", "")
				if !strings.HasPrefix(url, "http") {
					url = "https://support.microsoft.com" + url
				}
				seeAlsoLinks = append(seeAlsoLinks, "["+link.Text()+"]"+"("+url+")")
			}

			function.SeeAlso = strings.Join(seeAlsoLinks, "\n")
		})

		// Function Compatibility
		s.Find("#supAppliesToTableContainer #supAppliesToList .appliesToItem").Each(func(i int, item *goquery.Selection) {
			function.Version = append(function.Version, item.Text())
		})
	})
}

func parseSyntaxSection(s *goquery.Selection, function *Function) {
	syntaxSection := s.Find("section .ocpSection h2:contains('Syntax')").First().Parent()

	// Transform the Syntax section into the desired Raw format
	var rawBuilder strings.Builder

	syntaxSection.Contents().Each(func(i int, child *goquery.Selection) {
		switch {
		case child.Is("p"):
			text := child.Text()
			text = strings.ReplaceAll(text, "Required", "__Required__")
			text = strings.ReplaceAll(text, "Optional", "__Optional__")
			text = strings.TrimSpace(child.Find("b.ocpLegacyBold").Text()) + strings.ReplaceAll(text, child.Find("b.ocpLegacyBold").Text(), "")

			// Add `function()` syntax highlighting on the first line
			if rawBuilder.Len() == 0 {
				text = "`" + text + "`"
			}

			rawBuilder.WriteString(strings.TrimSpace(text))
			rawBuilder.WriteString("\n")

		case child.Is("ul"):
			selection := child.Find("li")
			selection.Each(func(j int, item *goquery.Selection) {
				text := strings.TrimSpace(item.Text())
				text = strings.ReplaceAll(text, "Required", "__Required__")
				text = strings.ReplaceAll(text, "Optional", "__Optional__")
				text = "`" + strings.TrimSpace(item.Find("b.ocpRunInHead").Text()) + "`" + strings.ReplaceAll(text, item.Find("b.ocpRunInHead").Text(), "")
				rawBuilder.WriteString(strings.TrimSpace(text))
				if selection.Size() < j {
					rawBuilder.WriteString("\n")
				}
			})

		}
	})

	// Remove any sequence of more than two newlines
	//raw := rawBuilder.String()
	//reg := regexp.MustCompile(`\n{3,}`)
	//function.Syntax.Raw = reg.ReplaceAllString(raw, "\n\n")
	function.Syntax.Raw = rawBuilder.String()

	function.Syntax.Layout = strings.TrimSpace(syntaxSection.Find("p b.ocpLegacyBold").First().Text())

	function.Syntax.Args = []Args{}
	syntaxSection.Find("ul > li").Each(func(i int, argItem *goquery.Selection) {
		var arg Args

		// Grab the name and primary description
		parts := strings.SplitN(strings.TrimSpace(argItem.Find(".ocpRunInHead").First().Text()), "\u00A0", 2)
		if len(parts) > 1 {
			arg.Description = parts[1]
		}

		// If there's an internal ul, append its text to the description
		argItem.Find("ul li").Each(func(j int, subItem *goquery.Selection) {
			arg.Description += "; " + strings.TrimSpace(subItem.Text())
		})

		function.Syntax.Args = append(function.Syntax.Args, arg)
	})
}
