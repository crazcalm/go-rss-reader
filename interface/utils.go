package gui

import (
	"strings"
	"errors"
)

func leftPad(str, pad string, times int) string {
	return strings.Repeat(pad, times) + str
}

func rightPad(str, pad string, times int) string {
	return str + strings.Repeat(pad, times)
}

func leftPadExactLength(str, pad string, length int) (string, error) {
	count := len(str)
	if count > length {
		return str, errors.New("Original string is longer than specified length")
	}
	return leftPad(str, " ", length - count), nil
}

func rightPadExactLength(str, pad string, length int) (string, error) {
	count := len(str)
	if count > length {
		return str, errors.New("Original string is longer than specified length")
	}
	return rightPad(str, " ", length - count), nil
}