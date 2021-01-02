package mobi

import (
	"bytes"
	"html/template"
)

const defaultTemplateString = `<?xml version="1.0" encoding="UTF-8"?>
<html xmlns="http://www.w3.org/1999/xhtml">
  <head>
    <title>{{ .Mobi.Title }}</title>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
    {{ range $i, $_ := .Mobi.CssFlows }}
    <link rel="stylesheet" type="text/css" href="kindle:flow:{{ $i | inc | printf "%03v" }}?mime=text/css"/>
    {{ end }}
  </head>
  <body aid="{{ .Chunk.Id | printf "%04v" }}">
  </body>
</html>`

var funcMap = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
	},
}

var defaultTemplate = template.Must(template.New("default").Funcs(funcMap).Parse(defaultTemplateString))

type inventory struct {
	Mobi    MobiBook
	Chapter struct {
		Title string
		Id    int
	}
	Chunk struct {
		Id int
	}
}

func newInventory(m MobiBook, c Chapter, chapId int, chunkId int) inventory {
	return inventory{
		Mobi: m,
		Chapter: struct {
			Title string
			Id    int
		}{
			Title: c.Title,
			Id:    chapId,
		},
		Chunk: struct {
			Id int
		}{
			Id: chunkId,
		},
	}
}

func runTemplate(tpl template.Template, v interface{}) (string, error) {
	buf := bytes.NewBuffer(nil)
	err := tpl.Execute(buf, v)
	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}
