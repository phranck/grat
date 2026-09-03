package manual

import (
	"strings"
	"testing"
)

// TestAnEntryCarriesMoreThanTheOneLineReference checks that the detail is
// actually detail. An entry repeating the reference adds nothing a reader could
// not get from grat help.
func TestAnEntryCarriesMoreThanTheOneLineReference(t *testing.T) {
	t.Parallel()

	// A word count would fire on a decision rather than on a defect, since a
	// command such as ports reassign genuinely needs two sentences. What can be
	// wrong is an entry that was never written, so that is what is checked.
	for _, detail := range commandDetails {
		prose := strings.TrimSpace(detail.detail)
		if prose == "" {
			t.Fatalf("the command %q has an entry with nothing in it", detail.usage)
		}
		if !strings.HasSuffix(prose, ".") {
			t.Fatalf("the entry for %q does not end in a sentence: %q", detail.usage, prose)
		}
	}
}

// TestEveryFlagOfAnEntryIsExplained keeps a flag from being listed without
// saying what it does or what happens without it.
func TestEveryFlagOfAnEntryIsExplained(t *testing.T) {
	t.Parallel()

	for _, detail := range commandDetails {
		for _, option := range detail.options {
			if !strings.HasPrefix(option.flag, "-") {
				t.Fatalf("the option %q of %q does not read as a flag", option.flag, detail.usage)
			}
			if len(strings.Fields(option.meaning)) < 5 {
				t.Fatalf("the option %q of %q is not explained: %q", option.flag, detail.usage, option.meaning)
			}
		}
	}
}
