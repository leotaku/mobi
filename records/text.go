package records

import "io"

const TextRecordMaxSize = 4096 // 0x1000

type TextRecord struct {
	data  []byte
	trail []byte
}

func NewTextRecord(s string, trail TrailingData) TextRecord {
	if len(s) > TextRecordMaxSize {
		panic("TextRecord too large")
	}
	return TextRecord{
		data:  []byte(s),
		trail: trail.Encode(),
	}
}

func (r TextRecord) Write(w io.Writer) error {
	_, err := w.Write(r.data)
	if err != nil {
		return err
	}

	_, err = w.Write(r.trail)
	return err
}

func (r TextRecord) Length() int {
	return len(r.data) + len(r.trail)
}
