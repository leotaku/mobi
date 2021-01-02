package types

import "math"

const (
	CBNCXSingle byte = 15  // 0x0F
	CBNCXParent      = 111 // 0x6F
	CBNCXChild       = 31  // 0x1F
	CBSkeleton       = 10  // 0x0A
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
		IndexType:        0,   // TODO: normal, inflection
		IDXTStart:        232, // TODO 240
		IndexRecordCount: RecordCount,
		IndexEncoding:    65001,
		IndexLanguage:    math.MaxUint32,
		IndexEntryCount:  EntryCount,
		ORDTStart:        0,
		LIGTStart:        0,
		LIGTCount:        0,
		CNCXCount:        1, // TODO
		TAGXOffset:       INDXHeaderLength,
	}
}

const TAGXHeaderLength = 12 // 0x0C

type TAGXHeader struct {
	TAGX             [4]byte
	HeaderLength     uint32
	ControlByteCount uint32
}

func NewTAGXHeader() TAGXHeader {
	return TAGXHeader{
		TAGX:             [4]byte{'T', 'A', 'G', 'X'},
		HeaderLength:     TAGXHeaderLength,
		ControlByteCount: 1,
	}
}

const TAGXSingleHeaderLength = 32 // 0x20

type TAGXSingleHeader struct {
	TAGX             [4]byte
	HeaderLength     uint32
	ControlByteCount uint32
	TagTable         [5]TAGXTag
}

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
			TAGXTagEnd,
		},
	}
}

const TAGXTagLength = 4 // 0x04

type TAGXTag uint32

const (
	TAGXTagEntryPosition       TAGXTag = 0x01010100
	TAGXTagEntryLength                 = 0x02010200
	TAGXTagEntryNameOffset             = 0x03010400
	TAGXTagEntryDepthLevel             = 0x04010800
	TAGXTagEntryParent                 = 0x15011000
	TAGXTagEntryChild1                 = 0x16012000
	TAGXTagEntryChildN                 = 0x17014000
	TAGXTagEntryPosFid                 = 0x06028000
	TAGXTagSkeletonChunkCount          = 0x01010300
	TAGXTagSkeletonGeometry            = 0x06020C00
	TAGXTagChunkCNCXOffset             = 0x02010100
	TAGXTagChunkFileNumber             = 0x03010200
	TAGXTagChunkSequenceNumber         = 0x04010400
	TAGXTagChunkGeometry               = 0x06020800
	TAGXTagGuideTitle                  = 0x01010100
	TAGXTagGuidePosFid                 = 0x06020200
	TAGXTagEnd                         = 0x00000001
)

type TAGXTagTable []TAGXTag

var TAGXTableNCXSingle = TAGXTagTable{
	TAGXTagEntryPosition,
	TAGXTagEntryLength,
	TAGXTagEntryNameOffset,
	TAGXTagEntryDepthLevel,
	TAGXTagEnd,
}

var TAGXTableSkeleton = TAGXTagTable{
	TAGXTagSkeletonChunkCount,
	TAGXTagSkeletonGeometry,
	TAGXTagEnd,
}

var TAGXTableChunk = TAGXTagTable{
	TAGXTagChunkCNCXOffset,
	TAGXTagChunkFileNumber,
	TAGXTagChunkSequenceNumber,
	TAGXTagChunkGeometry,
	TAGXTagEnd,
}

var TAGXTableGuide = TAGXTagTable{
	TAGXTagGuideTitle,
	TAGXTagGuidePosFid,
	TAGXTagEnd,
}

const IDXTSingleHeaderLength = 6 // 0x06

type IDXTSingleHeader struct {
	IDXT   [4]byte
	Offset uint16
}

func NewIDXTSingleHeader(Offset uint16) IDXTSingleHeader {
	return IDXTSingleHeader{
		IDXT:   [4]byte{'I', 'D', 'X', 'T'},
		Offset: Offset,
	}
}

const IDXTHeaderLength = 4 // 0x04

type IDXTHeader struct {
	IDXT [4]byte
}

func NewIDXTHeader() IDXTHeader {
	return IDXTHeader{
		IDXT: [4]byte{'I', 'D', 'X', 'T'},
	}
}
