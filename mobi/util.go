package mobi

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/leotaku/manki/mobi/pdb"
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

type ChunkInfo struct {
	PreStart      int
	PreLength     int
	ContentStart  int
	ContentLength int
}

func chaptersToText(c []Chapter) (string, []ChunkInfo) {
	text := ""
	info := make([]ChunkInfo, 0)
	for i, chap := range c {
		pre := fmt.Sprintf(preHTML, i)
		info = append(info, ChunkInfo{
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
	recordCount := len(html) / TextRecordMaxSize
	if len(html)%TextRecordMaxSize != 0 {
		recordCount += 1
	}

	for i := 0; i < int(recordCount); i++ {
		from := i * TextRecordMaxSize
		to := from + TextRecordMaxSize

		record := NewTBSTextRecord(html[from:min(to, len(html))])
		records = append(records, record)
	}

	len := records[len(records)-1].Length()
	if len%4 != 0 {
		pad := make(pdb.RawRecord, len%4)
		records = append(records, pad)
	}

	return records
}

func encodeVwi(x int) []byte {
	buf := make([]byte, 64)
	z := 0
	for {
		buf[z] = byte(x) & 0x7f
		x >>= 7
		z++
		if x == 0 {
			buf[0] |= 0x80
			break
		}
	}

	relevant := buf[:z]
	reverseBytes(relevant)
	return relevant
}

func reverseBytes(buf []byte) {
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
}

func invMod(dividend int, divisor int) int {
	return (divisor/2 + dividend) % divisor
}

func writeSequential(w io.Writer, bo binary.ByteOrder, vs ...interface{}) error {
	for _, v := range vs {
		err := binary.Write(w, bo, v)
		if err != nil {
			return err
		}
	}
	return nil
}
