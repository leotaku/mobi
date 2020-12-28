package mobi

import (
	"encoding/binary"
	"io"
)

const preChapOffset = 25

func chaptersToText(c []Chapter) string {
	text := "<html><head></head><body>" // TODO
	for _, chap := range c {
		text = text + chap.TextContent
	}
	text = text + "</body></html>"

	return text
}

func genTextRecords(html string) []TBSTextRecord {
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	records := []TBSTextRecord{}
	recordCount := len(html) / TextRecordMaxSize
	if len(html)%TextRecordMaxSize != 0 {
		recordCount += 1
	}

	for i := 0; i < int(recordCount); i++ {
		from := i * TextRecordMaxSize
		to := from + TextRecordMaxSize

		record := NewTBSTextRecord(html[from:min(to, len(html))])
		records = append(records, record)
	}

	return records
}

func encodeVwi(x int) []byte {
	buf := make([]byte, 64)
	z := 0
	for {
		buf[z] = byte(x) & 0x7f
		x >>= 7
		z++
		if x == 0 {
			buf[0] |= 0x80
			break
		}
	}

	relevant := buf[:z]
	reverseBytes(relevant)
	return relevant
}

func reverseBytes(buf []byte) {
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
}

func invMod(dividend int, divisor int) int {
	return (divisor/2 + dividend) % divisor
}

func writeSequential(w io.Writer, bo binary.ByteOrder, vs ...interface{}) error {
	for _, v := range vs {
		err := binary.Write(w, bo, v)
		if err != nil {
			return err
		}
	}
	return nil
}
