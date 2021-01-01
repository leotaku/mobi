package records

import (
	"encoding/binary"
	"io"

	"github.com/leotaku/manki/mobi/pdb"
	t "github.com/leotaku/manki/mobi/templates"
)

type EXTHSection struct {
	entries []EXTHEntry
}

func NewEXTHSection() EXTHSection {
	return EXTHSection{
		entries: []EXTHEntry{},
	}
}

func (sec *EXTHSection) AddString(tp t.EXTHEntryType, s string) {
	data := []byte(s)
	for i, entry := range sec.entries {
		if entry.EntryType == tp {
			sec.entries[i].Data = data
			return
		}
	}
	entry := NewEXTHEntry(tp, data)
	sec.entries = append(sec.entries, entry)
}

func (sec *EXTHSection) AddInt(tp t.EXTHEntryType, i int) {
	data := make([]byte, 4)
	pdb.Endian.PutUint32(data, uint32(i))
	for i, entry := range sec.entries {
		if entry.EntryType == tp {
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
	EntryType t.EXTHEntryType
	Data      []byte
}

func NewEXTHEntry(tp t.EXTHEntryType, data []byte) EXTHEntry {
	return EXTHEntry{
		EntryType: tp,
		Data:      data,
	}
}

func (e EXTHEntry) Length() int {
	return len(e.Data) + t.EXTHEntryHeaderLength
}

func (e EXTHEntry) Write(w io.Writer) error {
	h := t.NewEXTHEntryHeader(e.EntryType, uint32(e.Length()))
	err := binary.Write(w, pdb.Endian, h)
	if err != nil {
		return err
	}

	_, err = w.Write(e.Data)
	return err
}
