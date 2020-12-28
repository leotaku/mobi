package mobi

import (
	"encoding/binary"
	"io"

	"github.com/leotaku/manki/mobi/pdb"
	t "github.com/leotaku/manki/mobi/templates"
)

type EXTHEntryType = uint32

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
	EXTHPublishingdate                     = 106
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
	EXTHCountResources                     = 125
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
	EXTHDoctype                            = 501
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

type EXTHSection struct {
	entries []EXTHEntry
}

func NewEXTHSection() EXTHSection {
	return EXTHSection{
		entries: []EXTHEntry{},
	}
}

func (sec *EXTHSection) AddString(tp EXTHEntryType, s string) {
	data := []byte(s)
	for i, entry := range sec.entries {
		if entry.RecordType == tp {
			sec.entries[i].Data = data
			return
		}
	}
	entry := NewEXTHEntry(tp, data)
	sec.entries = append(sec.entries, entry)
}

func (sec *EXTHSection) AddInt(tp EXTHEntryType, i int) {
	data := make([]byte, 4)
	pdb.Endian.PutUint32(data, uint32(i))
	for i, entry := range sec.entries {
		if entry.RecordType == tp {
			sec.entries[i].Data = data
			return
		}
	}
	entry := NewEXTHEntry(tp, data)
	sec.entries = append(sec.entries, entry)
}

func (e EXTHSection) LengthWithoutPadding() int {
	len := t.EXTHHeaderLength
	for _, entry := range e.entries {
		len += entry.Length()
	}

	return len
}

func (e EXTHSection) Length() int {
	len := e.LengthWithoutPadding()
	return len + invMod(len, 4)
}

func (e EXTHSection) Write(w io.Writer) error {
	lenNoPadding := e.LengthWithoutPadding()

	// Write fixed start of header
	h := t.NewEXTHHeader(uint32(lenNoPadding), uint32(len(e.entries)))
	err := binary.Write(w, pdb.Endian, h)
	if err != nil {
		return err
	}

	// Write entries
	for _, entry := range e.entries {
		err := entry.Write(w)
		if err != nil {
			return err
		}
	}

	// Write padding
	pad := make([]byte, invMod(lenNoPadding, 4))
	_, err = w.Write(pad)
	return err
}

type EXTHEntry struct {
	RecordType EXTHEntryType
	Data       []byte
}

func NewEXTHEntry(tp EXTHEntryType, data []byte) EXTHEntry {
	return EXTHEntry{
		RecordType: tp,
		Data:       data,
	}
}

func (e EXTHEntry) Length() int {
	return len(e.Data) + t.EXTHEntryHeaderLength
}

func (e EXTHEntry) Write(w io.Writer) error {
	h := t.NewEXTHEntryHeader(e.RecordType, uint32(e.Length()))
	err := binary.Write(w, pdb.Endian, h)
	if err != nil {
		return err
	}

	_, err = w.Write(e.Data)
	return err
}
