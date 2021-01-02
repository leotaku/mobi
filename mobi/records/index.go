package records

import (
	"bytes"
	"fmt"

	"github.com/leotaku/manki/mobi/pdb"
	t "github.com/leotaku/manki/mobi/types"
)

func NCXHeaderIndexRecord(entryCount int) IndexRecord {
	bs := encodeINDXString(fmt.Sprintf("%03v", entryCount-1))
	pad := make([]byte, 5)
	pdb.Endian.PutUint16(pad, uint16(entryCount))
	bs = append(bs, pad...)

	return IndexRecord{
		TAGXTable:     t.TAGXTableNCXSingle,
		Type:          2,
		IDXTEntries:   [][]byte{bs},
		SubEntryCount: uint32(entryCount),
		CNCXCount:     1,
	}
}

func SkeletonHeaderIndexRecord(entryCount int) IndexRecord {
	bs := encodeINDXString(fmt.Sprintf("SKEL%010v", entryCount-1))
	pad := make([]byte, 5)
	pdb.Endian.PutUint16(pad, uint16(entryCount))
	bs = append(bs, pad...)

	return IndexRecord{
		TAGXTable:     t.TAGXTableSkeleton,
		Type:          2,
		IDXTEntries:   [][]byte{bs},
		SubEntryCount: uint32(entryCount),
	}
}

func ChunkHeaderIndexRecord(lastPos int, entryCount int) IndexRecord {
	bs := encodeINDXString(fmt.Sprintf("%010v", lastPos))
	pad := make([]byte, 5)
	pdb.Endian.PutUint16(pad, uint16(entryCount))
	bs = append(bs, pad...)

	return IndexRecord{
		TAGXTable:     t.TAGXTableChunk,
		Type:          2,
		IDXTEntries:   [][]byte{bs},
		SubEntryCount: uint32(entryCount),
		CNCXCount:     1,
	}
}

func NCXIndexRecord(info []ChapterInfo) (IndexRecord, CNCXRecord) {
	idxtEntries := make([][]byte, 0)
	cncxEntries := make([][]byte, 0)
	cncxOffset := 0
	for _, chap := range info {
		// CNCX entries
		s := fmt.Sprintf(chap.Title)
		cncx := encodeCNCXString(s)
		cncxEntries = append(cncxEntries, cncx)

		label := encodeINDXString(fmt.Sprintf("%03v", chap.Start))
		buf := bytes.NewBuffer(nil)
		writeSequential(buf, pdb.Endian,
			label,
			calculateControlByte(t.TAGXTableChunk),
			encodeVwi(chap.Start),  // Record offset
			encodeVwi(chap.Length), // Lenght of a record
			encodeVwi(cncxOffset),  // Label offset relative to CNXC record
			encodeVwi(0),           // Null
		)
		idxtEntries = append(idxtEntries, buf.Bytes())
		cncxOffset += len(cncx)
	}

	return IndexRecord{
			Type:          0,
			HeaderType:    1,
			IDXTEntries:   idxtEntries,
			SubEntryCount: 0,
		}, CNCXRecord{
			entries: cncxEntries,
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
			encodeVwi(chunk.PreLength),
			encodeVwi(chunk.PreStart),
			encodeVwi(chunk.PreLength),
		)
		entries = append(entries, buf.Bytes())
	}

	return IndexRecord{
		Type:        0,
		HeaderType:  1,
		IDXTEntries: entries,
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
			Type:          0,
			HeaderType:    1,
			IDXTEntries:   idxtEntries,
			SubEntryCount: 0,
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

type ChapterInfo struct {
	Title  string
	Start  int
	Length int
}

func encodeINDXString(s string) []byte {
	len := byte(len(s))
	return append([]byte{len}, s...)
}

func encodeCNCXString(label string) []byte {
	len := len(label)
	return append(encodeVwi(len), label...)
}
