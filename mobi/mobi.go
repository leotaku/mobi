package mobi

import (
	"fmt"
	"image"

	"github.com/leotaku/manki/mobi/pdb"
	t "github.com/leotaku/manki/mobi/templates"
)

type MobiBook struct {
	Title      string
	Author     string
	Chapters   []Chapter
	Images     []image.Image
	CoverImage int
	UniqueID   uint32

	html        string
	textRecords []TBSTextRecord
}

type Chapter struct {
	Name        string
	TextContent string
}

func (m MobiBook) Realize() pdb.Database {
	db := pdb.NewDatabase(m.Title)
	m.html = chaptersToText(m.Chapters)
	m.textRecords = genTextRecords(m.html)

	// Null record
	null := m.createNullRecord()
	db.AddRecord(null)

	// Text records
	for _, rec := range m.textRecords {
		db.AddRecord(rec)
	}

	// Index and CNCX records
	first, second, cncx := createIndexRecords(preChapOffset, m.Chapters)
	db.AddRecord(first)
	db.AddRecord(second)
	db.AddRecord(cncx)

	// Image records
	for _, img := range m.Images {
		rec := NewImageRecord(img)
		db.AddRecord(rec)
	}

	// Trailing records
	db.AddRecord(t.NewFLISRecord())
	db.AddRecord(t.NewFCISRecord(uint32(len(m.html))))
	db.AddRecord(t.EOFRecord)

	return db
}

func (m *MobiBook) createNullRecord() NullRecord {
	// Variables
	null := NewNullRecord(m.Title)
	ltr := len(m.textRecords)
	lir := ltr + 3
	lcr := ltr + len(m.Images)

	// PalmDoc header
	null.PalmDocHeader.TextLength = uint32(len(m.html))
	null.PalmDocHeader.TextRecordCount = uint16(ltr)

	// MOBI header
	null.MOBIHeader.LastContentRecordNumber = uint16(ltr)
	if len(m.Images) > 0 {
		null.MOBIHeader.FirstImageIndex = uint32(lir + 1)
	}
	null.MOBIHeader.FirstNonBookIndex = uint32(ltr + 1)
	null.MOBIHeader.INDXRecordOffset = uint32(ltr + 1)
	null.MOBIHeader.FLISRecordCount = 1
	null.MOBIHeader.LastContentRecordNumber = uint16(lcr + 1)
	null.MOBIHeader.FLISRecordNumber = uint32(lcr + 4)
	null.MOBIHeader.FCISRecordCount = 1
	null.MOBIHeader.FCISRecordNumber = uint32(lcr + 5)
	null.MOBIHeader.UniqueID = m.UniqueID

	// EXTH header
	null.EXTHSection.AddString(EXTHTitle, m.Title)
	null.EXTHSection.AddString(EXTHAuthor, m.Author)
	null.EXTHSection.AddString(EXTHUpdatedTitle, m.Title)
	if m.CoverImage <= len(m.Images) {
		null.EXTHSection.AddInt(EXTHCoverOffset, m.CoverImage)
		null.EXTHSection.AddInt(EXTHThumbOffset, m.CoverImage)
		null.EXTHSection.AddString(EXTHKF8CoverURI, fmt.Sprintf("kindle:embed:%03v", m.CoverImage+1))
	}

	return null
}
