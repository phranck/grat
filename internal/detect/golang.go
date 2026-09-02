package detect

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/phranck/grat/internal/project"
)

// mainPackagePattern matches the package clause of a runnable program, allowing
// for the build constraints and comments that can precede it.
var mainPackagePattern = regexp.MustCompile(`(?m)^package\s+main\b`)

// portReaders are the two functions of the os package that read an environment
// variable.
var portReaders = map[string]struct{}{"Getenv": {}, "LookupEnv": {}}

// skippedGoDirectories are what a Go module keeps that the shared project scan
// has no reason to know about. Everything else it leaves out, such as vendor,
// node_modules and .git, is already in project.SkipsScanning, and repeating a
// name here would mean two lists that can disagree.
//
// testdata is Go's own convention for input the toolchain ignores, so what it
// declares says nothing about what this module runs.
var skippedGoDirectories = map[string]struct{}{"testdata": {}}

// goProgram is one runnable program of a module.
type goProgram struct {
	// name becomes the service name, and with it the role.
	name string
	// path is the package path passed to go run.
	path string
	// evidence names what the program was found through.
	evidence string
}

// detectGo recognises a Go module and the programs it can run.
//
// A module is only a marker; what can be started is a main package. Those live
// under cmd by convention, and the directory name there becomes the service
// name. A module whose main package sits at the root yields one service named
// after the module directory.
//
// A program is only offered where the module reads the port from the
// environment. Go has no framework and therefore no flag, so a program listens
// on whatever address its own source names. Without that read, grat would wait
// for a port nothing binds and the service would never become ready.
func detectGo(root string) ([]Service, []Unresolved) {
	if _, ok := readBounded(join(root, "go.mod")); !ok {
		return nil, nil
	}

	programs := goPrograms(root)
	if len(programs) == 0 {
		return nil, []Unresolved{{
			Marker: "go.mod",
			Reason: "the module declares no main package to run, at the root or below cmd",
		}}
	}
	if !readsPortFromGoSource(root) {
		return nil, []Unresolved{{
			Marker: "go.mod",
			Reason: `no source in the module calls os.Getenv("PORT"), so the programs would ignore the port grat assigns`,
		}}
	}

	services := make([]Service, 0, len(programs))
	for _, program := range programs {
		services = append(services, service(program.name, "go run "+program.path, program.evidence))
	}
	return services, nil
}

// goPrograms lists what the module can run, preferring the programs under cmd
// because a module that has them keeps its root for the library.
func goPrograms(root string) []goProgram {
	programs := make([]goProgram, 0, 2)
	for _, entry := range entries(join(root, "cmd")) {
		if !entry.IsDir() || !holdsMainPackage(join(root, "cmd", entry.Name())) {
			continue
		}
		programs = append(programs, goProgram{
			name:     entry.Name(),
			path:     "./cmd/" + entry.Name(),
			evidence: "cmd/" + entry.Name(),
		})
	}
	if len(programs) > 0 {
		return programs
	}

	if holdsMainPackage(root) {
		return []goProgram{{name: filepath.Base(root), path: ".", evidence: "go.mod"}}
	}
	return nil
}

// holdsMainPackage reports whether any Go file directly in directory declares
// package main. Only that directory is read, because a main package in a
// subdirectory is a different program with its own path.
func holdsMainPackage(directory string) bool {
	for _, entry := range entries(directory) {
		if entry.IsDir() || !goSourceFile(entry.Name()) {
			continue
		}
		data, ok := readBounded(filepath.Join(directory, entry.Name()))
		if !ok {
			continue
		}
		if mainPackagePattern.Match(data) {
			return true
		}
	}
	return false
}

// readsPortFromGoSource reports whether any source in the module reads the port
// from the environment.
//
// The whole module is searched rather than the program's own directory, because
// a Go server conventionally builds its configuration in a package of its own
// and the main function only calls it. Searching one directory would answer the
// wrong question and miss almost every real project.
//
// The syntax tree decides rather than the text, so that the name appearing in a
// comment, an error message or a documentation line is not read as a call. A
// text search reports every tool that merely writes about the variable, and grat
// itself is one of them.
//
// It walks through project.Walk, which is the bounded traversal every other scan
// uses. Its own walk read the whole module without a limit on entries or on
// depth, so grat discover over a folder holding a large unpacked module tree
// spent its time in there.
func readsPortFromGoSource(root string) bool {
	found := false
	fileSet := token.NewFileSet()
	_, _ = project.Walk(root, project.MaxScanEntries, func(path string, entry fs.DirEntry) error {
		if found {
			return filepath.SkipAll
		}
		if entry.IsDir() {
			if _, skipped := skippedGoDirectories[entry.Name()]; skipped {
				return fs.SkipDir
			}
			return nil
		}
		if !goSourceFile(entry.Name()) {
			return nil
		}
		data, ok := readBounded(path)
		if !ok {
			return nil
		}
		file, parseError := parser.ParseFile(fileSet, path, data, parser.SkipObjectResolution)
		if parseError != nil {
			// A file that does not parse says nothing either way.
			return nil
		}
		ast.Inspect(file, func(node ast.Node) bool {
			if found {
				return false
			}
			if isPortEnvironmentCall(node) {
				found = true
				return false
			}
			return true
		})
		return nil
	})
	return found
}

// isPortEnvironmentCall reports whether node is a call to os.Getenv("PORT") or
// os.LookupEnv("PORT").
func isPortEnvironmentCall(node ast.Node) bool {
	call, isCall := node.(*ast.CallExpr)
	if !isCall || len(call.Args) != 1 {
		return false
	}
	selector, isSelector := call.Fun.(*ast.SelectorExpr)
	if !isSelector {
		return false
	}
	packageName, isIdentifier := selector.X.(*ast.Ident)
	if !isIdentifier || packageName.Name != "os" {
		return false
	}
	if _, reads := portReaders[selector.Sel.Name]; !reads {
		return false
	}
	literal, isLiteral := call.Args[0].(*ast.BasicLit)
	if !isLiteral || literal.Kind != token.STRING {
		return false
	}
	name, unquoteError := strconv.Unquote(literal.Value)
	return unquoteError == nil && name == "PORT"
}

// goSourceFile reports whether a name is a Go source file this reads. Test files
// are left out, because what they declare says nothing about what the module
// runs.
func goSourceFile(name string) bool {
	return filepath.Ext(name) == ".go" && !strings.HasSuffix(name, "_test.go")
}
