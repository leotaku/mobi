package mobi

import (
	"fmt"
	"image"
	"strings"
	"time"

	"github.com/leotaku/manki/mobi/pdb"
	r "github.com/leotaku/manki/mobi/records"
	t "github.com/leotaku/manki/mobi/templates"
	"golang.org/x/text/language"
)

type MobiBook struct {
	Title         string
	Author        string
	Publisher     string
	Subject       string
	CreatedDate   time.Time
	PublishedDate time.Time
	DocType       string
	Language      language.Tag
	FixedLayout   bool
	RightToLeft   bool
	Chapters      []Chapter
	CssFlows      []string
	Images        []image.Image
	CoverImage    image.Image
	ThumbImage    image.Image
	UniqueID      uint32
}
type Chapter struct {
	Title  string
	Chunks []Chunk
}

type Chunk struct {
	Head string
	Body string
}

func (m MobiBook) Realize() pdb.Database {
	db := pdb.NewDatabase(m.Title, m.CreatedDate)
	html, chunks, chaps, err := chaptersToText(m)
	text := html + strings.Join(m.CssFlows, "")
	textRecords := genTextRecords(text)

	// Handle possible template errors
	if err != nil {
		panic(err)
	}

	// Null record
	null := m.createNullRecord()
	db.AddRecord(null)

	// Text records
	null.PalmDocHeader.TextRecordCount = uint16(len(textRecords))
	null.PalmDocHeader.TextLength = uint32(len(text))
	for _, rec := range textRecords {
		db.AddRecord(rec)
	}

	// Padding
	lastLength := textRecords[len(textRecords)-1].Length()
	if lastLength%4 != 0 {
		pad := make(pdb.RawRecord, lastLength%4)
		db.AddRecord(pad)
	}
	null.MOBIHeader.FirstNonBookIndex = uint32(db.Idx() + 1)

	// Chunk record
	chunk, cncx := r.ChunkIndexRecord(chunks)
	ch := r.ChunkHeaderIndexRecord(len(text), len(m.Chapters))
	null.MOBIHeader.ChunkIndex = uint32(db.AddRecord(ch))
	db.AddRecord(chunk)
	db.AddRecord(cncx)

	// Skeleton record
	skeleton := r.SkeletonIndexRecord(chunks)
	sh := r.SkeletonHeaderIndexRecord(len(skeleton.IDXTEntries))
	null.MOBIHeader.SkeletonIndex = uint32(db.AddRecord(sh))
	db.AddRecord(skeleton)

	// NCX record
	ncx, cncx := r.NCXIndexRecord(chaps)
	nh := r.NCXHeaderIndexRecord(len(skeleton.IDXTEntries))
	null.MOBIHeader.INDXRecordOffset = uint32(db.AddRecord(nh))
	db.AddRecord(ncx)
	db.AddRecord(cncx)

	// Image records
	images := m.Images
	if m.CoverImage != nil {
		images = append(images, m.CoverImage)
	}
	if m.ThumbImage != nil {
		images = append(images, m.ThumbImage)
	}
	if len(images) > 0 {
		null.MOBIHeader.FirstImageIndex = uint32(db.Idx() + 1)
		null.EXTHSection.AddInt(t.EXTHKF8CountResources, len(images))
	}
	for _, img := range images {
		rec := r.NewImageRecord(img)
		db.AddRecord(rec)
	}

	// FDST Record
	flows := append([]string{html}, m.CssFlows...)
	db.AddRecord(r.NewFDSTRecord(flows...))
	null.MOBIHeader.Unknown3OrFDSTEntryCount = uint32(len(m.CssFlows) + 1)
	null.MOBIHeader.FirstContentRecordNumberOrFDSTNumberMSB = 0
	null.MOBIHeader.LastContentRecordNumberOrFDSTNumberLSB = uint16(db.Idx())

	// FLIS Record
	db.AddRecord(t.NewFLISRecord())
	null.MOBIHeader.FLISRecordCount = 1
	null.MOBIHeader.FLISRecordNumber = uint32(db.Idx())

	// FCIS Record
	db.AddRecord(t.NewFCISRecord(uint32(len(text))))
	null.MOBIHeader.FCISRecordCount = 1
	null.MOBIHeader.FCISRecordNumber = uint32(db.Idx())

	// Replace updated Null record
	db.AddRecord(t.EOFRecord)
	db.ReplaceRecord(0, null)

	return db
}

func (m MobiBook) createNullRecord() r.NullRecord {
	// Variables
	null := r.NewNullRecord(m.Title)
	lastImageId := len(m.Images)
	null.MOBIHeader.UniqueID = m.UniqueID
	null.MOBIHeader.Locale = matchLocale(m.Language)

	// EXTH header
	langString := fmt.Sprint(m.Language)
	null.EXTHSection.AddString(t.EXTHTitle, m.Title)
	null.EXTHSection.AddString(t.EXTHUpdatedTitle, m.Title)
	null.EXTHSection.AddString(t.EXTHAuthor, m.Author)
	null.EXTHSection.AddString(t.EXTHPublisher, m.Publisher)
	null.EXTHSection.AddString(t.EXTHSubject, m.Subject)
	null.EXTHSection.AddString(t.EXTHASIN, encodeID(m.UniqueID))
	null.EXTHSection.AddString(t.EXTHLanguage, langString)
	if m.PublishedDate != (time.Time{}) {
		dateString := m.PublishedDate.Format("2006-01-02T15:04:05.000000+07:00")
		null.EXTHSection.AddString(t.EXTHPublishingDate, dateString)
	}
	if len(m.DocType) > 0 {
		null.EXTHSection.AddString(t.EXTHDocType, m.DocType)
	} else {
		null.EXTHSection.AddString(t.EXTHDocType, "EBOK")
	}
	if m.FixedLayout {
		null.EXTHSection.AddString(t.EXTHFixedLayout, "true")
	}
	if m.RightToLeft {
		null.EXTHSection.AddString(t.EXTHPrimaryWritingMode, "horizontal-rl")
		null.EXTHSection.AddString(t.EXTHPageProgressionDirection, "rtl")
	}
	if m.CoverImage != nil {
		null.EXTHSection.AddInt(t.EXTHCoverOffset, lastImageId)
		null.EXTHSection.AddInt(t.EXTHHasFakeCover, 0)
		null.EXTHSection.AddString(t.EXTHKF8CoverURI, fmt.Sprintf("kindle:embed:%04v", lastImageId+1))
		lastImageId += 1
	}
	if m.ThumbImage != nil {
		null.EXTHSection.AddInt(t.EXTHThumbOffset, lastImageId)
	}

	return null
}

func encodeID(id uint32) string {
	// b := make([]byte, 4)
	// pdb.Endian.PutUint32(b, id)
	// return base32.HexEncoding.EncodeToString(b)

	return fmt.Sprintf("%020v", id)
}

func watermark(null *r.NullRecord) {
	null.EXTHSection.AddString(t.EXTHContributor, "go-mobi [https://github.com/leotaku/go-mobi]")
	null.EXTHSection.AddString(t.EXTHSource, "go-mobi:00000000-0000-0000-0000-000000000000")
	null.EXTHSection.AddInt(t.EXTHCreatorSoftware, 202)
	null.EXTHSection.AddInt(t.EXTHCreatorMajor, 0)
	null.EXTHSection.AddInt(t.EXTHCreatorMinor, 0)
	null.EXTHSection.AddInt(t.EXTHCreatorBuild, 0)
	null.EXTHSection.AddString(t.EXTHCreatorBuildRev, "0000-0000000")
}
