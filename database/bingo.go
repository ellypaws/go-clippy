package bingo

import (
	"github.com/nokusukun/bingo"
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

func GetFunction(s string, platform string) Function {
	//result := functions.Query(bingo.Query[Functions]{
	//	Filter: func(doc Functions) bool {
	//		switch platform.(type) {
	//		case excelFunctions:
	//			return strings.Contains(strings.ToLower(doc.Excel.Name), strings.ToLower(s))
	//		case sheetsFunctions:
	//			return strings.Contains(strings.ToLower(doc.Sheets.Name), strings.ToLower(s))
	//		}
	//		return false
	//	},
	//})
	//if !result.Any() {
	//	fmt.Printf("No function found for %s\n", s)
	//}
	var db *bingo.Collection[Function]
	switch platform {
	case "sheets":
		db = sheetsFunctions
	case "excel":
		db = excelFunctions
	}

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

var driver *bingo.Driver

func GetDriver() (*bingo.Driver, error) {
	if driver != nil {
		return driver, nil
	}
	config := bingo.DriverConfiguration{
		DeleteNoVerify: false,
		Filename:       "clippy.bingo",
	}
	return bingo.NewDriver(config)
}

func init() {
	driver, err := GetDriver()
	if err != nil {
		panic(err)
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
