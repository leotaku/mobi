package mobi

import (
	"bytes"
	"fmt"
	"io"

	"github.com/leotaku/manki/mobi/pdb"
	t "github.com/leotaku/manki/mobi/templates"
)

type PrimaryIndexRecord struct {
	totalIndexEntries int
}

func NewPrimaryIndexRecord(totalIndexEntries int) PrimaryIndexRecord {
	return PrimaryIndexRecord{
		totalIndexEntries: totalIndexEntries,
	}
}

func (r PrimaryIndexRecord) Write(w io.Writer) error {
	indx := t.NewINDXHeader(1, uint32(r.totalIndexEntries))
	tagx := t.NewTAGXSingleHeader()
	label := encodeLabelIdentifier(r.totalIndexEntries - 1)
	li := uint16(r.totalIndexEntries)
	pad := [2]byte{}
	idtx := t.NewIDXTHeader(t.INDXHeaderLength + t.TAGXSingleHeaderLength)
	pad2 := [2]byte{}

	return writeSequential(w, pdb.Endian, indx, tagx, label, li, pad, idtx, pad2)
}

func (r PrimaryIndexRecord) Length() int {
	return t.INDXHeaderLength + t.TAGXSingleHeaderLength + t.IDXTHeaderLength + 10
}

type SecondaryIndexRecord struct {
	indexEntries []chapterIndexEntry
}

type chapterIndexEntry struct {
	Offset int
	Length int
	Label  string
}

func NewSecondaryIndexRecord(offset int, chaps []Chapter) SecondaryIndexRecord {
	indexes := make([]chapterIndexEntry, 0)
	for _, chap := range chaps {
		len := len(chap.TextContent)
		fmt.Println("Chapter Offset:", offset)
		fmt.Println("Chapter Length:", len)
		indexes = append(indexes, chapterIndexEntry{
			Offset: offset,
			Length: len,
			Label:  chap.Name,
		})
		offset += len
	}

	return SecondaryIndexRecord{
		indexEntries: indexes,
	}
}

func (r SecondaryIndexRecord) Write(w io.Writer) error {
	buf, offsets := writeCNCXOffsets(r.indexEntries)
	h := t.NewINDXHeader(uint32(len(r.indexEntries)), uint32(len(r.indexEntries)))
	h.IDXTStart = uint32(t.INDXHeaderLength + len(buf))
	h.HeaderType = 1
	h.IndexType = 0
	idxt := [4]byte{'I', 'D', 'X', 'T'}

	// INDX header, CNCX offsets and IDXT
	err := writeSequential(w, pdb.Endian, h, buf, idxt)
	if err != nil {
		return err
	}

	// Write IDXT offsets
	for _, off := range offsets {
		err := writeSequential(w, pdb.Endian, uint16(off+t.INDXHeaderLength))
		if err != nil {
			return err
		}
	}

	// Write padding
	pad := [2]byte{}
	err = writeSequential(w, pdb.Endian, pad)
	if err != nil {
		return err
	}

	return nil
}

func writeCNCXOffsets(indexes []chapterIndexEntry) ([]byte, []uint16) {
	buf := bytes.NewBuffer(nil)
	idxtOffsets := make([]uint16, 0) // TODO
	labelOffset := 0
	for i, idx := range indexes {
		id := encodeLabelIdentifier(i)
		label := encodeLabel(idx.Label)
		idxtOffsets = append(idxtOffsets, uint16(buf.Len()))
		writeSequential(buf, pdb.Endian,
			id,                     // Len of ID and ID
			t.CBSingle,             // Control Byte
			encodeVwi(idx.Offset),  // Record offset
			encodeVwi(idx.Length),  // Lenght of a record
			encodeVwi(labelOffset), // Label offset relative to CNXC record
			encodeVwi(0),           // Null
		)
		labelOffset += len(label)
	}

	return buf.Bytes(), idxtOffsets
}

func (r SecondaryIndexRecord) Length() int {
	buf := bytes.NewBuffer(nil)
	r.Write(buf)
	return len(buf.Bytes())
}

type CNCXRecord struct {
	labels []string
}

func NewCNCXRecord(labels []string) CNCXRecord {
	return CNCXRecord{labels: labels}
}

func (r CNCXRecord) Write(w io.Writer) error {
	for _, label := range r.labels {
		_, err := w.Write(encodeLabel(label))
		if err != nil {
			return err
		}
	}

	return nil
}

func (r CNCXRecord) Length() int {
	result := 0
	for _, label := range r.labels {
		ll := len(label)
		result += ll + len(encodeVwi(ll))
	}

	return result
}

func createIndexRecords(offset int, chaps []Chapter) (PrimaryIndexRecord, SecondaryIndexRecord, CNCXRecord) {
	labels := make([]string, 0)
	for _, chap := range chaps {
		labels = append(labels, chap.Name)
	}

	first := NewPrimaryIndexRecord(len(labels))
	second := NewSecondaryIndexRecord(offset, chaps)
	cncx := NewCNCXRecord(labels)

	return first, second, cncx
}

func encodeLabelIdentifier(indx int) []byte {
	id := fmt.Sprintf("%03v", indx)
	len := byte(len(id))
	return append([]byte{len}, id...)
}

func encodeLabel(label string) []byte {
	len := len(label)
	return append(encodeVwi(len), label...)
}
