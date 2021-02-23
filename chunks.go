package mobi

import "github.com/leotaku/mobi/records"

// SingleChunks produces a list of chunks from several strings.  In
// the resulting list of chunks, each chunk exactly corresponds the
// given string with the same index.
//
// This is useful in cases where the size of a chunk has influence on
// the presentation of the generated book.  In the particular case of
// fixed-layout books, a single chunk also represents a single page.
func SingleChunks(ss ...string) []Chunk {
	result := make([]Chunk, 0)
	for _, s := range ss {
		result = append(result, Chunk{
			Body: s,
		})
	}
	return result
}

// Chunks splits a string into a list of chunks sized at most 4096
// bytes.  In most cases, this is the preferred method of converting
// book text into a list of chunks.
//
// Using this method, chunks should be very roughly aligned to the
// maximum text record size, which may or may not be beneficial for
// performance.
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
