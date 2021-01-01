package pdb

import "strings"

func trimZeroes(s string) string {
	return strings.TrimRight(s, "\0000")
}

func underscoreSpaces(s string) string {
	return strings.ReplaceAll(s, " ", "_")
}
