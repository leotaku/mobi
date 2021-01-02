package types

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/leotaku/manki/mobi/pdb"
)

const FLISRecordLength = 36

type FLISRecord struct {
	FLIS    [4]byte
	Fixed1  uint32
	Fixed2  uint16
	Fixed3  uint16
	Fixed4  uint32
	Fixed5  uint32
	Fixed6  uint16
	Fixed7  uint16
	Fixed8  uint32
	Fixed9  uint32
	Fixed10 uint32
}

func NewFLISRecord() FLISRecord {
	return FLISRecord{
		FLIS:    [4]byte{'F', 'L', 'I', 'S'},
		Fixed1:  8,
		Fixed2:  65,
		Fixed3:  0,
		Fixed4:  0,
		Fixed5:  math.MaxUint32,
		Fixed6:  1,
		Fixed7:  3,
		Fixed8:  3,
		Fixed9:  1,
		Fixed10: math.MaxUint32,
	}
}

func (r FLISRecord) Write(w io.Writer) error {
	return binary.Write(w, pdb.Endian, r)
}

func (r FLISRecord) Length() int {
	return FLISRecordLength
}

const FCISRecordLength = 52

type FCISRecord struct {
	FCIS       [4]byte
	Fixed1     uint32
	Fixed2     uint32
	Fixed3     uint32
	Fixed4     uint32
	TextLength uint32
	Fixed5     uint32
	Fixed6     uint32
	Fixed7     uint32
	Fixed8     uint32
	Fixed9     uint32
	Fixed10    uint16
	Fixed11    uint16
	Fixed12    uint32
}

func NewFCISRecord(TextLength uint32) FCISRecord {
	return FCISRecord{
		FCIS:       [4]byte{'F', 'C', 'I', 'S'},
		Fixed1:     20,
		Fixed2:     16,
		Fixed3:     2,
		Fixed4:     0,
		TextLength: TextLength,
		Fixed5:     0,
		Fixed6:     40,
		Fixed7:     0,
		Fixed8:     40,
		Fixed9:     8,
		Fixed10:    1,
		Fixed11:    1,
		Fixed12:    0,
	}
}

func (r FCISRecord) Write(w io.Writer) error {
	return binary.Write(w, pdb.Endian, r)
}

func (r FCISRecord) Length() int {
	return FCISRecordLength
}

const EOFRecordLength = 4

var EOFRecord = pdb.RawRecord{0xE9, 0x8E, 0x0D, 0x0A}
