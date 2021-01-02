package records

import "io"

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
