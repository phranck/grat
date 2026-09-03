package detect

import (
	"strings"
	"testing"

	"github.com/phranck/grat/internal/config"
)

// TestAStorybookIsReportedBesideTheApplication is the case this pair exists for:
// a repository holding a component library and its Storybook yields both, rather
// than the first detector that matches winning.
func TestAStorybookIsReportedBesideTheApplication(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, map[string]string{
		"package.json":       `{"name": "library", "devDependencies": {"vite": "^7.0.0", "storybook": "^10.0.0", "@storybook/react-vite": "^10.0.0"}}`,
		".storybook/main.ts": "export default { framework: '@storybook/react-vite' };\n",
		"vite.config.ts":     "export default {};\n",
	}))
	if len(finding.Services) != 2 {
		t.Fatalf("services = %+v; want the application and the Storybook", finding.Services)
	}

	byName := map[string]Service{}
	for _, found := range finding.Services {
		byName[found.Name] = found
	}
	application, ok := byName["frontend"]
	if !ok {
		t.Fatalf("services = %+v; want the application under frontend", finding.Services)
	}
	if !strings.Contains(application.Command, "vite dev") {
		t.Fatalf("application command = %q, want the Vite server", application.Command)
	}

	storybook, ok := byName["storybook"]
	if !ok {
		t.Fatalf("services = %+v; want the Storybook", finding.Services)
	}
	if storybook.Role != config.RoleDeveloper {
		t.Fatalf("storybook role = %q, want the developer lane rather than the application's", storybook.Role)
	}
	for _, wanted := range []string{"storybook dev", "-p $PORT", "-h 127.0.0.1", "--exact-port"} {
		if !strings.Contains(storybook.Command, wanted) {
			t.Fatalf("storybook command = %q, want it to carry %q", storybook.Command, wanted)
		}
	}
}

// TestTheStorybookCommandRefusesToMovePorts pins the one flag without which
// Storybook relocates on a taken port and grat waits on one nothing serves.
func TestTheStorybookCommandRefusesToMovePorts(t *testing.T) {
	t.Parallel()

	command := commandFor(t, map[string]string{
		"package.json":       `{"name": "workshop", "devDependencies": {"storybook": "^10.0.0"}}`,
		".storybook/main.js": "export default {};\n",
	})
	if !strings.Contains(command, "--exact-port") {
		t.Fatalf("command = %q, and without --exact-port Storybook walks to another port", command)
	}
}

// TestAnOlderStorybookIsRecognisedByItsScopedPackage covers the shape from
// before the framework packages were consolidated into one.
func TestAnOlderStorybookIsRecognisedByItsScopedPackage(t *testing.T) {
	t.Parallel()

	command := commandFor(t, map[string]string{
		"package.json":       `{"name": "workshop", "devDependencies": {"@storybook/react": "^8.0.0"}}`,
		".storybook/main.js": "module.exports = {};\n",
	})
	if !strings.Contains(command, "storybook dev") {
		t.Fatalf("command = %q, want the Storybook server", command)
	}
}

// TestADocusaurusSiteIsStartedOnItsOwnPort covers the documentation site, whose
// service is named for what it is rather than for a product role.
func TestADocusaurusSiteIsStartedOnItsOwnPort(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, map[string]string{
		"package.json":         `{"name": "site", "dependencies": {"@docusaurus/core": "^3.10.0", "@docusaurus/preset-classic": "^3.10.0"}}`,
		"docusaurus.config.ts": "export default { title: 'Example' };\n",
	}))
	if len(finding.Services) != 1 {
		t.Fatalf("services = %+v; want just the documentation site", finding.Services)
	}
	found := finding.Services[0]
	if found.Name != "docs" || found.Role != config.RoleDeveloper {
		t.Fatalf("service = %+v, want docs in the developer lane", found)
	}
	for _, wanted := range []string{"docusaurus start", "--port $PORT", "--host 127.0.0.1"} {
		if !strings.Contains(found.Command, wanted) {
			t.Fatalf("command = %q, want it to carry %q", found.Command, wanted)
		}
	}
}

// TestADocusaurusConfigurationUnderAnyExtensionIsFound covers the names beyond
// the two obvious ones, which real sites use.
func TestADocusaurusConfigurationUnderAnyExtensionIsFound(t *testing.T) {
	t.Parallel()

	for _, name := range docusaurusConfigNames {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			command := commandFor(t, map[string]string{
				"package.json": `{"name": "site", "dependencies": {"@docusaurus/core": "^3.10.0"}}`,
				name:           "export default {};\n",
			})
			if !strings.Contains(command, "docusaurus start") {
				t.Fatalf("command = %q, want the Docusaurus server", command)
			}
		})
	}
}

// TestARepositoryHoldingOnlyASiteYieldsOneService covers the case the Node
// detector's last resort would answer a second time. Its dev script starts the
// documentation site the companion detector already answered, and it would
// carry no port.
func TestARepositoryHoldingOnlyASiteYieldsOneService(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, map[string]string{
		"package.json":         `{"name": "site", "scripts": {"dev": "docusaurus start"}, "dependencies": {"@docusaurus/core": "^3.10.0"}}`,
		"docusaurus.config.js": "module.exports = {};\n",
	}))
	if len(finding.Services) != 1 {
		t.Fatalf("services = %+v; want the site once, with its port", finding.Services)
	}
	if finding.Services[0].Name != "docs" {
		t.Fatalf("service = %+v, want the one carrying the port", finding.Services[0])
	}
}

// TestAConfigurationWithoutItsPackageIsReported covers both tools at once: the
// configuration file says what the project is, and the package says whether
// there is anything to run.
func TestAConfigurationWithoutItsPackageIsReported(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, map[string]string{
		"package.json":         `{"name": "site", "scripts": {"build": "echo"}}`,
		"docusaurus.config.js": "module.exports = {};\n",
		".storybook/main.js":   "module.exports = {};\n",
	}))
	if len(finding.Services) != 0 {
		t.Fatalf("services = %+v; want none without the packages", finding.Services)
	}

	reported := map[string]bool{}
	for _, unresolved := range finding.Unresolved {
		reported[unresolved.Marker] = true
	}
	for _, marker := range []string{"docusaurus.config.js", ".storybook/main.js"} {
		if !reported[marker] {
			t.Fatalf("unresolved = %+v; want %s reported", finding.Unresolved, marker)
		}
	}
}
