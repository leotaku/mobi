package records

import (
	"bytes"
	"fmt"
	"io"

	"github.com/leotaku/manki/mobi/pdb"
	t "github.com/leotaku/manki/mobi/templates"
)

type IndexRecord struct {
	Type           uint32
	HeaderType     uint32
	RecordCount    uint32
	TAGXTable      t.TAGXTagTable
	IDXTEntries    [][]byte
	IDXTEntryCount uint32
	CNCXCount      uint32
}

func SkeletonHeaderIndexRecord(entryCount int) IndexRecord {
	bs := encodeINDXString(fmt.Sprintf("SKEL%010v", entryCount-1))
	pad := make([]byte, 5)
	pdb.Endian.PutUint16(pad, uint16(entryCount))
	bs = append(bs, pad...)

	return IndexRecord{
		TAGXTable:      t.TAGXTableSkeleton,
		Type:           2, // Inflection
		RecordCount:    1,
		IDXTEntries:    [][]byte{bs},
		IDXTEntryCount: uint32(entryCount),
	}
}

func SkeletonIndexRecord(info []ChunkInfo) IndexRecord {
	entries := make([][]byte, 0)
	for i, chunk := range info {
		buf := bytes.NewBuffer(nil)
		label := encodeINDXString(fmt.Sprintf("SKEL%010v", i))

		writeSequential(buf, pdb.Endian,
			label,
			calculateControlByte(t.TAGXTableSkeleton),
			encodeVwi(1),
			encodeVwi(1),
			encodeVwi(chunk.PreStart),
			encodeVwi(chunk.PreLength-18),
			encodeVwi(chunk.PreStart),
			encodeVwi(chunk.PreLength-18),
		)
		entries = append(entries, buf.Bytes())
	}

	return IndexRecord{
		Type:        0,
		HeaderType:  1,
		RecordCount: uint32(len(info)),
		IDXTEntries: entries,
	}
}

func ChunkHeaderIndexRecord(lastPos int, entryCount int) IndexRecord {
	bs := encodeINDXString(fmt.Sprintf("%010v", lastPos))
	pad := make([]byte, 5)
	pdb.Endian.PutUint16(pad, uint16(entryCount))
	bs = append(bs, pad...)

	return IndexRecord{
		TAGXTable:      t.TAGXTableChunk,
		Type:           2,
		RecordCount:    1,
		IDXTEntries:    [][]byte{bs},
		IDXTEntryCount: uint32(entryCount),
		CNCXCount:      1,
	}
}

func ChunkIndexRecord(info []ChunkInfo) (IndexRecord, CNCXRecord) {
	idxtEntries := make([][]byte, 0)
	cncxEntries := make([][]byte, 0)
	cncxOffset := 0
	for i, chunk := range info {
		// CNCX entries
		s := fmt.Sprintf("P-//*[@aid='%04v']", i)
		cncx := encodeCNCXString(s)
		cncxEntries = append(cncxEntries, cncx)

		label := encodeINDXString(fmt.Sprintf("%010v", chunk.ContentStart))
		buf := bytes.NewBuffer(nil)
		writeSequential(buf, pdb.Endian,
			label,
			calculateControlByte(t.TAGXTableChunk),
			encodeVwi(cncxOffset),          // CNCX offset
			encodeVwi(i),                   // File number
			encodeVwi(i),                   // Sequence number
			encodeVwi(0),                   // Geometry start
			encodeVwi(chunk.ContentLength), // Geometry length
		)
		idxtEntries = append(idxtEntries, buf.Bytes())
		cncxOffset += len(cncx)
	}

	return IndexRecord{
			Type:           0,
			HeaderType:     1,
			RecordCount:    uint32(len(info)),
			IDXTEntries:    idxtEntries,
			IDXTEntryCount: 0,
		}, CNCXRecord{
			entries: cncxEntries,
		}
}

type ChunkInfo struct {
	PreStart      int
	PreLength     int
	ContentStart  int
	ContentLength int
}

type CNCXRecord struct {
	entries [][]byte
}

func (r CNCXRecord) Write(w io.Writer) error {
	for _, entry := range r.entries {
		_, err := w.Write(entry)
		if err != nil {
			return err
		}
	}

	pad := make([]byte, r.LengthNoPadding()%4)
	_, err := w.Write(pad)
	if err != nil {
		return err
	}

	return nil
}

func (r CNCXRecord) Length() int {
	length := r.LengthNoPadding()
	return length + length%4
}

func (r CNCXRecord) LengthNoPadding() int {
	result := 0
	for _, entry := range r.entries {
		result += len(entry)
	}
	return result
}

