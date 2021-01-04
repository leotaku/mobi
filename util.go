package mobi

import (
	"strings"

	"github.com/leotaku/mobi/pdb"
	r "github.com/leotaku/mobi/records"
)

func chaptersToText(m Book) (string, []r.ChunkInfo, []r.ChapterInfo, error) {
	text := new(strings.Builder)
	chunks := make([]r.ChunkInfo, 0)
	chaps := make([]r.ChapterInfo, 0)
	if m.tpl == nil {
		m.tpl = defaultTemplate
	}

	for I, chap := range m.Chapters {
		for i, chunk := range chap.Chunks {
			inv := newInventory(m, chap, I, I+i)
		chapStart := text.Len()
			head, err := runTemplate(*m.tpl, inv)
			if err != nil {
				return "", nil, nil, err
			}
			chunks = append(chunks, r.ChunkInfo{
				PreStart:      text.Len(),
				PreLength:     len(head),
				ContentStart:  text.Len() + len(head),
				ContentLength: len(chunk.Body),
			})
			text.WriteString(head)
			text.WriteString(chunk.Body)
		}
		chaps = append(chaps, r.ChapterInfo{
			Title:  chap.Title,
			Start:  chapStart,
			Length: text.Len() - chapStart,
		})
	}

	return text.String(), chunks, chaps, nil
}

func textToRecords(html string) []pdb.Record {
	records := []pdb.Record{}
	recordCount := len(html) / r.TextRecordMaxSize
	if len(html)%r.TextRecordMaxSize != 0 {
		recordCount++
	}

	for i := 0; i < recordCount; i++ {
		from := i * r.TextRecordMaxSize
		to := from + r.TextRecordMaxSize

		record := r.NewTextRecord(html[from:min(to, len(html))])
		records = append(records, record)
	}

	return records
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
