// Package textsafe identifies and neutralizes Unicode characters that can
// change terminal layout or conceal the visual order of untrusted text.
package textsafe

import (
	"strings"
	"unicode"
)

// UnsafeRune reports control and format characters, including bidi controls,
// isolates, zero-width format characters, and byte-order marks.
func UnsafeRune(character rune) bool {
	return unicode.IsControl(character) || unicode.In(character, unicode.Cf)
}

// ContainsUnsafe reports whether value contains a control or format character.
func ContainsUnsafe(value string) bool {
	for _, character := range value {
		if UnsafeRune(character) {
			return true
		}
	}
	return false
}

// Sanitize replaces unsafe characters with a visible replacement marker.
func Sanitize(value string) string {
	return strings.Map(func(character rune) rune {
		if UnsafeRune(character) {
			return '\uFFFD'
		}
		return character
	}, value)
}
