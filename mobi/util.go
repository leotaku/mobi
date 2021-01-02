package mobi

import (
	"github.com/leotaku/manki/mobi/pdb"
	r "github.com/leotaku/manki/mobi/records"
)

func chaptersToText(m MobiBook) (string, []r.ChunkInfo, []r.ChapterInfo, error) {
	text := ""
	chunks := make([]r.ChunkInfo, 0)
	chaps := make([]r.ChapterInfo, 0)
	if m.tpl == nil {
		m.tpl = defaultTemplate
	}

	for I, chap := range m.Chapters {
		chapStart := len(text)
		for i, chunk := range chap.Chunks {
			inv := newInventory(m, chap, I, I+i)
			head, err := runTemplate(*m.tpl, inv)
			if err != nil {
				return "", nil, nil, err
			}
			chunks = append(chunks, r.ChunkInfo{
				PreStart:      len(text),
				PreLength:     len(head),
				ContentStart:  len(text) + len(head),
				ContentLength: len(chunk.Body),
			})
			text += head + chunk.Body
		}
		chaps = append(chaps, r.ChapterInfo{
			Title:  chap.Title,
			Start:  chapStart,
			Length: len(text) - chapStart,
		})
	}

	return text, chunks, chaps, nil
}

func textToRecords(html string) []pdb.Record {
	records := []pdb.Record{}
	recordCount := len(html) / r.TextRecordMaxSize
	if len(html)%r.TextRecordMaxSize != 0 {
		recordCount += 1
	}

	for i := 0; i < int(recordCount); i++ {
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
