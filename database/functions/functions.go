package functions

import (
	"github.com/nokusukun/bingo"
	"strings"
)

var ExcelCollection *bingo.Collection[Function]
var SheetsCollection *bingo.Collection[Function]

type Function struct {
	Name        string   `json:"name,omitempty"`
	Category    string   `json:"category,omitempty"`
	Syntax      Syntax   `json:"syntax,omitempty"`
	Example     string   `json:"example,omitempty"`
	Description string   `json:"description,omitempty"`
	URL         string   `json:"url,omitempty"`
	SeeAlso     string   `json:"seealso,omitempty"`
	Version     []string `json:"version,omitempty"`
}

type Syntax struct {
	Layout string `json:"layout"`
	Raw    string `json:"raw"`
	Args   []Args `json:"args,omitempty"`
}

type Args struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"` // string, int, boolean, range, array, function (lambda)
	Variadic    bool   `json:"variadic"`
	Optional    bool   `json:"optional"`
}

func (f Function) Key() []byte {
	return []byte(strings.ToUpper(f.Name))
}

func GetCollection(platform string) *bingo.Collection[Function] {
	switch platform {
	case "sheets":
		return SheetsCollection
	case "excel":
		return ExcelCollection
	default:
		return ExcelCollection
	}
}
