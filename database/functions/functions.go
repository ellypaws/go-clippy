package functions

import (
	"github.com/nokusukun/bingo"
)

var ExcelCollection *bingo.Collection[Function]
var SheetsCollection *bingo.Collection[Function]

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
	Layout string          `json:"layout,omitempty"`
	Raw    string          `json:"raw,omitempty"`
	Args   map[string]Args `json:"args,omitempty"`
}

type Args struct {
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"` // string, int, boolean, range, array, function (lambda)
	Variadic    bool   `json:"variadic,omitempty"`
	Optional    bool   `json:"optional,omitempty"`
}

func (f Function) Key() []byte {
	return []byte(f.Name)
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
