package pdb

import "io"

// Record represents a generic Palm database record.
type Record interface {
	Write(io.Writer) error
	Length() int
}

// RawRecord represents an uninterpreted Palm database record.
type RawRecord []byte

func (r RawRecord) Write(w io.Writer) error {
	_, err := w.Write(r)
	return err
}

func (r RawRecord) Length() int {
	return len(r)
}

type Info struct {
	mapping map[string]int
}

func (i *Info) Set(s string, v int) {
	i.mapping[s] = int(v)
}
