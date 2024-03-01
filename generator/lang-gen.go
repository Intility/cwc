package main

import (
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Language struct {
	Type               string   `yaml:"type"`
	TmScope            string   `yaml:"tm_scope,omitempty"`
	AceMode            string   `yaml:"ace_mode"`
	CodemirrorMode     string   `yaml:"codemirror_mode,omitempty"`
	CodemirrorMimeType string   `yaml:"codemirror_mime_type,omitempty"`
	Color              string   `yaml:"color,omitempty"`
	Aliases            []string `yaml:"aliases,omitempty"`
	Extensions         []string `yaml:"extensions"`
	Filenames          []string `yaml:"filenames,omitempty"`
	Interpreters       []string `yaml:"interpreters,omitempty"`
	Group              string   `yaml:"group,omitempty"`
	LanguageID         int64    `yaml:"language_id"`
}

type Languages map[string]Language

var languageTemplate = `// Code generated with lang-gen. DO NOT EDIT.

package filetree

type Language struct {
	Type               string   ` + "`yaml:\"type\"`" + `
	TmScope            string   ` + "`yaml:\"tm_scope,omitempty\"`" + `
	AceMode            string   ` + "`yaml:\"ace_mode\"`" + `
	CodemirrorMode     string   ` + "`yaml:\"codemirror_mode,omitempty\"`" + `
	CodemirrorMimeType string   ` + "`yaml:\"codemirror_mime_type,omitempty\"`" + `
	Color              string   ` + "`yaml:\"color,omitempty\"`" + `
	Aliases            []string ` + "`yaml:\"aliases,omitempty\"`" + `
	Extensions         []string ` + "`yaml:\"extensions\"`" + `
	Filenames          []string ` + "`yaml:\"filenames,omitempty\"`" + `
	Interpreters       []string ` + "`yaml:\"interpreters,omitempty\"`" + `
	Group              string   ` + "`yaml:\"group,omitempty\"`" + `
	LanguageID         int64    ` + "`yaml:\"language_id\"`" + `
}

var languages = map[string]Language{
	{{- range $key, $value := .}}
	"{{ $key }}": {
		Type: "{{ $value.Type }}",
		TmScope: "{{ $value.TmScope }}",
		AceMode: "{{ $value.AceMode }}",
		CodemirrorMode: "{{ $value.CodemirrorMode }}",
		CodemirrorMimeType: "{{ $value.CodemirrorMimeType }}",
		Color: "{{ $value.Color }}",
		Aliases: []string{ {{ join "\", \"" $value.Aliases "\"" }} },
		Extensions: []string{ {{ join "\", \"" $value.Extensions "\"" }} },
		Filenames: []string{ {{ join "\", \"" $value.Filenames "\"" }} },
		Interpreters: []string{ {{ join "\", \"" $value.Interpreters "\"" }} },
		Group: "{{ $value.Group }}",
		LanguageID: {{ $value.LanguageID }},
	},
	{{- end}}
}
`

// join function will join the string slice with a comma
func join(sep string, s []string, surroundingStr string) string {
	if len(s) == 0 {
		return ""
	}
	if surroundingStr == "" {
		return strings.Join(s, sep)
	}
	return surroundingStr + strings.Join(s, sep) + surroundingStr
}

func main() {

	var languages Languages

	// read the file from web
	client := &http.Client{}
	resp, err := client.Get("https://raw.githubusercontent.com/github/linguist/master/lib/linguist/languages.yml")

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal(data, &languages)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// write the data to a go file using the template
	// so that we can embed the data in the binary

	// create a file
	file, err := os.Create("pkg/filetree/languages.go")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer file.Close()

	// use the data to render the template to disk
	tmpl, err := template.New("languages").Funcs(template.FuncMap{"join": join}).Parse(languageTemplate)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = tmpl.Execute(file, languages)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
