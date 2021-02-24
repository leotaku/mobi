package mobi

// Chunks produces a list of chunks from one or more strings.
//
// In the resulting list of chunks, each chunk exactly corresponds the
// given string with the same index.  In almost all cases, this is the
// preferred method of converting KF8 HTML into a list of chunks.
func Chunks(ss ...string) []Chunk {
	result := make([]Chunk, 0)
	for _, s := range ss {
		result = append(result, Chunk{
			Body: s,
		})
	}
	return result
}
