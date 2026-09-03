package detect

import (
	"strings"
	"testing"
)

// commandFor returns the single command detected for a project, or fails saying
// what came back instead.
func commandFor(t *testing.T, files map[string]string) string {
	t.Helper()
	finding := Directory(writeProject(t, files))
	if len(finding.Services) != 1 {
		t.Fatalf("services = %+v, unresolved = %+v; want exactly one", finding.Services, finding.Unresolved)
	}
	return finding.Services[0].Command
}

// TestSymfonyIsStartedByItsOwnServer covers the stack whose server needs no
// extra reading: the bundle says it serves, and the CLI takes the port.
func TestSymfonyIsStartedByItsOwnServer(t *testing.T) {
	t.Parallel()

	command := commandFor(t, map[string]string{
		"composer.json": `{"require": {"php": "^8.2", "symfony/framework-bundle": "^7.0"}}`,
	})
	for _, wanted := range []string{"symfony server:start", "--no-tls", "--port=$PORT"} {
		if !strings.Contains(command, wanted) {
			t.Fatalf("command = %q, want it to carry %q", command, wanted)
		}
	}
	// -d would put the server in the background, where grat has nothing to watch.
	if strings.Contains(command, " -d") {
		t.Fatalf("command = %q, and -d daemonises the server", command)
	}
}

// TestAPhpLibraryIsNotASymfonyApplication keeps the marker from firing on any
// project that merely uses Composer.
func TestAPhpLibraryIsNotASymfonyApplication(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, map[string]string{
		"composer.json": `{"require": {"php": "^8.2", "psr/log": "^3.0"}}`,
	}))
	if finding.Any() {
		t.Fatalf("a Composer library was recognised: %+v", finding)
	}
}

// TestPhoenixIsStartedWhereItReadsThePort covers the stack that needs a second
// file read, because Phoenix has no port flag at all.
func TestPhoenixIsStartedWhereItReadsThePort(t *testing.T) {
	t.Parallel()

	command := commandFor(t, map[string]string{
		"mix.exs": `defmodule Example.MixProject do
  use Mix.Project
  defp deps do
    [{:phoenix, "~> 1.8"}]
  end
end
`,
		"config/runtime.exs": `import Config
config :example, ExampleWeb.Endpoint,
  http: [port: String.to_integer(System.get_env("PORT", "4000"))]
`,
	})
	if command != phoenixServeCommand {
		t.Fatalf("command = %q, want %q", command, phoenixServeCommand)
	}
}

// TestAPhoenixProjectThatIgnoresThePortIsReported is the case the issue asks
// for. Such a project runs perfectly on its own fixed port and never answers on
// the one grat assigned, which is the failure worth naming rather than guessing
// at.
func TestAPhoenixProjectThatIgnoresThePortIsReported(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, map[string]string{
		"mix.exs": `defmodule Example.MixProject do
  use Mix.Project
  defp deps do
    [{:phoenix, "~> 1.8"}]
  end
end
`,
		"config/runtime.exs": `import Config
config :example, ExampleWeb.Endpoint, http: [port: 4000]
`,
	}))
	if len(finding.Services) != 0 {
		t.Fatalf("a command was proposed for a project that ignores the port: %+v", finding.Services)
	}
	if len(finding.Unresolved) != 1 {
		t.Fatalf("unresolved = %+v, want exactly one", finding.Unresolved)
	}
	if !strings.Contains(finding.Unresolved[0].Reason, "PORT") {
		t.Fatalf("the reason does not say what is missing: %+v", finding.Unresolved[0])
	}
}

// TestSpringBootOnGradleCarriesNoDaemon is the one flag this whole detector
// turns on. bootRun forks the application either way; the flag decides whose
// child it is, and without it the application hangs below Gradle's long-lived
// daemon, which grat did not start.
func TestSpringBootOnGradleCarriesNoDaemon(t *testing.T) {
	t.Parallel()

	for _, manifest := range []string{"build.gradle", "build.gradle.kts"} {
		command := commandFor(t, map[string]string{
			manifest:  `plugins { id("org.springframework.boot") version "3.4.0" }`,
			"gradlew": "#!/bin/sh\nexec gradle \"$@\"\n",
		})
		if !strings.Contains(command, "--no-daemon") {
			t.Fatalf("%s: command = %q, and without --no-daemon the application runs below Gradle's daemon", manifest, command)
		}
		if !strings.Contains(command, "SERVER_PORT=$PORT") {
			t.Fatalf("%s: command = %q, and Spring Boot reads SERVER_PORT rather than PORT", manifest, command)
		}
		if !strings.Contains(command, "./gradlew bootRun") {
			t.Fatalf("%s: command = %q, want the wrapper", manifest, command)
		}
	}
}

// TestSpringBootOnMavenForksItsOwn covers the other build tool, where the
// plugin forks a child JVM of mvn and no flag is needed.
func TestSpringBootOnMavenForksItsOwn(t *testing.T) {
	t.Parallel()

	command := commandFor(t, map[string]string{
		"pom.xml": `<project><parent><artifactId>spring-boot-starter-parent</artifactId></parent></project>`,
	})
	if command != springBootMavenCommand {
		t.Fatalf("command = %q, want %q", command, springBootMavenCommand)
	}
}

// TestSpringBootOnGradleWithoutAWrapperIsReported keeps grat from proposing a
// command that depends on a Gradle somewhere else on the machine.
func TestSpringBootOnGradleWithoutAWrapperIsReported(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, map[string]string{
		"build.gradle": `plugins { id("org.springframework.boot") version "3.4.0" }`,
	}))
	if len(finding.Services) != 0 {
		t.Fatalf("a command was proposed without a wrapper: %+v", finding.Services)
	}
	if len(finding.Unresolved) != 1 || !strings.Contains(finding.Unresolved[0].Reason, "gradlew") {
		t.Fatalf("unresolved = %+v, want the missing wrapper named", finding.Unresolved)
	}
}

// TestAPlainMavenProjectIsNotSpringBoot keeps the marker from firing on any
// project that merely uses Maven.
func TestAPlainMavenProjectIsNotSpringBoot(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, map[string]string{
		"pom.xml": `<project><artifactId>example</artifactId></project>`,
	}))
	if finding.Any() {
		t.Fatalf("a plain Maven project was recognised: %+v", finding)
	}
}
