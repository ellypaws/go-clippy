package bingo

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/nokusukun/bingo"
	"log"
	"net/http"
	"strings"
)

var clippy *bingo.Collection[Clippy]
var excelFunctions *bingo.Collection[Function]
var sheetsFunctions *bingo.Collection[Function]

type Clippy struct {
	Username string `json:"username,omitempty" validate:"required,email"`
	Points   int    `json:"points,omitempty" validate:"required,min=3,max=64"`
}

type Function struct {
	Name        string   `json:"name,omitempty"`
	Category    string   `json:"category,omitempty"`
	Syntax      Syntax   `json:"args,omitempty"`
	Example     string   `json:"example,omitempty"`
	Description string   `json:"description,omitempty"`
	URL         string   `json:"url,omitempty"`
	SeeAlso     string   `json:"seealso,omitempty"`
	Version     []string `json:"version,omitempty"`
}

type Syntax struct {
	Layout string `json:"layout,omitempty"`
	Raw    string `json:"raw,omitempty"`
	Args   []Args `json:"args,omitempty"`
}

type Args struct {
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"` // string, int, boolean, range, array, function (lambda)
	Variadic    bool   `json:"variadic,omitempty"`
	Optional    bool   `json:"optional,omitempty"`
}

func getCollection(platform string) *bingo.Collection[Function] {
	switch platform {
	case "sheets":
		return sheetsFunctions
	case "excel":
		return excelFunctions
	}
	return nil
}

func GetFunction(s string, platform string) Function {
	db := getCollection(platform)

	result := db.Query(bingo.Query[Function]{
		Filter: func(doc Function) bool {
			return strings.Contains(strings.ToLower(doc.Name), strings.ToLower(s))
		},
	})
	// -------------------------------------- or ------------------------------------------
	_, _ = db.FindOne(func(doc Function) bool {
		return strings.Contains(strings.ToLower(doc.Name), strings.ToLower(s)) //
	})

	if !result.Any() {
		return Function{}
	}

	return *result.First()
}

func Record(f Function, platform string) {
	db := getCollection(platform)
	id, err := db.Insert(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted", id)
}

var driver *bingo.Driver

func GetDriver() (*bingo.Driver, error) {
	if driver != nil {
		return driver, nil
	}
	config := bingo.DriverConfiguration{
		DeleteNoVerify: false,
		Filename:       "clippy.bingo.db",
	}
	return bingo.NewDriver(config)
}

func init() {
	driver, err := GetDriver()
	if err != nil {
		log.Fatal(err)
	}
	clippy = ClippyCollection(driver)

	sheetsFunctions = bingo.CollectionFrom[Function](driver, "sheets")
	excelFunctions = bingo.CollectionFrom[Function](driver, "excel")
}

func ClippyCollection(d *bingo.Driver) *bingo.Collection[Clippy] {
	return bingo.CollectionFrom[Clippy](d, "clippy")
}

//func FunctionCollection(d *bingo.Driver) *bingo.Collection[Functions] {
//	return bingo.CollectionFrom[Functions](d, "functions")
//}

func (c Clippy) Key() []byte {
	return []byte(c.Username)
}

//func (f Functions) Key() []byte {
//	return []byte("functions")
//}

func (f Function) Key() []byte {
	return []byte(f.Name)
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
