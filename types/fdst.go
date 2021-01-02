package types

const FDSTHeaderLength = 12 // 0x0C

type FDSTHeader struct {
	FDST       [4]byte
	Fixed      uint32
	EntryCount uint32
}

func NewFDSTHeader() FDSTHeader {
	return FDSTHeader{
		FDST:       [4]byte{'F', 'D', 'S', 'T'},
		Fixed:      12,
		EntryCount: 0,
	}
}

const FDSTEntryLength = 8 // 0x08

type FDSTEntry struct {
	Start uint32
	End   uint32
}
