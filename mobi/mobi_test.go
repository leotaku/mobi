package mobi_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/leotaku/manki/mobi/pdb"
	"github.com/leotaku/manki/mobi/records"
	"github.com/leotaku/manki/mobi/templates"
)

func TestPDBHeaderLength(t *testing.T) {
	len := measure(pdb.PalmDBHeader{})
	assertEq(t, len, pdb.PalmDBHeaderLength)
}

func TestRecordHeaderLength(t *testing.T) {
	len := measure(pdb.RecordHeader{})
	assertEq(t, len, pdb.RecordHeaderLength)
}

func TestNullRecordLength(t *testing.T) {
	nr := records.NewNullRecord("Foo")
	buf := bytes.NewBuffer(nil)
	nr.Write(buf)

	assertEq(t, nr.Length(), len(buf.Bytes()))
}

func TestIndexSectionLength(t *testing.T) {
	indx := templates.NewINDXHeader(0, 0)
	tagx := templates.NewTAGXSingleHeader()
	idtx := templates.NewIDXTSingleHeader(0)

	assertEq(t, measure(indx), templates.INDXHeaderLength)
	assertEq(t, measure(tagx), templates.TAGXSingleHeaderLength)
	assertEq(t, measure(idtx), templates.IDXTSingleHeaderLength)
}

func TestNullRecordLengthWithEXTH(t *testing.T) {
	nr := records.NewNullRecord("Foo")
	nr.EXTHSection.AddString(templates.EXTHTitle, "BookTitle")
	nr.EXTHSection.AddInt(templates.EXTHAdult, 0)
	buf := bytes.NewBuffer(nil)
	nr.Write(buf)

	assertEq(t, nr.Length(), len(buf.Bytes()))
}

func TestIndexHeaderRecordLength(t *testing.T) {
	hr := records.SkeletonHeaderIndexRecord(0)
	buf := bytes.NewBuffer(nil)
	hr.Write(buf)

	assertEq(t, hr.Length(), len(buf.Bytes()))
	assertEq(t, hr.Length()%4, 0)
}

func TestSkeletonHeaderRecordLength(t *testing.T) {
	sr := records.SkeletonIndexRecord(nil)
	buf := bytes.NewBuffer(nil)
	sr.Write(buf)

	assertEq(t, sr.Length(), len(buf.Bytes()))
	assertEq(t, sr.Length()%4, 0)
}

func TestSkeletonHeaderRecordLengthWithChapter(t *testing.T) {
	chunk := records.ChunkInfo{
		PreStart:      0,
		PreLength:     100,
		ContentStart:  100,
		ContentLength: 100,
	}
	sr := records.SkeletonIndexRecord([]records.ChunkInfo{chunk})
	buf := bytes.NewBuffer(nil)
	sr.Write(buf)

	assertEq(t, sr.Length(), len(buf.Bytes()))
	assertEq(t, sr.Length()%4, 0)
}

func TestSkeletonHeaderRecordLengthWithChapters(t *testing.T) {
	chunk := records.ChunkInfo{
		PreStart:      0,
		PreLength:     100,
		ContentStart:  100,
		ContentLength: 100,
	}
	sr := records.SkeletonIndexRecord([]records.ChunkInfo{chunk, chunk, chunk})
	buf := bytes.NewBuffer(nil)
	sr.Write(buf)

	assertEq(t, sr.Length(), len(buf.Bytes()))
	assertEq(t, sr.Length()%4, 0)
}

func TestReadWrite(t *testing.T) {
	// Write
	w := bytes.NewBuffer(nil)
	db := pdb.NewDatabase("Test_Book")
	db.AddRecord(pdb.RawRecord{'o'})
	db.AddRecord(pdb.RawRecord{'h', 'i'})
	db.AddRecord(pdb.RawRecord{'c', 'a', 't'})
	db.AddRecord(pdb.RawRecord{'t', 'r', 'e', 'e'})
	db.Write(w)

	// Read
	r := bytes.NewReader(w.Bytes())
	db2, _ := pdb.ReadDatabase(r)

	// Compare
	assertEq(t, db.Name, db2.Name)
	assertEq(t, len(db.Records), len(db2.Records))
	for i := 0; i < len(db.Records); i++ {
		assertEq(t, db.Records[i].Length(), db.Records[i].Length())
	}
}

func assertEq(t *testing.T, v1 interface{}, v2 interface{}) {
	if v1 != v2 {
		t.Errorf("Not equal: %v, %v", v1, v2)
	}
}

func measure(v interface{}) int {
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, pdb.Endian, v)

	return len(buf.Bytes())
}
