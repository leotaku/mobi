package mobi

import (
	"bytes"
	"html/template"

	"github.com/leotaku/manki/mobi/pdb"
	r "github.com/leotaku/manki/mobi/records"
)

func chaptersToText(m MobiBook) (string, []r.ChunkInfo, []r.ChapterInfo, error) {
	text := ""
	chunks := make([]r.ChunkInfo, 0)
	chaps := make([]r.ChapterInfo, 0)

	for I, chap := range m.Chapters {
		chapStart := len(text)
		for i, chunk := range chap.Chunks {
			head, body, err := runChunk(chunk, m, I+i)
			if err != nil {
				return "", nil, nil, err
			}
			chunks = append(chunks, r.ChunkInfo{
				PreStart:      len(text),
				PreLength:     len(head),
				ContentStart:  len(text) + len(head),
				ContentLength: len(body),
			})
			text += head + body
		}
		chaps = append(chaps, r.ChapterInfo{
			Title:  chap.Title,
			Start:  chapStart,
			Length: len(text) - chapStart,
		})
	}

	return text, chunks, chaps, nil
}

func runChunk(c Chunk, m MobiBook, id int) (string, string, error) {
	inventory := struct {
		Mobi MobiBook
		Id   int
	}{
		Mobi: m,
		Id:   id,
	}

	head, err := runTemplate(c.Head, inventory)
	if err != nil {
		return "", "", err
	}

	body, err := runTemplate(c.Body, inventory)
	if err != nil {
		return "", "", err
	}

	return head, body, nil
}

var funcMap = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
	},
}

func runTemplate(s string, v interface{}) (string, error) {
	t, err := template.New("placeholder").Funcs(funcMap).Parse(s)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, v)
	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
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

	return records
}
