// Package mobi implements writing KF8-style formatted MOBI and AZW3 books.
package mobi

import (
	"fmt"
	"image"
	"strings"
	"text/template"
	"time"

	"github.com/leotaku/mobi/pdb"
	r "github.com/leotaku/mobi/records"
	t "github.com/leotaku/mobi/types"
	"golang.org/x/text/language"
)

// Book represents all the information necessary to generate a
// KF8-style formatted MOBI or AZW3 book.
//
// You are expected to initialize the Book struct using the required
// variables and/or builder pattern, then convert the resulting
// structure into a PalmDB database.  This database can then be
// written out to any io.Writer.
type Book struct {
	Title         string
	Authors       []string
	Contributors  []string
	Publisher     string
	Subject       string
	CreatedDate   time.Time
	PublishedDate time.Time
	DocType       string
	Language      language.Tag
	FixedLayout   bool
	RightToLeft   bool
	Chapters      []Chapter
	CSSFlows      []string
	Images        []image.Image
	CoverImage    image.Image
	ThumbImage    image.Image
	UniqueID      uint32

	// hidden
	tpl *template.Template
}

// OverrideTemplate overrides the template used in order to generate
// the skeleton section of a KF8 HTML chunk.
//
// During conversion to a PalmDB database, this template is passed the
// internal inventory type.  If the template cannot successfully be
// applied, the conversion will panic.
//
// The skeleton section generally consists of a complete HTML document
// including head and body, with the body tag expected to contain an
// 'aid' attribute that indicates the identifier of the corresponding
// chunk.  As it is relatively easy to end up with an invalid KF8
// document by generating invalid skeleton sections, this option is
// private and hidden behind a setter function.
func (m *Book) OverrideTemplate(tpl template.Template) Book {
	m.tpl = &tpl
	return *m
}

// GetThumbFilename returns the filename a JPEG image in the Kindle
// thumbnails folder would need to have in order to be correctly
// associated with the generated book.
//
// Newer Kindle readers no longer respect the thumbnail embedded in
// the Palm database but instead try to download the cover art from
// Amazon's servers using the embedded ASIN identification string.
// However, we can still get covers for sideloaded books by embedding
// a fake ASIN identification string and manually copying an image to
// the corresponding location inside the Kindle thumbnails folder.
func (m Book) GetThumbFilename() string {
	asin := encodeASIN(m.UniqueID)
	return fmt.Sprintf("thumbnail_%v_EBOK_portrait.jpg", asin)
}

// Chapter represents a chapter in a Book.
type Chapter struct {
	Title  string
	Chunks []Chunk
}

// Chunk represents a chunk of text in a Chapter.
//
// Chunks are mostly an implementation detail that is exposed for
// maximum control over the final book output.  Generally, you should
// use one of the various "Chunks" functions in order to generate the
// correct amount of chunks for a chapter.
type Chunk struct {
	Body string
}

// Realize converts a Book to a PalmDB Database.
func (m Book) Realize() pdb.Database {
	db := pdb.NewDatabase(m.Title, m.CreatedDate)
	html, chunks, chaps, err := chaptersToText(m)
	text := html + strings.Join(m.CSSFlows, "")
	textRecords := textToRecords(text, chaps)

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
	ch := r.ChunkHeaderIndexRecord(len(text), len(chunks))
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
	nh := r.NCXHeaderIndexRecord(len(chaps))
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
	flows := append([]string{html}, m.CSSFlows...)
	db.AddRecord(r.NewFDSTRecord(flows...))
	null.MOBIHeader.Unknown3OrFDSTEntryCount = uint32(len(m.CSSFlows) + 1)
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

func (m Book) createNullRecord() r.NullRecord {
	// Variables
	null := r.NewNullRecord(m.Title)
	lastImageID := len(m.Images)
	null.MOBIHeader.UniqueID = m.UniqueID
	null.MOBIHeader.Locale = matchLocale(m.Language)

	// EXTH header
	lang, _ := m.Language.Base()
	null.EXTHSection.AddString(t.EXTHTitle, m.Title)
	null.EXTHSection.AddString(t.EXTHUpdatedTitle, m.Title)
	null.EXTHSection.AddString(t.EXTHAuthor, m.Authors...)
	null.EXTHSection.AddString(t.EXTHContributor, m.Contributors...)
	null.EXTHSection.AddString(t.EXTHPublisher, m.Publisher)
	null.EXTHSection.AddString(t.EXTHSubject, m.Subject)
	null.EXTHSection.AddString(t.EXTHASIN, encodeASIN(m.UniqueID))
	null.EXTHSection.AddString(t.EXTHLanguage, lang.String())
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
		null.EXTHSection.AddInt(t.EXTHCoverOffset, lastImageID)
		null.EXTHSection.AddInt(t.EXTHHasFakeCover, 0)
		null.EXTHSection.AddString(t.EXTHKF8CoverURI, fmt.Sprintf("kindle:embed:%v", r.To32(lastImageID+1)))
		lastImageID++
	}
	if m.ThumbImage != nil {
		null.EXTHSection.AddInt(t.EXTHThumbOffset, lastImageID)
	}

	return null
}

func encodeASIN(id uint32) string {
	return fmt.Sprintf("%015x", id)
}
