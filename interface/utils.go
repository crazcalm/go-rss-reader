package gui

import (
	"strings"
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
		str = str[0:length]
		count = length
	}
	return leftPad(str, " ", length-count), nil
}

func rightPadExactLength(str, pad string, length int) (string, error) {
	count := len(str)
	if count > length {
		str = str[0:length]
		count = length
	}
	return rightPad(str, " ", length-count), nil
}
