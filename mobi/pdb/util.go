package pdb

import "strings"

func trimZeroes(s string) string {
	return strings.TrimRight(s, "\0000")
}
