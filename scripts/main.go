package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func getPatternFromFile(filename string) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Warning: Could not read %s. Returning empty string.\n", filename)
		return ""
	}

	items := strings.Split(string(content), "|") // already split by |
	var validItems []string

	for _, item := range items {
		cleaned := strings.TrimSpace(item)
		if len(cleaned) > 0 {
			validItems = append(validItems, cleaned)
		}
	}

	// Sort by length descending to prevent partial matches
	sort.Slice(validItems, func(i, j int) bool {
		return len(validItems[i]) > len(validItems[j])
	})

	fmt.Printf("Total %s: %d patterns\n", strings.Split(filename, ".")[0], len(validItems))

	return strings.Join(validItems, "|")
}

func main() {
	// Generate the regex strings
	functionsStr := getPatternFromFile("functions.txt")
	variablesStr := getPatternFromFile("variables.txt")
	constantsStr := getPatternFromFile("constants.txt")

	// Build the JSON structure (copy of old tmLanguage.json)
	tmLanguage := map[string]any{
		"$schema":   "https://raw.githubusercontent.com/martinring/tmlanguage/master/tmlanguage.json",
		"name":      "Maxima",
		"scopeName": "source.mac",
		"patterns": []map[string]string{
			{"include": "#comments"},
			{"include": "#built_in_functions"},
			{"include": "#built_in_variables"},
			{"include": "#built_in_constants"},
			{"include": "#keywords"},
			{"include": "#operators"},
			{"include": "#double_quoted_strings"},
			{"include": "#numbers"},
			{"include": "#user_created_functions"},
			{"include": "#user_created_variables"},
		},
		"repository": map[string]any{
			"comments": map[string]string{
				"name":  "comment.block.mac",
				"begin": `/\*`,
				"end":   `\*/`,
			},
			"keywords": map[string]any{
				"name":  "keyword.control.mac",
				"match": `\b(if|then|else|elseif|while|do|for|to|return|block|quote|load|unless|from|From|thru|Thru|step|Step|next|Next|in|go|catch|throw)\b`,
			},
			"built_in_functions": map[string]string{
				"name":  "support.function.mac",
				"match": fmt.Sprintf(`\b(%s)\b(?=\s*\()`, functionsStr),
			},
			"user_created_functions": map[string]any{
				"name":  "entity.name.function.mac",
				"match": `\b[a-zA-Z_][a-zA-Z0-9_']*\b(?=\s*\()`,
				"captures": map[string]any{
					"0": map[string]string{"name": "entity.name.function.mac"},
				},
			},
			"built_in_variables": map[string]string{
				"name":  "support.variable.mac",
				"match": fmt.Sprintf(`\b(%s)\b`, variablesStr),
			},
			"built_in_constants": map[string]string{
				"name":  "support.constant.mac",
				"match": fmt.Sprintf(`(?<![a-zA-Z0-9_])(%s)\b`, constantsStr),
			},
			"user_created_variables": map[string]string{
				"name":  "variable.other.mac",
				"match": `\b[a-zA-Z_][a-zA-Z0-9_']*\b`,
			},
			"operators": map[string]string{
				"name":  "keyword.operator.mac",
				"match": `(?:\+|-|\*|/|\^|%|:|:=|=|#|<|>|<=|>=|/=)|\b(and|or|not)\b`,
			},
			"double_quoted_strings": map[string]any{
				"name":  "string.quoted.double.mac",
				"begin": `"`,
				"end":   `"`,
				"patterns": []map[string]string{
					{"name": "constant.character.escape.mac", "match": `\\.`},
				},
			},
			"numbers": map[string]string{
				"name":  "constant.numeric.mac",
				"match": `\b\d+(\.\d+)?([eE][+-]?\d+)?\b`,
			},
		},
	}

	jsonData, err := json.MarshalIndent(tmLanguage, "", "    ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}

	outDir := filepath.Join("..", "syntaxes")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Println("Error creating syntaxes directory:", err)
		return
	}

	outPath := filepath.Join(outDir, "mac.tmLanguage.json")
	if err := os.WriteFile(outPath, jsonData, 0644); err != nil {
		fmt.Println("Error writing output file:", err)
		return
	}

	fmt.Printf("Successfully generated %s!\n", outPath)
}
