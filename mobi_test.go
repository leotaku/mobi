package mobi_test

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"

	"github.com/leotaku/mobi/pdb"
	"github.com/leotaku/mobi/records"
	"github.com/leotaku/mobi/types"
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
	bs := writeRecord(nr)

	assertEq(t, nr.Length(), len(bs))
}

func TestIndexSectionLength(t *testing.T) {
	indx := types.NewINDXHeader(0, 0)
	tagx := types.NewTAGXSingleHeader()
	idtx := types.NewIDXTSingleHeader(0)

	assertEq(t, measure(indx), types.INDXHeaderLength)
	assertEq(t, measure(tagx), types.TAGXSingleHeaderLength)
	assertEq(t, measure(idtx), types.IDXTSingleHeaderLength)
}

func TestNullRecordLengthWithEXTH(t *testing.T) {
	nr := records.NewNullRecord("Foo")
	nr.EXTHSection.AddString(types.EXTHTitle, "BookTitle")
	nr.EXTHSection.AddInt(types.EXTHAdult, 0)
	bs := writeRecord(nr)

	assertEq(t, nr.Length(), len(bs))
}

func TestIndexHeaderRecordLength(t *testing.T) {
	hr := records.SkeletonHeaderIndexRecord(0)
	bs := writeRecord(hr)

	assertEq(t, hr.Length(), len(bs))
	assertEq(t, hr.Length()%4, 0)
}

func TestSkeletonHeaderRecordLength(t *testing.T) {
	sr := records.SkeletonIndexRecord(nil)
	bs := writeRecord(sr)

	assertEq(t, sr.Length(), len(bs))
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
	bs := writeRecord(sr)

	assertEq(t, sr.Length(), len(bs))
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
	bs := writeRecord(sr)

	assertEq(t, sr.Length(), len(bs))
	assertEq(t, sr.Length()%4, 0)
}

func TestReadWrite(t *testing.T) {
	// Write
	w := bytes.NewBuffer(nil)
	db := pdb.NewDatabase("Test Book", time.Unix(0, 0))
	db.AddRecord(pdb.RawRecord{'o'})
	db.AddRecord(pdb.RawRecord{'h', 'i'})
	db.AddRecord(pdb.RawRecord{'c', 'a', 't'})
	db.AddRecord(pdb.RawRecord{'t', 'r', 'e', 'e'})
	err := db.Write(w)
	if err != nil {
		t.Fatal(err)
	}

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

func writeRecord(r pdb.Record) []byte {
	buf := bytes.NewBuffer(nil)
	err := r.Write(buf)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func measure(v interface{}) int {
	buf := bytes.NewBuffer(nil)
	err := binary.Write(buf, pdb.Endian, v)
	if err != nil {
		panic(err)
	}

	return len(buf.Bytes())
}
