package types

type EXTHEntryType uint32

const (
	EXTHDRMServer                EXTHEntryType = 1
	EXTHDRMCommerce              EXTHEntryType = 2
	EXTHDRMEbookBase             EXTHEntryType = 3
	EXTHTitle                    EXTHEntryType = 99
	EXTHAuthor                   EXTHEntryType = 100
	EXTHPublisher                EXTHEntryType = 101
	EXTHImprint                  EXTHEntryType = 102
	EXTHDescription              EXTHEntryType = 103
	EXTHISBN                     EXTHEntryType = 104
	EXTHSubject                  EXTHEntryType = 105
	EXTHPublishingDate           EXTHEntryType = 106
	EXTHReview                   EXTHEntryType = 107
	EXTHContributor              EXTHEntryType = 108
	EXTHRights                   EXTHEntryType = 109
	EXTHSubjectCode              EXTHEntryType = 110
	EXTHType_                    EXTHEntryType = 111
	EXTHSource                   EXTHEntryType = 112
	EXTHASIN                     EXTHEntryType = 113
	EXTHVersion                  EXTHEntryType = 114
	EXTHSample                   EXTHEntryType = 115
	EXTHStartReading             EXTHEntryType = 116
	EXTHAdult                    EXTHEntryType = 117
	EXTHPrice                    EXTHEntryType = 118
	EXTHCurrency                 EXTHEntryType = 119
	EXTHKF8Boundary              EXTHEntryType = 121
	EXTHFixedLayout              EXTHEntryType = 122
	EXTHBookType                 EXTHEntryType = 123
	EXTHOrientationLock          EXTHEntryType = 124
	EXTHKF8CountResources        EXTHEntryType = 125
	EXTHOrigResolution           EXTHEntryType = 126
	EXTHZeroGutter               EXTHEntryType = 127
	EXTHZeroMargin               EXTHEntryType = 128
	EXTHKF8CoverURI              EXTHEntryType = 129
	EXTHKF8UnidentifiedCount     EXTHEntryType = 131
	EXTHRegionMagnification      EXTHEntryType = 132
	EXTHDictName                 EXTHEntryType = 200
	EXTHCoverOffset              EXTHEntryType = 201
	EXTHThumbOffset              EXTHEntryType = 202
	EXTHHasFakeCover             EXTHEntryType = 203
	EXTHCreatorSoftware          EXTHEntryType = 204
	EXTHCreatorMajor             EXTHEntryType = 205
	EXTHCreatorMinor             EXTHEntryType = 206
	EXTHCreatorBuild             EXTHEntryType = 207
	EXTHWatermark                EXTHEntryType = 208
	EXTHTamperKeys               EXTHEntryType = 209
	EXTHFontSignature            EXTHEntryType = 300
	EXTHClippingLimit3XX         EXTHEntryType = 301
	EXTHClippingLimit            EXTHEntryType = 401
	EXTHPublisherLimit           EXTHEntryType = 402
	EXTHTtsDisable               EXTHEntryType = 404
	EXTHRental                   EXTHEntryType = 406
	EXTHDocType                  EXTHEntryType = 501
	EXTHLastUpdate               EXTHEntryType = 502
	EXTHUpdatedTitle             EXTHEntryType = 503
	EXTHASIN5XX                  EXTHEntryType = 504
	EXTHTitleFurigana            EXTHEntryType = 508
	EXTHCreatorFurigana          EXTHEntryType = 517
	EXTHPublisherFurigana        EXTHEntryType = 522
	EXTHLanguage                 EXTHEntryType = 524
	EXTHPrimaryWritingMode       EXTHEntryType = 525
	EXTHPageProgressionDirection EXTHEntryType = 527
	EXTHOverrideFonts            EXTHEntryType = 528
	EXTHSourceDescription        EXTHEntryType = 529
	EXTHDictLangInput            EXTHEntryType = 531
	EXTHDictLangOutput           EXTHEntryType = 532
	EXTHInputSourceType          EXTHEntryType = 534
	EXTHCreatorBuildRev          EXTHEntryType = 535
	EXTHContainerInfo            EXTHEntryType = 536
	EXTHContainerResolution      EXTHEntryType = 538
	EXTHContainerMimetype        EXTHEntryType = 539
	EXTHContainerID              EXTHEntryType = 543
)

const EXTHEntryHeaderLength = 8 // 0x08

type EXTHEntryHeader struct {
	RecordType   EXTHEntryType
	RecordLength uint32
}

func NewEXTHEntryHeader(RecordType EXTHEntryType, RecordLength uint32) EXTHEntryHeader {
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
