package templates

type EXTHEntryType uint32

const (
	EXTHDrmserver            EXTHEntryType = 1
	EXTHDrmcommerce                        = 2
	EXTHDrmebookbase                       = 3
	EXTHTitle                              = 99
	EXTHAuthor                             = 100
	EXTHPublisher                          = 101
	EXTHImprint                            = 102
	EXTHDescription                        = 103
	EXTHIsbn                               = 104
	EXTHSubject                            = 105
	EXTHPublishingDate                     = 106
	EXTHReview                             = 107
	EXTHContributor                        = 108
	EXTHRights                             = 109
	EXTHSubjectCode                        = 110
	EXTHType_                              = 111
	EXTHSource                             = 112
	EXTHAsin                               = 113
	EXTHVersion                            = 114
	EXTHSample                             = 115
	EXTHStartReading                       = 116
	EXTHAdult                              = 117
	EXTHPrice                              = 118
	EXTHCurrency                           = 119
	EXTHKF8Boundary                        = 121
	EXTHFixedLayout                        = 122
	EXTHBookType                           = 123
	EXTHOrientationLock                    = 124
	EXTHKF8CountResources                  = 125
	EXTHOrigResolution                     = 126
	EXTHZeroGutter                         = 127
	EXTHZeroMargin                         = 128
	EXTHKF8CoverURI                        = 129
	EXTHKF8UnidentifiedCount               = 131
	EXTHRegionMagni                        = 132
	EXTHDictName                           = 200
	EXTHCoverOffset                        = 201
	EXTHThumbOffset                        = 202
	EXTHHasFakeCover                       = 203
	EXTHCreatorSoftware                    = 204
	EXTHCreatorMajor                       = 205
	EXTHCreatorMinor                       = 206
	EXTHCreatorBuild                       = 207
	EXTHWatermark                          = 208
	EXTHTamperKeys                         = 209
	EXTHFontSignature                      = 300
	EXTHClippingLimit                      = 401
	EXTHPublisherLimit                     = 402
	EXTHUnknown403                         = 403
	EXTHTtsDisable                         = 404
	EXTHUnknown405                         = 405
	EXTHRental                             = 406
	EXTHUnknown407                         = 407
	EXTHUnknown450                         = 450
	EXTHUnknown451                         = 451
	EXTHUnknown452                         = 452
	EXTHUnknown453                         = 453
	EXTHDocType                            = 501
	EXTHLastUpdate                         = 502
	EXTHUpdatedTitle                       = 503
	EXTHAsin504                            = 504
	EXTHTitleFileAs                        = 508
	EXTHCreatorFileAs                      = 517
	EXTHPublisherFileAs                    = 522
	EXTHLanguage                           = 524
	EXTHAlignment                          = 525
	EXTHPagedir                            = 527
	EXTHOverrideFonts                      = 528
	EXTHSorceDescription                   = 529
	EXTHDictLangInput                      = 531
	EXTHDictLangOutput                     = 532
	EXTHUnknown534                         = 534
	EXTHCreatorBuildRev                    = 535
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
