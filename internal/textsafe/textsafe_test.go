package textsafe

import (
	"strings"
	"testing"
)

func TestUnsafeRuneRejectsControlsAndUnicodeFormatCharacters(t *testing.T) {
	t.Parallel()

	for _, character := range []rune{'\n', '\u009b', '\u200b', '\u202e', '\u2066', '\ufeff'} {
		if !UnsafeRune(character) {
			t.Errorf("UnsafeRune(%U) = false, want true", character)
		}
	}
	for _, character := range []rune{'a', 'ä', '界', '🙂', '\u0301'} {
		if UnsafeRune(character) {
			t.Errorf("UnsafeRune(%U) = true, want false", character)
		}
	}
}

func TestSanitizeReplacesEveryUnsafeRune(t *testing.T) {
	t.Parallel()

	got := Sanitize("safe\u202eevil\nline")
	if strings.ContainsAny(got, "\u202e\n") || got != "safe�evil�line" {
		t.Fatalf("Sanitize() = %q, want unsafe runes replaced", got)
	}
}
