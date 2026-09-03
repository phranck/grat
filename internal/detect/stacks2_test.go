package detect

import (
	"strings"
	"testing"
)

// TestAspNetCoreTakesTheUrlAfterTheSeparator covers the stack that settles the
// port completely from the command line.
func TestAspNetCoreTakesTheUrlAfterTheSeparator(t *testing.T) {
	t.Parallel()

	command := commandFor(t, map[string]string{
		"Example.csproj": `<Project Sdk="Microsoft.NET.Sdk.Web"><PropertyGroup><TargetFramework>net9.0</TargetFramework></PropertyGroup></Project>`,
	})
	// The separator is the point: everything after it reaches the application
	// unread, which is what the documentation asks for.
	if !strings.Contains(command, "dotnet run -- --urls") {
		t.Fatalf("command = %q, want the arguments after a bare separator", command)
	}
	if !strings.Contains(command, "127.0.0.1:$PORT") {
		t.Fatalf("command = %q, want the assigned port on the loopback address", command)
	}
}

// TestAConsoleProjectIsNotAWebApplication keeps the marker on the SDK rather
// than on the file extension.
func TestAConsoleProjectIsNotAWebApplication(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, map[string]string{
		"Example.csproj": `<Project Sdk="Microsoft.NET.Sdk"><PropertyGroup><OutputType>Exe</OutputType></PropertyGroup></Project>`,
	}))
	if finding.Any() {
		t.Fatalf("a console project was recognised: %+v", finding)
	}
}

// rustCrate is a crate with one binary that reads the port, which is the shape
// the tests below vary from.
func rustCrate(manifest string, main string) map[string]string {
	return map[string]string{"Cargo.toml": manifest, "src/main.rs": main}
}

const rustMainReadingPort = `use std::env;

#[tokio::main]
async fn main() {
    let port = env::var("PORT").unwrap_or_else(|_| "3000".to_string());
    let listener = tokio::net::TcpListener::bind(format!("127.0.0.1:{port}")).await.unwrap();
    axum::serve(listener, axum::Router::new()).await.unwrap();
}
`

// TestARustServiceThatReadsThePortIsStarted covers the stack that settles
// nothing, so the source has to answer.
func TestARustServiceThatReadsThePortIsStarted(t *testing.T) {
	t.Parallel()

	command := commandFor(t, rustCrate(
		"[package]\nname = \"example\"\n\n[dependencies]\naxum = \"0.8\"\n",
		rustMainReadingPort,
	))
	if command != "cargo run" {
		t.Fatalf("command = %q, want cargo run", command)
	}
}

// TestARustServiceThatFixesItsPortIsReported is the negative case. Such a crate
// runs perfectly on its own address and never answers where grat is waiting.
func TestARustServiceThatFixesItsPortIsReported(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, rustCrate(
		"[package]\nname = \"example\"\n\n[dependencies]\nactix-web = \"4\"\n",
		"fn main() {\n    // env::var(\"PORT\") would be read here one day\n    let _ = 8080;\n}\n",
	)))
	if len(finding.Services) != 0 {
		t.Fatalf("a command was proposed for a crate that fixes its port: %+v", finding.Services)
	}
	if len(finding.Unresolved) != 1 || !strings.Contains(finding.Unresolved[0].Reason, "PORT") {
		t.Fatalf("unresolved = %+v, want the missing read named", finding.Unresolved)
	}
}

// TestAMentionInARustCommentIsNotARead is why the comments come out before the
// search. grat detected itself that way once, through its own documentation.
func TestAMentionInARustCommentIsNotARead(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, rustCrate(
		"[package]\nname = \"example\"\n\n[dependencies]\naxum = \"0.8\"\n",
		"/* read the port with env::var(\"PORT\") some day */\nfn main() {\n    // env::var(\"PORT\")\n}\n",
	)))
	if len(finding.Services) != 0 {
		t.Fatalf("a mention in a comment was taken for a read: %+v", finding.Services)
	}
}

// TestARustCrateWithSeveralBinariesIsReported covers what cargo itself refuses.
func TestARustCrateWithSeveralBinariesIsReported(t *testing.T) {
	t.Parallel()

	files := rustCrate("[package]\nname = \"example\"\n\n[dependencies]\naxum = \"0.8\"\n", rustMainReadingPort)
	files["src/bin/worker.rs"] = "fn main() {}\n"
	finding := Directory(writeProject(t, files))
	if len(finding.Services) != 0 {
		t.Fatalf("a command was proposed where cargo run would refuse: %+v", finding.Services)
	}
	if len(finding.Unresolved) != 1 || !strings.Contains(finding.Unresolved[0].Reason, "default-run") {
		t.Fatalf("unresolved = %+v, want the ambiguity named", finding.Unresolved)
	}
}

// TestARustCrateWithADefaultRunIsStarted is the other side: the ambiguity is
// resolved in the manifest, so cargo run knows what to do and so does grat.
func TestARustCrateWithADefaultRunIsStarted(t *testing.T) {
	t.Parallel()

	files := rustCrate(
		"[package]\nname = \"example\"\ndefault-run = \"example\"\n\n[dependencies]\naxum = \"0.8\"\n",
		rustMainReadingPort,
	)
	files["src/bin/worker.rs"] = "fn main() {}\n"
	if command := commandFor(t, files); command != "cargo run" {
		t.Fatalf("command = %q, want cargo run", command)
	}
}

// TestFlaskNamesItsModule covers the stack that settles the port but not the
// entry point.
func TestFlaskNamesItsModule(t *testing.T) {
	t.Parallel()

	command := commandFor(t, map[string]string{
		"requirements.txt": "flask\n",
		"app.py":           "from flask import Flask\n\napp = Flask(__name__)\n",
	})
	for _, wanted := range []string{"flask --app app:app", "--host 127.0.0.1", "--port $PORT"} {
		if !strings.Contains(command, wanted) {
			t.Fatalf("command = %q, want it to carry %q", command, wanted)
		}
	}
}

// TestAFlaskProjectWithNoApplicationIsReported is the negative case: the
// dependency says it could serve, and nothing says what.
func TestAFlaskProjectWithNoApplicationIsReported(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, map[string]string{
		"requirements.txt": "flask\n",
		"helpers.py":       "def add(left, right):\n    return left + right\n",
	}))
	if len(finding.Services) != 0 {
		t.Fatalf("a command was proposed with no application: %+v", finding.Services)
	}
	if len(finding.Unresolved) != 1 || !strings.Contains(finding.Unresolved[0].Reason, "--app") {
		t.Fatalf("unresolved = %+v, want the missing module named", finding.Unresolved)
	}
}

// TestAFastApiProjectIsNotAlsoAFlaskOne keeps two detectors from proposing two
// services under the same name.
func TestAFastApiProjectIsNotAlsoAFlaskOne(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, map[string]string{
		"requirements.txt": "fastapi\nuvicorn\nflask\n",
		"main.py":          "from fastapi import FastAPI\n\napp = FastAPI()\n",
	}))
	if len(finding.Services) != 1 {
		t.Fatalf("services = %+v, want exactly one", finding.Services)
	}
	if !strings.Contains(finding.Services[0].Command, "uvicorn") {
		t.Fatalf("command = %q, want the FastAPI one", finding.Services[0].Command)
	}
}
