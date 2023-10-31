package pdb

import (
	"strings"
	"time"
)

func trimZeroes(s string) string {
	return strings.TrimRight(s, "\x00")
}

func convertToPalmTime(t time.Time) uint32 {
	delta := t.Sub(time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC))
	return uint32(delta.Seconds())
}

func convertFromPalmTime(t uint32) time.Time {
	start := time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC)
	return start.Add(time.Duration(t) * time.Second)
}
