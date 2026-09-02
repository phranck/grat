package detect

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/project"
)

const readsPort = `package config

import "os"

func Port() string { return os.Getenv("PORT") }
`

// TestTheGoScanStaysInsideTheBoundsEveryOtherScanHas is the fault this guards.
// The search walked the whole module with a traversal of its own, without the
// entry count and the depth every other scan carries, so grat discover over a
// folder holding a large unpacked module tree spent its time in there.
func TestTheGoScanStaysInsideTheBoundsEveryOtherScanHas(t *testing.T) {
	t.Parallel()

	for name, testCase := range map[string]struct {
		where string
		want  bool
	}{
		"a package of the module":   {where: "internal/config/config.go", want: true},
		"an unpacked dependency":    {where: "vendor/example.com/lib/lib.go", want: false},
		"input the toolchain skips": {where: "testdata/fixture/main.go", want: false},
		"deeper than a scan looks": {
			where: strings.Repeat("nested/", project.MaxScanDepth+1) + "config.go",
			want:  false,
		},
	} {
		root := writeProject(t, map[string]string{
			"go.mod":       "module example.com/service\n\ngo 1.25\n",
			testCase.where: readsPort,
		})
		if got := readsPortFromGoSource(root); got != testCase.want {
			t.Fatalf("%s: readsPortFromGoSource() = %v for %s, want %v",
				name, got, filepath.ToSlash(testCase.where), testCase.want)
		}
	}
}
