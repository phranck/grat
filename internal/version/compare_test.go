package version

import "testing"

// TestCompareOrdersReleases pins the ordering grat update turns on. Installing
// whatever GitHub marks as the latest, without asking whether it is newer, is
// what this comparison exists to stop.
func TestCompareOrdersReleases(t *testing.T) {
	t.Parallel()

	for name, testCase := range map[string]struct {
		left  string
		right string
		want  int
	}{
		"older patch":            {left: "v1.9.0", right: "v1.9.1", want: -1},
		"older minor":            {left: "v1.8.9", right: "v1.9.0", want: -1},
		"older major":            {left: "v1.99.99", right: "v2.0.0", want: -1},
		"the same release":       {left: "v1.9.1", right: "v1.9.1", want: 0},
		"the same without a v":   {left: "1.9.1", right: "v1.9.1", want: 0},
		"newer patch":            {left: "v1.9.2", right: "v1.9.1", want: 1},
		"a pre-release is older": {left: "v2.0.0-rc.1", right: "v2.0.0", want: -1},
		"an unreadable tag":      {left: "latest", right: "v1.9.1", want: -1},
	} {
		if got := Compare(testCase.left, testCase.right); got != testCase.want {
			t.Fatalf("%s: Compare(%q, %q) = %d, want %d", name, testCase.left, testCase.right, got, testCase.want)
		}
	}
}
