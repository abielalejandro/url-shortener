package utils

import (
	"math/big"
	"strings"
)

func PaddingRight(str string, char string, size int) string {
	return strings.Join([]string{str, strings.Repeat(char, size-len(str))}, "")
}

func PaddingLeft(str string, char string, size int) string {
	return strings.Join([]string{strings.Repeat(char, size-len(str)), str}, "")
}

func ToBase62(str string) string {
	var i big.Int
	i.SetBytes([]byte(str))
	return i.Text(62)
}
