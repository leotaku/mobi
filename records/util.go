package records

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/leotaku/mobi/pdb"
	t "github.com/leotaku/mobi/types"
)

func To32(i int) string {
	s := strconv.FormatInt(int64(i), 32)
	return fmt.Sprintf("%04v", strings.ToUpper(s))
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

func encodeTrailingBytes(data []byte) []byte {
	return append(data, encodeVwi(len(data))...)
}

func reverseBytes(buf []byte) {
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
}

func invMod(a, n int) int {
	return n - maxMod(a, n)
}

func maxMod(a, n int) int {
	switch mod := a % n; mod {
	case 0:
		return n
	default:
		return mod
	}
}

var bitmaskToShiftMap = map[uint8]uint8{1: 0, 2: 1, 3: 0, 4: 2, 8: 3, 12: 2, 16: 4, 32: 5, 48: 4, 64: 6, 128: 7, 192: 6}

func calculateControlByte(tagx t.TAGXTagTable) byte {
	cbs := make([]byte, 0)
	ans := uint8(0)
	for _, tag := range tagx {
		_, tagnum, bm, cb := deconstructTag(tag)
		if cb == 1 {
			cbs = append(cbs, ans)
			ans = 0
			continue
		}
		nvals := mapTagToNvals(tag)
		nentries := nvals / tagnum
		shifts := bitmaskToShiftMap[bm]
		ans |= bm & (nentries << shifts)
	}

	return cbs[0]
}

func mapTagToNvals(tag t.TAGXTag) byte {
	switch tag {
	case t.TAGXTagSkeletonGeometry:
		return 4
	case t.TAGXTagChunkGeometry, t.TAGXTagSkeletonChunkCount:
		return 2
	default:
		return 1
	}
}

func deconstructTag(tag t.TAGXTag) (byte, byte, byte, byte) {
	bs := make([]byte, 4)
	pdb.Endian.PutUint32(bs, uint32(tag))

	return bs[0], bs[1], bs[2], bs[3]
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

func bytesSequential(bo binary.ByteOrder, vs ...interface{}) []byte {
	buf := bytes.NewBuffer(nil)
	err := writeSequential(buf, bo, vs...)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}
