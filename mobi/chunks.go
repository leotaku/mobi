package mobi

var defaultTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<html xmlns="http://www.w3.org/1999/xhtml">
  <head>
    <title>{{ .Mobi.Title }}</title>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
    {{ range $i, $_ := .Mobi.CssFlows }}
    <link rel="stylesheet" type="text/css" href="kindle:flow:{{ $i | inc | printf "%03v" }}?mime=text/css"/>
    {{ end }}
  </head>
  <body class="calibre" aid="{{ .Id | printf "%04v" }}">
  </body>
</html>`

func SingleChunk(ss ...string) []Chunk {
	return SingleChunkTemplated(defaultTemplate, ss...)
}

func SingleChunkTemplated(tpl string, ss ...string) []Chunk {
	result := make([]Chunk, 0)
	for _, s := range ss {
		result = append(result, Chunk{
			Head: tpl,
			Body: s,
		})
	}
	return result
}

func Chunks(s string) []Chunk {
	return ChunksTemplated(defaultTemplate, s)
}

func ChunksTemplated(tpl string, s string) []Chunk {
	panic("unimplemented")
}
