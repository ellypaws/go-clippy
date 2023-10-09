package bingo

import (
	"github.com/PuerkitoBio/goquery"
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

type ExcelScraper struct{}

func (e *ExcelScraper) Scrape() []Function {
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

func (e *ExcelScraper) UpdateUrls(functions []Function) {
	for i, function := range functions {
		if function.Syntax.Raw != "" {
			continue
		}
		e.UpdateSingleUrl(&functions[i])
	}
}

func (e *ExcelScraper) UpdateSingleUrl(function *Function) {
	doc, err := urlToDocument(function.URL)
	if err != nil {
		return
	}

	// Parse the article content
	doc.Find("#supArticleContent").Each(func(i int, s *goquery.Selection) {

		// Function syntax
		e.parseSyntaxSectionExcel(s, function)

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

func (e *ExcelScraper) parseSyntaxSectionExcel(s *goquery.Selection, function *Function) {
	syntaxSection := s.Find("section .ocpSection h2:contains('Syntax')").First().Parent()

	// Transform the Syntax section into the desired Raw format
	var rawBuilder strings.Builder

	var section []string
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

			text = strings.TrimSpace(text)
			section = append(section, text)

		case child.Is("ul"):
			selection := child.Find("li")
			selection.Each(func(j int, item *goquery.Selection) {
				text := strings.TrimSpace(item.Text())
				text = strings.ReplaceAll(text, "Required", "__Required__")
				text = strings.ReplaceAll(text, "Optional", "__Optional__")
				text = "`" + strings.TrimSpace(item.Find("b.ocpRunInHead").Text()) + "`" + strings.ReplaceAll(text, item.Find("b.ocpRunInHead").Text(), "")

				section = append(section, text)
			})
		}
	})

	function.Syntax.Raw = strings.Join(section, "\n")

	function.Syntax.Layout = strings.TrimSpace(syntaxSection.Find("p").First().Text())

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
