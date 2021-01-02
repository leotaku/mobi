package mobi

import (
	"bytes"
	"text/template"
)

const defaultTemplateString = `<?xml version="1.0" encoding="UTF-8"?>
<html xmlns="http://www.w3.org/1999/xhtml">
  <head>
    <title>{{ .Mobi.Title }}</title>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
    {{ range $i, $_ := .Mobi.CSSFlows }}
    <link rel="stylesheet" type="text/css" href="kindle:flow:{{ $i | inc | printf "%03v" }}?mime=text/css"/>
    {{ end }}
  </head>
  <body aid="{{ .Chunk.ID | printf "%04v" }}">
  </body>
</html>`

var funcMap = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
	},
}

var defaultTemplate = template.Must(template.New("default").Funcs(funcMap).Parse(defaultTemplateString))

type inventory struct {
	Mobi    Book
	Chapter struct {
		Title string
		ID    int
	}
	Chunk struct {
		ID int
	}
}

func newInventory(m Book, c Chapter, chapID int, chunkID int) inventory {
	return inventory{
		Mobi: m,
		Chapter: struct {
			Title string
			ID    int
		}{
			Title: c.Title,
			ID:    chapID,
		},
		Chunk: struct {
			ID int
		}{
			ID: chunkID,
		},
	}
}

func runTemplate(tpl template.Template, v interface{}) (string, error) {
	buf := bytes.NewBuffer(nil)
	err := tpl.Execute(buf, v)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
