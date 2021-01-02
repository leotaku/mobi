package templates

import "math"

const PalmDocHeaderLength = 16 // 0x10

type PalmDocHeader struct {
	Compression     uint16
	Unused1         uint16
	TextLength      uint32
	TextRecordCount uint16
	RecordSize      uint16
	Encryption      uint16
	Unknown1        uint16
}

func NewPalmDocHeader() PalmDocHeader {
	return PalmDocHeader{
		Compression:     1, // TODO
		Unused1:         0,
		TextLength:      0,
		TextRecordCount: 0,
		RecordSize:      0x1000,
		Encryption:      0, // TODO
		Unknown1:        0,
	}
}

const MOBIHeaderLength = 232 // 0xE8

type MOBIHeader struct {
	MOBI                                    [4]byte
	HeaderLength                            uint32
	MOBIType                                uint32
	TextEncoding                            uint32
	UniqueID                                uint32
	FileVersion                             uint32
	OrthographicIndex                       uint32
	InflectionIndex                         uint32
	IndexNames                              uint32
	IndexKeys                               uint32
	ExtraIndex0                             uint32
	ExtraIndex1                             uint32
	ExtraIndex2                             uint32
	ExtraIndex3                             uint32
	ExtraIndex4                             uint32
	ExtraIndex5                             uint32
	FirstNonBookIndex                       uint32
	FullNameOffset                          uint32
	FullNameLength                          uint32
	Locale                                  uint32
	InputLanguage                           uint32
	OutputLanguage                          uint32
	MinVersion                              uint32
	FirstImageIndex                         uint32
	HuffmanRecordOffset                     uint32
	HuffmanRecordCount                      uint32
	HuffmanTableOffset                      uint32 // DATP
	HuffmanTableLength                      uint32 // DATP
	EXTHFlags                               uint32
	Unknown1                                [32]byte
	DRMOffset                               uint32
	DRMCount                                uint32
	DRMSize                                 uint32
	DRMFlags                                uint32
	Unknown2                                [12]byte
	FirstContentRecordNumberOrFDSTNumberMSB uint16
	LastContentRecordNumberOrFDSTNumberLSB  uint16
	Unknown3OrFDSTEntryCount                uint32
	FCISRecordNumber                        uint32
	FCISRecordCount                         uint32
	FLISRecordNumber                        uint32
	FLISRecordCount                         uint32
	Unknown4                                uint64
	Unknown5                                uint32
	FirstCompilationSectionCount            uint32
	CompilationSectionCount                 uint32
	Unknown6                                uint32
	ExtraRecordDataFlags                    uint32
	INDXRecordOffset                        uint32
}

func NewMOBIHeader() MOBIHeader {
	return MOBIHeader{
		MOBI:                                    [4]byte{'M', 'O', 'B', 'I'},
		HeaderLength:                            MOBIHeaderLength,
		MOBIType:                                2,     // Book
		TextEncoding:                            65001, // Unicode
		UniqueID:                                0,
		FileVersion:                             6,
		OrthographicIndex:                       math.MaxUint32,
		InflectionIndex:                         math.MaxUint32,
		IndexNames:                              math.MaxUint32,
		IndexKeys:                               math.MaxUint32,
		ExtraIndex0:                             math.MaxUint32,
		ExtraIndex1:                             math.MaxUint32,
		ExtraIndex2:                             math.MaxUint32,
		ExtraIndex3:                             math.MaxUint32,
		ExtraIndex4:                             math.MaxUint32,
		ExtraIndex5:                             math.MaxUint32,
		FirstNonBookIndex:                       0,
		FullNameOffset:                          0,
		FullNameLength:                          0,
		Locale:                                  0, // NEUTRAL-NEUTRAL
		InputLanguage:                           0,
		OutputLanguage:                          0,
		MinVersion:                              6,
		FirstImageIndex:                         math.MaxUint32,
		HuffmanRecordOffset:                     0, // TODO
		HuffmanRecordCount:                      0,
		HuffmanTableOffset:                      0,
		HuffmanTableLength:                      0,
		EXTHFlags:                               0b1010000, // TODO
		Unknown1:                                [32]byte{},
		DRMOffset:                               math.MaxUint32,
		DRMCount:                                math.MaxUint32,
		DRMSize:                                 0,
		DRMFlags:                                0,
		Unknown2:                                [12]byte{},
		FirstContentRecordNumberOrFDSTNumberMSB: 1,
		LastContentRecordNumberOrFDSTNumberLSB:  1,
		Unknown3OrFDSTEntryCount:                1,
		FCISRecordNumber:                        0,
		FCISRecordCount:                         0,
		FLISRecordNumber:                        0,
		FLISRecordCount:                         0,
		Unknown4:                                0,
		Unknown5:                                math.MaxUint32,
		FirstCompilationSectionCount:            0,
		CompilationSectionCount:                 math.MaxUint32,
		Unknown6:                                math.MaxUint32,
		ExtraRecordDataFlags:                    0b01, // Should be '0b11' with TBS
		INDXRecordOffset:                        math.MaxUint32,
	}
}

const KF8HeaderLength = MOBIHeaderLength + 32

type KF8Header struct {
	MOBIHeader
	ChunkIndex        uint32 // SECT
	SkeletonIndex     uint32 // SKEL
	HuffmanTableIndex uint32 // DATP
	GuideIndex        uint32 // OTH
	Unknown           [16]byte
}

func NewKF8Header() KF8Header {
	mh := NewMOBIHeader()
	mh.FileVersion = 8
	mh.MinVersion = 8
	mh.HeaderLength = KF8HeaderLength
	mh.FirstContentRecordNumberOrFDSTNumberMSB = math.MaxUint16
	mh.LastContentRecordNumberOrFDSTNumberLSB = math.MaxUint16
	mh.Unknown3OrFDSTEntryCount = 0
	return KF8Header{
		MOBIHeader:        mh,
		ChunkIndex:        math.MaxUint32,
		SkeletonIndex:     math.MaxUint32,
		HuffmanTableIndex: math.MaxUint32,
		GuideIndex:        math.MaxUint32,
	}
}