var bitmaskToShiftMap = map[uint8]uint8{1: 0, 2: 1, 3: 0, 4: 2, 8: 3, 12: 2, 16: 4, 32: 5, 48: 4, 64: 6, 128: 7, 192: 6}

func calculateControlByte(tagx t.TAGXTagTable) byte {
	cbs := make([]byte, 0)
	ans := uint8(0)
	for _, tag := range tagx {
		_, tagnum, bm, cb := deconstructTag(tag)
		if cb == 1 {
			cbs = append(cbs, ans)
			ans = 0
			continue
		}
		nvals := mapTagToNvals(tag)
		nentries := nvals / tagnum
		shifts := bitmaskToShiftMap[bm]
		ans |= bm & (nentries << shifts)
	}

	return cbs[0]
}

func mapTagToNvals(tag t.TAGXTag) byte {
	if tag == t.TAGXTagSkeletonGeometry {
		return 4
	} else if tag == t.TAGXTagChunkGeometry || tag == t.TAGXTagSkeletonChunkCount {
		return 2
	} else {
		return 1
	}
}

func deconstructTag(tag t.TAGXTag) (byte, byte, byte, byte) {
	bs := make([]byte, 4)
	pdb.Endian.PutUint32(bs, uint32(tag))

	return bs[0], bs[1], bs[2], bs[3]
}

func encodeINDXString(s string) []byte {
	len := byte(len(s))
	return append([]byte{len}, s...)
}

func encodeCNCXString(label string) []byte {
	len := len(label)
	return append(encodeVwi(len), label...)
}

func (r IndexRecord) Write(w io.Writer) error {
	// Headers
	inh := t.NewINDXHeader(0, 0)
	th := t.NewTAGXHeader()
	idh := t.NewIDXTHeader()
	offset := t.INDXHeaderLength

	// INDX variables
	inh.IndexRecordCount = r.RecordCount
	inh.IndexType = r.Type
	inh.HeaderType = r.HeaderType
	inh.CNCXCount = r.CNCXCount

	// TAGX variables
	if len(r.TAGXTable) > 0 {
		inh.TAGXOffset = t.INDXHeaderLength
		th.HeaderLength += uint32(len(r.TAGXTable) * t.TAGXTagLength)
		offset += int(th.HeaderLength)
		fmt.Printf("TAGX offset: %#x\n", t.INDXHeaderLength)
		fmt.Printf("TAGX size: %#x\n", uint32(len(r.TAGXTable)*t.TAGXTagLength))
		fmt.Printf("TAGX end: %#x\n", offset)
	} else {
		inh.TAGXOffset = 0
	}

	// IDXT variables
	idxtLength := 0
	idxtOffsets := make([]uint16, 0)
	for _, entry := range r.IDXTEntries {
		len := len(entry)
		idxtOffsets = append(idxtOffsets, uint16(offset))
		offset += len
		idxtLength += len
	}
	inh.IDXTStart = uint32(offset + idxtLength%4)
	inh.IndexEntryCount = r.IDXTEntryCount
	fmt.Printf("INDX offset: %#x\n", offset)
	fmt.Printf("INDX size: %#x\n", 4+len(r.IDXTEntries)*2)
	fmt.Printf("INDX end: %#x\n", offset+4+len(r.IDXTEntries)*2)

	// Write INDX header
	err := writeSequential(w, pdb.Endian, inh)
	if err != nil {
		return err
	}

	if len(r.TAGXTable) > 0 {
		// Write TAGX section
		err := writeSequential(w, pdb.Endian, th, r.TAGXTable)
		if err != nil {
			return err
		}

	}

	// Write entries?
	for _, entry := range r.IDXTEntries {
		err := writeSequential(w, pdb.Endian, entry)
		if err != nil {
			return err
		}
	}

	// Write IDXT and padding
	pad1 := make([]byte, idxtLength%4)
	pad2 := make([]byte, r.LengthNoPadding()%4)
	return writeSequential(w, pdb.Endian, pad1, idh, idxtOffsets, pad2)
}

func (r IndexRecord) Length() int {
	length := r.LengthNoPadding()
	return length + length%4
}

func (r IndexRecord) LengthNoPadding() int {
	length := t.INDXHeaderLength + t.IDXTHeaderLength + len(r.IDXTEntries)*2
	entriesLength := 0
	for _, entry := range r.IDXTEntries {
		length += len(entry)
		entriesLength += len(entry)
	}
	length += entriesLength % 4
	if len(r.TAGXTable) > 0 {
		length += t.TAGXHeaderLength + len(r.TAGXTable)*t.TAGXTagLength
	}

	return length
}
