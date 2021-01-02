package records

import "io"

const TextRecordMaxSize = 4096 // 0x1000

type TextRecord struct {
	data []byte
}

func NewTextRecord(s string) TextRecord {
	if len(s) > TextRecordMaxSize {
		panic("TBSTextRecord too large")
	}
	return TextRecord{data: []byte(s)}
}

func (r TextRecord) Write(w io.Writer) error {
	_, err := w.Write(append(r.data, encodeVwi(0)...))
	return err
}

func (r TextRecord) Length() int {
	return len(r.data) + 1
}
