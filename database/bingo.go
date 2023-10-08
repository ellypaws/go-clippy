package bingo

import (
	"github.com/nokusukun/bingo"
	"log"
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
	Args        []string `json:"args,omitempty"`
	Example     string   `json:"example,omitempty"`
	Description string   `json:"description,omitempty"`
	URL         string   `json:"url,omitempty"`
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
	ir := db.Insert(f)
	if ir.Error() != nil {
		log.Fatal(ir.Error())
	}
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
