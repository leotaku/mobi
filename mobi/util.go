package mobi

import (
	"fmt"

	"github.com/leotaku/manki/mobi/pdb"
	r "github.com/leotaku/manki/mobi/records"
)

const preHTML = `<?xml version="1.0" encoding="UTF-8"?>
<html xmlns="http://www.w3.org/1999/xhtml">
  <head>
    <title>Unknown</title>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
  <link rel="stylesheet" type="text/css" href="kindle:flow:0001?mime=text/css"/>
<link rel="stylesheet" type="text/css" href="kindle:flow:0002?mime=text/css"/>
</head>
  <body class="calibre" aid="%04v">
    </body>
</html>`

func chaptersToText(chaps []Chapter) (string, []r.ChunkInfo) {
	text := ""
	info := make([]r.ChunkInfo, 0)
	for i, chap := range chaps {
		pre := fmt.Sprintf(preHTML, i)
		info = append(info, r.ChunkInfo{
			PreStart:      len(text),
			PreLength:     len(pre),
			ContentStart:  len(text) + len(pre),
			ContentLength: len(chap.TextContent),
		})
		text = text + pre + chap.TextContent
	}

	return text, info
}

func genTextRecords(html string) []pdb.Record {
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	records := []pdb.Record{}
	recordCount := len(html) / r.TextRecordMaxSize
	if len(html)%r.TextRecordMaxSize != 0 {
		recordCount += 1
	}

	for i := 0; i < int(recordCount); i++ {
		from := i * r.TextRecordMaxSize
		to := from + r.TextRecordMaxSize

		record := r.NewTBSTextRecord(html[from:min(to, len(html))])
		records = append(records, record)
	}

	len := records[len(records)-1].Length()
	if len%4 != 0 {
		pad := make(pdb.RawRecord, len%4)
		records = append(records, pad)
	}

	return records
}
