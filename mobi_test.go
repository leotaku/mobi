package mobi_test

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"

	"github.com/leotaku/mobi/pdb"
)

func TestPDBHeaderLength(t *testing.T) {
	len := measure(pdb.PalmDBHeader{})
	assertEq(t, len, pdb.PalmDBHeaderLength)
}

func TestRecordHeaderLength(t *testing.T) {
	len := measure(pdb.RecordHeader{})
	assertEq(t, len, pdb.RecordHeaderLength)
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
	rdb, _ := pdb.ReadDatabase(r)

	// Compare
	assertEq(t, db.Name, rdb.Name)
	assertEq(t, len(db.Records), len(rdb.Records))
	for i := 0; i < len(db.Records); i++ {
		buf := bytes.NewBuffer(nil)
		rbuf := bytes.NewBuffer(nil)
		db.Records[i].Write(buf)
		rdb.Records[i].Write(rbuf)
		assertEq(t, buf.String(), rbuf.String())
	}
}

func assertEq(t *testing.T, v1 interface{}, v2 interface{}) {
	if v1 != v2 {
		t.Errorf("Not equal: %v, %v", v1, v2)
	}
}

func measure(v interface{}) int {
	buf := bytes.NewBuffer(nil)
	err := binary.Write(buf, pdb.Endian, v)
	if err != nil {
		panic(err)
	}

	return len(buf.Bytes())
}
