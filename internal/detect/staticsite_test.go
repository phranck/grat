package detect

import (
	"strings"
	"testing"
)

// TestHugoIsStartedFromItsConfiguration covers the simplest of the three static
// site generators: one configuration file settles it, and the command carries
// both the port and the interface.
func TestHugoIsStartedFromItsConfiguration(t *testing.T) {
	t.Parallel()

	command := commandFor(t, map[string]string{
		"hugo.toml": "baseURL = 'https://example.com/'\ntitle = 'Example'\n",
	})
	for _, wanted := range []string{"hugo server", "--port $PORT", "--bind 127.0.0.1"} {
		if !strings.Contains(command, wanted) {
			t.Fatalf("command = %q, want it to carry %q", command, wanted)
		}
	}
}

// TestHugoIsRecognisedUnderItsOlderConfigurationName covers the names Hugo still
// reads after the ones carrying its own name, which its documentation no longer
// mentions.
func TestHugoIsRecognisedUnderItsOlderConfigurationName(t *testing.T) {
	t.Parallel()

	command := commandFor(t, map[string]string{
		"config.toml":      "baseURL = 'https://example.com/'\n",
		"content/index.md": "# Example\n",
	})
	if !strings.Contains(command, "hugo server") {
		t.Fatalf("command = %q, want a Hugo server", command)
	}
}

// TestAGenericConfigFileIsNotAHugoSite keeps the older names from firing on any
// project that happens to keep a config.json in its root.
func TestAGenericConfigFileIsNotAHugoSite(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, map[string]string{
		"config.json": `{"setting": true}`,
	}))
	if finding.Any() {
		t.Fatalf("a bare config.json was recognised: %+v", finding)
	}
}

// TestJekyllRunsThroughBundler covers the stack whose command depends on a
// second file, because bundle exec is what the Gemfile makes possible.
func TestJekyllRunsThroughBundler(t *testing.T) {
	t.Parallel()

	command := commandFor(t, map[string]string{
		"_config.yml": "title: Example\n",
		"Gemfile":     "source 'https://rubygems.org'\ngem 'jekyll', '~> 4.3'\n",
	})
	for _, wanted := range []string{"bundle exec jekyll serve", "--port $PORT", "--host 127.0.0.1"} {
		if !strings.Contains(command, wanted) {
			t.Fatalf("command = %q, want it to carry %q", command, wanted)
		}
	}
	// --detach would put the server in the background, where grat has nothing
	// to watch.
	if strings.Contains(command, "--detach") {
		t.Fatalf("command = %q, and --detach backgrounds the server", command)
	}
}

// TestAJekyllConfigurationWithoutTheGemIsNotProposed covers the case the
// configuration file alone cannot settle: _config.yml is a name other tools use
// too, and without the gem bundle exec has nothing to run.
func TestAJekyllConfigurationWithoutTheGemIsNotProposed(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, map[string]string{
		"_config.yml": "title: Example\n",
		"Gemfile":     "source 'https://rubygems.org'\ngem 'rake'\n",
	}))
	if len(finding.Services) != 0 {
		t.Fatalf("services = %+v; want none without the jekyll gem", finding.Services)
	}
	if len(finding.Unresolved) != 1 {
		t.Fatalf("unresolved = %+v; want the missing gem reported", finding.Unresolved)
	}
	if !strings.Contains(finding.Unresolved[0].Reason, "jekyll gem") {
		t.Fatalf("reason = %q, want it to name the missing gem", finding.Unresolved[0].Reason)
	}
}

// TestEleventyIsReportedRatherThanProposed covers the third generator, which
// grat recognises in order to say why it cannot manage it.
func TestEleventyIsReportedRatherThanProposed(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, map[string]string{
		"package.json":       `{"name": "site", "scripts": {"dev": "eleventy --serve"}, "devDependencies": {"@11ty/eleventy": "^3.0.0"}}`,
		"eleventy.config.js": "export default function () {};\n",
	}))
	if len(finding.Services) != 0 {
		t.Fatalf("services = %+v; want no command for Eleventy", finding.Services)
	}
	if len(finding.Unresolved) != 1 {
		t.Fatalf("unresolved = %+v; want exactly the Eleventy report", finding.Unresolved)
	}
	if finding.Unresolved[0].Marker != "eleventy.config.js" {
		t.Fatalf("marker = %q, want the configuration file", finding.Unresolved[0].Marker)
	}
	for _, wanted := range []string{"every interface", "another port"} {
		if !strings.Contains(finding.Unresolved[0].Reason, wanted) {
			t.Fatalf("reason = %q, want it to name %q", finding.Unresolved[0].Reason, wanted)
		}
	}
}

// TestAnEleventyProjectWithoutAConfigurationIsStillRecognised covers the layout
// Eleventy permits, where the dependency is the only thing that says so.
func TestAnEleventyProjectWithoutAConfigurationIsStillRecognised(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, map[string]string{
		"package.json": `{"name": "site", "devDependencies": {"@11ty/eleventy": "^3.0.0"}}`,
	}))
	if len(finding.Unresolved) != 1 || finding.Unresolved[0].Marker != "package.json" {
		t.Fatalf("unresolved = %+v; want the manifest reported", finding.Unresolved)
	}
}
