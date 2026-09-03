package detect

import (
	"strings"
	"testing"
)

// TestADetectedNameNeverBecomesASecondCommand is the defect this guards.
// Detection builds command lines out of what a repository's own files say, and
// those files come with the repository. A name carrying a semicolon used to be
// copied straight into a line that grat start later runs through /bin/sh.
func TestADetectedNameNeverBecomesASecondCommand(t *testing.T) {
	t.Parallel()

	for name, project := range map[string]map[string]string{
		"a Swift executable target": {
			"Package.swift": `// swift-tools-version:5.9
import PackageDescription
let package = Package(
    name: "app",
    dependencies: [.package(url: "https://github.com/vapor/vapor.git", from: "4.0.0")],
    targets: [.executableTarget(name: "App; touch /tmp/pwned #")]
)
`,
		},
		"a directory below cmd": {
			"go.mod":              "module example.com/service\n\ngo 1.25\n",
			"cmd/api; id/main.go": "package main\n\nimport \"os\"\n\nfunc main() { _ = os.Getenv(\"PORT\") }\n",
		},
		"a Python application module": {
			"requirements.txt":     "fastapi\nuvicorn\n",
			"app; touch pwned;.py": "from fastapi import FastAPI\n\napp = FastAPI()\n",
		},
	} {
		root := writeProject(t, project)
		finding := Directory(root)
		services, unresolved := finding.Services, finding.Unresolved

		for _, service := range services {
			if strings.ContainsAny(service.Command, ";&|`$()") && !strings.Contains(service.Command, "$PORT") {
				t.Fatalf("%s: the command carries shell punctuation: %q", name, service.Command)
			}
			for _, character := range ";&|`" {
				if strings.ContainsRune(service.Command, character) {
					t.Fatalf("%s: the command carries %q: %q", name, character, service.Command)
				}
			}
		}
		if len(unresolved) == 0 && len(services) > 0 {
			t.Fatalf("%s: a command was proposed from an unsafe name: %+v", name, services)
		}
	}
}

// TestAnUnsafeNameSaysWhichCharacter is what makes the refusal actionable. A
// project that gets no command needs to know which character to change, and
// several of the ones that matter cannot be seen in a message.
func TestAnUnsafeNameSaysWhichCharacter(t *testing.T) {
	t.Parallel()

	root := writeProject(t, map[string]string{
		"Package.swift": `import PackageDescription
let package = Package(
    name: "app",
    dependencies: [.package(url: "https://github.com/vapor/vapor.git", from: "4.0.0")],
    targets: [.executableTarget(name: "App; id")]
)
`,
	})
	unresolved := Directory(root).Unresolved
	if len(unresolved) != 1 {
		t.Fatalf("unresolved = %+v, want exactly one", unresolved)
	}
	for _, wanted := range []string{"Package.swift", "App; id", "the character ;"} {
		if !strings.Contains(unresolved[0].Marker+" "+unresolved[0].Reason, wanted) {
			t.Fatalf("the reason does not carry %q: %+v", wanted, unresolved[0])
		}
	}
}

// TestAnOrdinaryNameStillYieldsACommand keeps the check from refusing what it
// is there to let through.
func TestAnOrdinaryNameStillYieldsACommand(t *testing.T) {
	t.Parallel()

	root := writeProject(t, map[string]string{
		"Package.swift": `import PackageDescription
let package = Package(
    name: "app",
    dependencies: [.package(url: "https://github.com/vapor/vapor.git", from: "4.0.0")],
    targets: [.executableTarget(name: "App")]
)
`,
	})
	finding := Directory(root)
	services, unresolved := finding.Services, finding.Unresolved
	if len(services) != 1 || !strings.Contains(services[0].Command, "swift run App serve") {
		t.Fatalf("services = %+v, unresolved = %+v; want the ordinary command", services, unresolved)
	}
}

// TestTheCharacterSetIsWhatItSays pins the set itself, since it is the whole of
// the rule and everything else reads from it.
func TestTheCharacterSetIsWhatItSays(t *testing.T) {
	t.Parallel()

	for _, safe := range []string{"App", "my-service", "my_service", "app.v2", "Api123"} {
		if _, ok := safeIdentifier(safe); !ok {
			t.Fatalf("%q was refused and is an ordinary name", safe)
		}
	}
	for _, unsafe := range []string{"", "a;b", "a b", "a`b", "a$b", "a|b", "a&b", "a(b", "a\nb", "../etc"} {
		if _, ok := safeIdentifier(unsafe); ok {
			t.Fatalf("%q was accepted and cannot go into a command", unsafe)
		}
	}
}
