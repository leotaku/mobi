package records

import (
	"encoding/binary"
	"io"
)

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
