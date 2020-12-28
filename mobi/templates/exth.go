package templates

const EXTHEntryHeaderLength = 8 // 0x08

type EXTHEntryHeader struct {
	RecordType   uint32
	RecordLength uint32
}

func NewEXTHEntryHeader(RecordType uint32, RecordLength uint32) EXTHEntryHeader {
	return EXTHEntryHeader{
		RecordType:   RecordType,
		RecordLength: RecordLength,
	}

}

const EXTHHeaderLength = 12 // 0x0C

type EXTHHeader struct {
	EXTH         [4]byte
	HeaderLength uint32
	EntryCount   uint32
}

func NewEXTHHeader(HeaderLength uint32, EntryCount uint32) EXTHHeader {
	return EXTHHeader{
		EXTH:         [4]byte{'E', 'X', 'T', 'H'},
		HeaderLength: HeaderLength,
		EntryCount:   EntryCount,
	}
}
