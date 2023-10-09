package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

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

func testSection(i int) []string {
	switch i {
	case 1:
		return []string{
			"ACCRINTM(issue, settlement, rate, par, [basis])",
			"The ACCRINTM function syntax has the following arguments:",
			"`Issue`    __Required__. The security's issue date.",
			"`Settlement`    __Required__. The security's maturity date.",
		}
	case 2:
		return []string{
			"SUM(number1, [number2], ...)",
			"The SUM function syntax has the following arguments:",
			"`Number1`    __Required__. The first number, cell reference, or range for which you want the total.",
			"`[number2-255]`    __Optional__. Additional numbers, cell references or ranges for which you want the total, up to a maximum of 255.",
		}
	}
	return nil
}

func main() {
	function := Function{
		Name:     "test",
		Category: "test",
		Syntax: Syntax{
			Layout: "function(args)",
			Raw:    "",
			Args:   nil,
		},
		Example:     "",
		Description: "this is a test function",
		URL:         "https://support.microsoft.com/en-us/office/sum-function-043e1c7d-7726-4e80-8f32-07b23e057f89",
		SeeAlso:     "[test2](https://support.microsoft.com/en-us/office/sum-function-043e1c7d-7726-4e80-8f32-07b23e057f89)",
		Version:     []string{"test"},
	}
	testParse(&function)
	json, _ := json.MarshalIndent(function, "", "    ")
	os.WriteFile("test.json", json, 0644)
	fmt.Println(string(json))
}

func testParse(function *Function) {
	function.Syntax.Args = map[string]Args{}
	// TODO: Parse arguments
	// Read from the raw syntax and parse the arguments
	// We can read from section[2:] to ignore the first 2 lines
	// Example:

	//`ACCRINTM(issue, settlement, rate, par, [basis])`
	//`The ACCRINTM function syntax has the following arguments:`
	//`Issue`    __Required__. The security's issue date.
	//`Settlement`    __Required__. The security's maturity date.

	// Start from the 3rd line.
	// Each line is a new argument in the map
	// Issue becomes the key, use the backticks for this
	// Required/Optional fills in the Args.Optional bool.
	// We can also check if the argument has square brackets. CONCAT(text1, [text2],…)
	// Take care of random spaces in each line
	// The rest of the line becomes the Args.Description
	// If there is an ellipsis, then the Args.Variadic bool is true. Check for ... or …

	section := testSection(2)

	for _, line := range section[2:] {
		// Ignore empty lines
		if line == "" {
			continue
		}

		line = strings.TrimSpace(line)

		// Split the line into words
		words := strings.Fields(line)

		// The first word is the argument name
		argName := strings.ToLower(words[0])
		argName = strings.ReplaceAll(argName, "`", "")

		description := strings.Join(words[2:], " ")         // The rest of the words are the description
		description = strings.TrimPrefix(description, ". ") // Remove the ". " from the description if it exists

		// Check if the argument is optional
		optional := strings.Contains(line, "Optional")

		// Check if the argument is variadic
		variadic :=
			strings.Contains(line, "...") ||
				strings.Contains(line, "…") ||
				strings.Contains(line, "-")

		// Remove the backticks from the argument name
		argName = strings.Trim(argName, "`")

		// Infer type from the argName
		infer := map[string][]string{
			"number":  {"number", "num", "digit", "integer", "int", "float", "double", "decimal"},
			"text":    {"text", "string", "str", "char", "character"},
			"range":   {"range", "rng"},
			"array":   {"array", "arr", "list", "collection", "set", "map"},
			"boolean": {"criteria", "condition", "logical", "boolean", "bool", "true", "false"},
		}
		var argType string
		for t, types := range infer {
			for _, accepted := range types {
				if strings.Contains(strings.ToLower(argName), accepted) {
					argType = t
					break
				}
			}
		}

		// Add the argument to the Syntax.Args map
		function.Syntax.Args[argName] = Args{
			Description: description,
			Type:        argType,
			Variadic:    variadic,
			Optional:    optional,
		}
	}
}
