package detect

import "path/filepath"

// This file holds the two tools that sit beside an application rather than
// instead of one. A repository with a Storybook usually also holds the component
// library it renders, and a repository with a Docusaurus site often holds the
// thing it documents, so both are reported in addition to whatever else the
// directory yields. That falls out of the detector registry, where every
// detector is asked and the answers are collected.
//
// Both are named for what they are rather than for a product role, and both
// therefore land in the developer range, which keeps them out of the lane the
// application itself allocates from.

// docusaurusConfigNames are the root configuration file names Docusaurus
// accepts, in the order it tries them.
var docusaurusConfigNames = []string{
	"docusaurus.config.ts",
	"docusaurus.config.mts",
	"docusaurus.config.cts",
	"docusaurus.config.js",
	"docusaurus.config.mjs",
	"docusaurus.config.cjs",
}

// docusaurusDependency is the package a Docusaurus site installs. The site's own
// name field says nothing, because a template leaves it arbitrary.
const docusaurusDependency = "@docusaurus/core"

// storybookConfigDirectory is where Storybook keeps its configuration.
const storybookConfigDirectory = ".storybook"

// storybookConfigExtensions are the extensions Storybook accepts for its main
// configuration file.
var storybookConfigExtensions = []string{
	".js", ".ts", ".jsx", ".tsx", ".mjs", ".mts", ".mtsx", ".cjs", ".cts", ".ctsx",
}

// storybookConfigNames are those extensions as paths below the project root.
var storybookConfigNames = storybookConfigFileNames()

// storybookPackage is the single package that has carried the command line
// since Storybook 7, and which absorbed the framework packages in 9.
const storybookPackage = "storybook"

// storybookScope covers the older shape, where the command line came from a
// framework package such as @storybook/react.
const storybookScope = "@storybook/"

// storybookConfigFileNames builds the configuration paths once, so the extension
// list stays the one thing that states them.
func storybookConfigFileNames() []string {
	names := make([]string, 0, len(storybookConfigExtensions))
	for _, extension := range storybookConfigExtensions {
		names = append(names, filepath.Join(storybookConfigDirectory, "main"+extension))
	}
	return names
}

// detectDocusaurus recognises a Docusaurus documentation site.
//
// The port and the host are both stated. Docusaurus reads PORT as well, and the
// flag wins over it, so the command says what it does without depending on what
// reaches its environment.
//
// It never relocates unattended: on a taken port with no terminal it prints an
// error and exits without binding, rather than moving as it does when somebody
// is there to answer the prompt. grat then sees the port stay unheld and reports
// the failure with the log, which is what covers the case where the exit status
// says nothing, because that path exits zero.
func detectDocusaurus(root string) ([]Service, []Unresolved) {
	value, ok := companionManifest(root)
	if !ok {
		return nil, nil
	}
	marker := firstExisting(root, docusaurusConfigNames)
	if marker == "" {
		return nil, nil
	}
	if !value.declares(docusaurusDependency) {
		return nil, []Unresolved{{
			Marker: marker,
			Reason: "the site is configured for Docusaurus but the manifest does not declare " +
				docusaurusDependency + ", so there is no command to run",
		}}
	}

	command := value.binaryRunner(root) + " docusaurus start --port $PORT --host 127.0.0.1"
	return []Service{service("docs", command, marker)}, nil
}

// detectStorybook recognises a Storybook.
//
// --exact-port is the whole reason this command is safe to write. Storybook
// otherwise runs the requested port through detect-port, which tries the next
// nine and then falls back to a random one, and the prompt that would catch it
// is suppressed whenever CI is set or no port was given. With the flag it exits
// instead, and it does so before the configuration is loaded, so a taken port is
// reported as a taken port.
//
// -h is doing real work too. Storybook passes its host option straight to listen
// and never defaults it, so without the flag Node binds every interface.
func detectStorybook(root string) ([]Service, []Unresolved) {
	value, ok := companionManifest(root)
	if !ok {
		return nil, nil
	}
	marker := firstExisting(root, storybookConfigNames)
	if marker == "" {
		return nil, nil
	}
	if !value.declares(storybookPackage) && !value.declaresWithin(storybookScope) {
		return nil, []Unresolved{{
			Marker: marker,
			Reason: "the project is configured for Storybook but the manifest declares no storybook package, " +
				"so there is no command line to run",
		}}
	}

	command := value.binaryRunner(root) + " storybook dev -p $PORT -h 127.0.0.1 --exact-port"
	return []Service{service("storybook", command, marker)}, nil
}

// companionManifest reads package.json for the two detectors above.
//
// A manifest that cannot be parsed is reported by the Node detector, which reads
// the same file, so it is passed over here rather than reported three times for
// one broken file.
func companionManifest(root string) (manifest, bool) {
	value, _, ok := readManifest(root)
	return value, ok
}

// servesOnlyCompanion reports whether a directory holds one of these tools and
// nothing else that starts a server.
//
// It exists for the Node detector's last resort, which builds a command from the
// conventional dev script and therefore carries no port. Where that script is
// the Storybook or the documentation site the detectors above have already
// answered with a port, grat would otherwise offer the same service twice and
// wait forever on the copy that cannot take one.
func servesOnlyCompanion(root string, value manifest) bool {
	docusaurus := firstExisting(root, docusaurusConfigNames) != "" && value.declares(docusaurusDependency)
	storybook := firstExisting(root, storybookConfigNames) != "" &&
		(value.declares(storybookPackage) || value.declaresWithin(storybookScope))
	return docusaurus || storybook
}
