package mobi

func SingleChunk(ss ...string) []Chunk {
	result := make([]Chunk, 0)
	for _, s := range ss {
		result = append(result, Chunk{
			Body: s,
		})
	}
	return result
}

func Chunks(s string) []Chunk {
	panic("unimplemented")
}
