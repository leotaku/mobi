package records

import (
	"io"

	"github.com/leotaku/manki/mobi/pdb"
	t "github.com/leotaku/manki/mobi/types"
)

type FDSTRecord struct {
	entries []t.FDSTEntry
}

func NewFDSTRecord(flows ...string) FDSTRecord {
	entries := make([]t.FDSTEntry, 0)
	offset := 0
	for _, s := range flows {
		entries = append(entries, t.FDSTEntry{
			Start: uint32(offset),
			End:   uint32(offset + len(s)),
		})
		offset += len(s)
	}

	return FDSTRecord{
		entries: entries,
	}
}

func (r FDSTRecord) Write(w io.Writer) error {
	h := t.NewFDSTHeader()
	h.EntryCount = uint32(len(r.entries))

	return writeSequential(w, pdb.Endian, h, r.entries)
}

func (r FDSTRecord) Length() int {
	return t.FDSTHeaderLength + len(r.entries)*t.FDSTEntryLength
}
