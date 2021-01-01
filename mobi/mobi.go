package mobi

import (
	"image"
	"math"
	"strings"

	"github.com/leotaku/manki/mobi/pdb"
	t "github.com/leotaku/manki/mobi/templates"
)

type MobiBook struct {
	Title      string
	Author     string
	Chapters   []Chapter
	ExtraFlows []string
	Images     []image.Image
	CoverImage int
	UniqueID   uint32

	text        string
	textRecords []pdb.Record
}

type Chapter struct {
	Name        string
	TextContent string
}

func (m MobiBook) Realize() pdb.Database {
	db := pdb.NewDatabase(m.Title)
	html, info := chaptersToText(m.Chapters)
	m.text = html + strings.Join(m.ExtraFlows, "")
	m.textRecords = genTextRecords(m.text)

	// Null record
	null := m.createNullRecord()
	db.AddRecord(null)

	// Text records
	for _, rec := range m.textRecords {
		db.AddRecord(rec)
	}

	// Chunk record
	chunk, cncx := ChunkIndexRecord(info)
	ch := ChunkHeaderIndexRecord(len(m.text), len(m.Chapters))
	db.AddRecord(ch)
	db.AddRecord(chunk)
	db.AddRecord(cncx)

	// Skeleton record
	skeleton := SkeletonIndexRecord(info)
	sh := SkeletonHeaderIndexRecord(len(skeleton.IDXTEntries))
	db.AddRecord(sh)
	db.AddRecord(skeleton)

	// Image records
	for _, img := range m.Images {
		rec := NewImageRecord(img)
		db.AddRecord(rec)
	}

	// Trailing records
	flows := append([]string{html}, m.ExtraFlows...)
	db.AddRecord(NewFDSTRecord(flows...))
	db.AddRecord(t.NewFLISRecord())
	db.AddRecord(t.NewFCISRecord(uint32(len(m.text))))
	db.AddRecord(t.EOFRecord)

	return db
}

func (m *MobiBook) createNullRecord() NullRecord {
	// Variables
	null := NewNullRecord(m.Title)
	ltr := len(m.textRecords)  // Last text record
	lir := ltr + 5             // Last index record
	lcr := lir + len(m.Images) // Last contentful record

	// PalmDoc header
	null.PalmDocHeader.TextLength = uint32(len(m.text))
	null.PalmDocHeader.TextRecordCount = uint16(ltr - 1)

	// MOBI header
	if len(m.Images) > 0 {
		null.MOBIHeader.FirstImageIndex = uint32(lir + 1)
	} else {
		null.MOBIHeader.FirstImageIndex = math.MaxUint32
	}
	null.MOBIHeader.FirstNonBookIndex = uint32(ltr + 1)
	null.MOBIHeader.FLISRecordCount = 1
	null.MOBIHeader.FLISRecordNumber = uint32(lcr + 2)
	null.MOBIHeader.FCISRecordCount = 1
	null.MOBIHeader.FCISRecordNumber = uint32(lcr + 3)
	null.MOBIHeader.UniqueID = m.UniqueID

	// KF8
	null.MOBIHeader.FirstContentRecordNumberOrFDSTNumberMSB = 0
	null.MOBIHeader.LastContentRecordNumberOrFDSTNumberLSB = uint16(lcr + 1)
	null.MOBIHeader.Unknown3OrFDSTEntryCount = uint32(len(m.ExtraFlows) + 1)
	null.MOBIHeader.HuffmanTableIndex = math.MaxUint32

	// Indexes
	null.MOBIHeader.INDXRecordOffset = math.MaxUint32
	null.MOBIHeader.ChunkIndex = uint32(ltr + 1)
	null.MOBIHeader.SkeletonIndex = uint32(ltr + 4)
	null.MOBIHeader.GuideIndex = math.MaxUint32

	// EXTH header
	// null.EXTHSection.AddString(EXTHTitle, m.Title)
	null.EXTHSection.AddString(EXTHAuthor, m.Author)
	null.EXTHSection.AddString(EXTHUpdatedTitle, m.Title)
	null.EXTHSection.AddString(EXTHDocType, "EBOK")
	null.EXTHSection.AddString(EXTHLanguage, "en")
	null.EXTHSection.AddString(EXTHPublishingDate, "2020-12-31T14:45:45.113077+00:00")
	null.EXTHSection.AddString(EXTHContributor, "calibre (4.23.0) [https://calibre-ebook.com]")
	null.EXTHSection.AddString(EXTHSource, "calibre:1c984ebb-d396-46c8-9773-456f57e012de")
	null.EXTHSection.AddString(EXTHAsin, "1c984ebb-d396-46c8-9773-456f57e012de")
	null.EXTHSection.AddInt(EXTHKF8CountResources, 0)
	null.EXTHSection.AddInt(EXTHCreatorSoftware, 201)
	null.EXTHSection.AddInt(EXTHCreatorMajor, 2)
	null.EXTHSection.AddInt(EXTHCreatorMinor, 9)
	null.EXTHSection.AddInt(EXTHCreatorBuild, 0)
	null.EXTHSection.AddString(EXTHCreatorBuildRev, "0730-890adc2")
	if m.CoverImage <= len(m.Images) {
		// null.EXTHSection.AddInt(EXTHCoverOffset, m.CoverImage)
		// null.EXTHSection.AddInt(EXTHThumbOffset, m.CoverImage)
		// null.EXTHSection.AddString(EXTHKF8CoverURI, fmt.Sprintf("kindle:embed:%03v", m.CoverImage+1))
	}

	return null
}
