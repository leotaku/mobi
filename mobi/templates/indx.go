package templates

const (
	CBSingle byte = 15  // 0x0F
	CBParent      = 111 // 0x6F
	CBChild       = 31  // 0x1F
)

const INDXHeaderLength = 192 // 0xC0

type INDXHeader struct {
	INDX             [4]byte
	HeaderLength     uint32
	Unknown1         [4]byte
	HeaderType       uint32
	IndexType        uint32
	IDXTStart        uint32
	IndexRecordCount uint32
	IndexEncoding    uint32
	IndexLanguage    uint32
	IndexEntryCount  uint32
	ORDTStart        uint32
	LIGTStart        uint32
	LIGTCount        uint32
	CNCXCount        uint32
	Unknown2         [124]byte
	TAGXOffset       uint32
	Unknown3         [8]byte
}

func NewINDXHeader(RecordCount uint32, EntryCount uint32) INDXHeader {
	return INDXHeader{
		INDX:             [4]byte{'I', 'N', 'D', 'X'},
		HeaderLength:     INDXHeaderLength,
		HeaderType:       0,
		IndexType:        2,   // TODO: normal, inflection
		IDXTStart:        232, // TODO 240
		IndexRecordCount: RecordCount,
		IndexEncoding:    65001,
		IndexLanguage:    0,
		IndexEntryCount:  EntryCount,
		ORDTStart:        0,
		LIGTStart:        0,
		LIGTCount:        0,
		CNCXCount:        1, // TODO
		TAGXOffset:       INDXHeaderLength,
	}
}

const TAGXSingleHeaderLength = 32 // 0x20

type TAGXSingleHeader struct {
	TAGX             [4]byte
	HeaderLength     uint32
	ControlByteCount uint32
	TagTable         [5]TAGXTag
}

type TAGXTag uint32

const (
	TAGXTagEntryPosition   TAGXTag = 0x01010100 // 01, 1, 001, 0
	TAGXTagEntryLength             = 0x02010200 // 02, 1, 002, 0
	TAGXTagEntryNameOffset         = 0x03010400 // 03, 1, 004, 0
	TAGXTagEntryDepthLevel         = 0x04010800 // 04, 1, 008, 0
	TAGXTagEntryParent             = 0x15011000 // 21, 1, 016, 0
	TAGXTagEntryChild1             = 0x16012000 // 22, 1, 032, 0
	TAGXTagEntryChildN             = 0x17014000 // 23, 1, 064, 0
	TAGXTagEntryPosFid             = 0x06028000 // 06, 2, 128, 0
	TAGXTagEntryEnd                = 0x00000001 // 00, 0, 000, 1
)

func NewTAGXSingleHeader() TAGXSingleHeader {
	return TAGXSingleHeader{
		TAGX:             [4]byte{'T', 'A', 'G', 'X'},
		HeaderLength:     TAGXSingleHeaderLength,
		ControlByteCount: 1,
		TagTable: [5]TAGXTag{
			TAGXTagEntryPosition,
			TAGXTagEntryLength,
			TAGXTagEntryNameOffset,
			TAGXTagEntryDepthLevel,
			TAGXTagEntryEnd,
		},
	}
}

const IDXTHeaderLength = 6 // 0x06

type IDXTHeader struct {
	IDXT   [4]byte
	Offset uint16
}

func NewIDXTHeader(Offset uint16) IDXTHeader {
	return IDXTHeader{
		IDXT:   [4]byte{'I', 'D', 'X', 'T'},
		Offset: Offset,
	}
}
