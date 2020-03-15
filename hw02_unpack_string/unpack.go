package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"io"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var builder strings.Builder
	reader := strings.NewReader(s)
	prevRune, isPrevDigit := '0', true

	for r, _, err := reader.ReadRune(); err != io.EOF; r, _, err = reader.ReadRune() {
		switch {
		case r == '\\':
			//Read,write escaped Rune
			if r, _, err = reader.ReadRune(); err != io.EOF {
				builder.WriteRune(r)
				prevRune, isPrevDigit = r, false
			}
		case unicode.IsDigit(r):
			//Two consecutive not escaped digits is error
			if isPrevDigit {
				return "", ErrInvalidString
			}

			//Repeat previous rune N-1 times
			n := int(r - '0')
			if n > 1 { //nolint:gomnd
				builder.WriteString(strings.Repeat(string(prevRune), n-1)) //nolint:gomnd
			}
			prevRune, isPrevDigit = r, true
		default:
			builder.WriteRune(r)
			prevRune, isPrevDigit = r, false
		}
	}
	return builder.String(), nil
}
