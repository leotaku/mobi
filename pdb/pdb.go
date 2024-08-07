// Package pdb implements reading and writing PalmDB databases.
package pdb

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"
)

// Endian describes the byte-order of integers in Palm databases.
var Endian = binary.BigEndian

// Database represents an in-memory Palm database.
type Database struct {
	Name    string
	Date    time.Time
	Records []Record
}

// NewDatabase creates an empty Palm database with name and date.
func NewDatabase(name string, date time.Time) Database {
	return Database{
		Name:    trimZeroes(name),
		Date:    date,
		Records: []Record{},
	}
}

// AddRecord adds a generic record to the Palm database.
//
// Returns the index of the inserted record.
func (d *Database) AddRecord(r Record) int {
	d.Records = append(d.Records, r)
	return len(d.Records) - 1
}

// Idx returns the index of the last inserted record.
func (d *Database) Idx() int {
	return len(d.Records) - 1
}

// ReplaceRecord overrides the record at index i in the Palm database.
//
// Panics if index i is out of range.
func (d *Database) ReplaceRecord(i int, r Record) {
	d.Records[i] = r
}

// Write writes out the binary representation of the Palm database to w.
func (d Database) Write(w io.Writer) error {
	rnum := len(d.Records)
	palmDBHeader := NewPalmDBHeader(d.Name, d.Date, uint16(rnum), uint32(rnum)*2-1)
	err := binary.Write(w, Endian, palmDBHeader)
	if err != nil {
		return err
	}

	// Offsets
	buf := bytes.NewBuffer(nil)
	initialOffset := PalmDBHeaderLength + RecordHeaderLength*len(d.Records) + 2
	offsets := make([]int, 0)

	// Write records
	for _, rec := range d.Records {
		offsets = append(offsets, initialOffset+buf.Len())
		err := rec.Write(buf)
		if err != nil {
			return err
		}
	}

	// Write record headers
	for i, offset := range offsets {
		h := RecordHeader{
			Offset:    uint32(offset),
			Attribute: 0,             // No idea
			Skip:      0,             // No idea
			UniqueID:  uint16(i * 2), // Calibre doubles UID for some reason
		}
		err := binary.Write(w, Endian, h)
		if err != nil {
			return err
		}
	}

	// Write 2-byte padding
	pad := make([]byte, 2)
	_, err = w.Write(pad)
	if err != nil {
		return err
	}

	// Write buffer
	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}

	// Success
	return nil
}

// ReadDatabase reads an uninterpreted Palm database from r.
//
// This is not a lossless routine, in particular a large portion of
// metadata stored in the Palm database header will be ignored.
func ReadDatabase(r io.Reader) (*Database, error) {
	data, err := io.ReadAll(r)
	b := bytes.NewReader(data)
	if err != nil {
		return nil, err
	}

	palmDBHeader := PalmDBHeader{}
	err = binary.Read(b, Endian, &palmDBHeader)
	if err != nil {
		return nil, err
	}

	offsets := make([]RecordHeader, palmDBHeader.NumRecords)
	err = binary.Read(b, Endian, &offsets)
	if err != nil {
		return nil, err
	}

	name := trimZeroes(string(palmDBHeader.Name[:]))
	date := convertFromPalmTime(palmDBHeader.CreationTime)

	records := make([]Record, 0)
	for i := 1; i < len(offsets); i++ {
		curr := offsets[i].Offset
		prev := offsets[i-1].Offset
		records = append(records, RawRecord(data[prev:curr]))
	}
	last := offsets[len(offsets)-1].Offset
	records = append(records, RawRecord(data[last:]))

	return &Database{
		Name:    name,
		Date:    date,
		Records: records,
	}, nil
}
