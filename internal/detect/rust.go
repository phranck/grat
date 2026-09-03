package detect

import (
	"regexp"
	"strings"
)

// rustServerCrates are the two web frameworks this recognises. Neither has a
// command line of its own, so neither can define a port flag; both official
// getting-started examples fix the address in the source instead.
var rustServerCrates = []string{"axum", "actix-web"}

// rustPortRead is the read that lets a crate take the port grat assigns.
var rustPortRead = regexp.MustCompile(`env::var(_os)?\s*\(\s*"PORT"`)

// rustLineComment and rustBlockComment are removed before the read above is
// looked for.
//
// Go's detector parses the syntax, which is the right way and is available
// there because the parser ships with the language. No Rust parser is at hand
// here, so this is a text search with the comments taken out first. That covers
// the case that actually happens, which is a mention in a comment, and it would
// still be fooled by the call appearing inside a string literal. Anything more
// would be a Rust parser, and the cost of the remaining gap is one project
// offered a command it turns out not to honour, which its first start reports.
var (
	rustLineComment  = regexp.MustCompile(`(?m)//.*$`)
	rustBlockComment = regexp.MustCompile(`(?s)/\*.*?\*/`)
)

// rustBinaryTarget matches a [[bin]] table, and rustDefaultRun the field that
// picks one where a crate has several.
var (
	rustBinaryTarget = regexp.MustCompile(`(?m)^\s*\[\[bin\]\]`)
	rustDefaultRun   = regexp.MustCompile(`(?m)^\s*default-run\s*=`)
)

// detectRust recognises a Rust web service and proposes a command only where
// cargo run would start one thing and that thing takes the port.
//
// Two questions decide it, and both have to be answered from files. Which
// binary does cargo run start, since a crate can have several and then needs a
// flag or a default-run to choose. And does the source read the port, since
// neither framework offers a way to pass one and every example fixes it.
func detectRust(root string) ([]Service, []Unresolved) {
	data, ok := readBounded(join(root, "Cargo.toml"))
	if !ok {
		return nil, nil
	}
	manifest := string(data)
	if !dependsOnRustServer(manifest) {
		return nil, nil
	}

	if reason := ambiguousRustBinary(root, manifest); reason != "" {
		return nil, []Unresolved{{Marker: "Cargo.toml", Reason: reason}}
	}
	if !readsPortFromRustSource(root) {
		return nil, []Unresolved{{
			Marker: "Cargo.toml",
			Reason: `no source in the crate reads std::env::var("PORT"), and neither axum nor actix-web takes a port any other way, so the service would listen on the address its own source fixes`,
		}}
	}
	return []Service{service("backend", "cargo run", "Cargo.toml")}, nil
}

// dependsOnRustServer reports whether the manifest names one of the two
// frameworks. A dependency alone is not enough to propose a command, which is
// why the caller goes on asking; it is enough to stop looking at every crate.
func dependsOnRustServer(manifest string) bool {
	for _, crate := range rustServerCrates {
		if regexp.MustCompile(`(?m)^\s*` + regexp.QuoteMeta(crate) + `\s*=`).MatchString(manifest) {
			return true
		}
	}
	return false
}

// ambiguousRustBinary reports why cargo run could not choose, or nothing where
// it can.
//
// A binary target comes from src/main.rs, from any file in src/bin, or from a
// [[bin]] table. With more than one and no default-run, cargo refuses without a
// --bin flag, so grat does not propose a command that would refuse.
func ambiguousRustBinary(root string, manifest string) string {
	targets := 0
	if fileExists(join(root, "src", "main.rs")) {
		targets++
	}
	for _, entry := range entries(join(root, "src", "bin")) {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".rs") {
			targets++
		}
	}
	if rustBinaryTarget.MatchString(manifest) {
		targets++
	}

	switch {
	case targets == 0:
		return "the crate depends on a web framework but declares no binary to run, so it is a library"
	case targets > 1 && !rustDefaultRun.MatchString(manifest):
		return "the crate has more than one binary and names no default-run, so cargo run would refuse without a --bin flag"
	default:
		return ""
	}
}

// readsPortFromRustSource reports whether anything below src reads the port out
// of the environment.
func readsPortFromRustSource(root string) bool {
	found := false
	var walk func(directory string, depth int)
	walk = func(directory string, depth int) {
		if found || depth > rustSourceDepth {
			return
		}
		for _, entry := range entries(directory) {
			if found {
				return
			}
			if entry.IsDir() {
				walk(join(directory, entry.Name()), depth+1)
				continue
			}
			if !strings.HasSuffix(entry.Name(), ".rs") {
				continue
			}
			data, ok := readBounded(join(directory, entry.Name()))
			if !ok {
				continue
			}
			source := rustBlockComment.ReplaceAllString(string(data), "")
			source = rustLineComment.ReplaceAllString(source, "")
			if rustPortRead.MatchString(source) {
				found = true
			}
		}
	}
	walk(join(root, "src"), 0)
	return found
}

// rustSourceDepth bounds the walk below src. A crate's modules sit a few levels
// down at most, and the bound is what keeps a vendored tree from costing the
// whole of a discovery run.
const rustSourceDepth = 6
