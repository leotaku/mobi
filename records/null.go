package records

import (
	"encoding/binary"
	"io"

	"github.com/leotaku/mobi/pdb"
	t "github.com/leotaku/mobi/types"
)

const NullPaddingLength = 8192 // 0x2000

type NullRecord struct {
	PalmDocHeader t.PalmDocHeader
	MOBIHeader    t.KF8Header
	FullName      string
	EXTHSection   EXTHSection
}

func NewNullRecord(name string) NullRecord {
	return NullRecord{
		PalmDocHeader: t.NewPalmDocHeader(),
		MOBIHeader:    t.NewKF8Header(),
		FullName:      name,
		EXTHSection:   NewEXTHSection(),
	}
}

func (n NullRecord) Length() int {
	return t.PalmDocHeaderLength + t.KF8HeaderLength + n.EXTHSection.Length() + len(n.FullName) + NullPaddingLength
}

func (n NullRecord) Write(w io.Writer) error {
	// Set full name offset and length
	n.MOBIHeader.FullNameOffset = uint32(t.PalmDocHeaderLength + t.KF8HeaderLength + n.EXTHSection.Length())
	n.MOBIHeader.FullNameLength = uint32(len(n.FullName))

	// Write PalmDoc header
	err := binary.Write(w, pdb.Endian, n.PalmDocHeader)
	if err != nil {
		return err
	}

	// Write MOBI header
	err = binary.Write(w, pdb.Endian, n.MOBIHeader)
	if err != nil {
		return err
	}

	// Write EXTH header
	err = n.EXTHSection.Write(w)
	if err != nil {
		return err
	}

	// Write full name
	_, err = w.Write([]byte(n.FullName))
	if err != nil {
		return err
	}

	// Write padding
	pad := make([]byte, NullPaddingLength)
	_, err = w.Write(pad)
	return err
}
