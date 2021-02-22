package mobi

import "github.com/leotaku/mobi/records"

func SingleChunks(ss ...string) []Chunk {
	result := make([]Chunk, 0)
	for _, s := range ss {
		result = append(result, Chunk{
			Body: s,
		})
	}
	return result
}

func Chunks(s string) []Chunk {
	result := make([]Chunk, 0)
	for len(s) > records.TextRecordMaxSize {
		body := s[:records.TextRecordMaxSize]
		s = s[records.TextRecordMaxSize:]
		result = append(result, Chunk{
			Body: body,
		})
	}

	return append(result, Chunk{
		Body: s,
	})
}
