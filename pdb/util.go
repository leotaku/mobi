package pdb

import (
	"strings"
	"time"
)

func trimZeroes(s string) string {
	return strings.TrimRight(s, "\x00")
}

func calculatePalmTime(t time.Time) uint32 {
	delta := t.Sub(time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC))
	return uint32(delta.Seconds())
}
