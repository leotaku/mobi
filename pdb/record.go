package pdb

import "io"

// Record represents a generic Palm database record.
type Record interface {
	Write(io.Writer) error
}

// RawRecord represents an uninterpreted Palm database record.
type RawRecord []byte

func (r RawRecord) Write(w io.Writer) error {
	_, err := w.Write(r)
	return err
}
