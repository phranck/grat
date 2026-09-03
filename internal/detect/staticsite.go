package detect

import (
	"os"
	"regexp"
)

// hugoConfigNames are the root configuration files whose name identifies a Hugo
// site on its own, in the order Hugo reads them.
//
// Hugo takes the first one it finds and ignores the rest, so a site carrying two
// of them is configured by exactly one. The list comes from DefaultConfigNames
// and ValidConfigFileExtensions in config/configLoader.go, which is why .yml is
// here whilst the configuration documentation lists only .yaml.
var hugoConfigNames = []string{"hugo.toml", "hugo.yaml", "hugo.yml", "hugo.json"}

// hugoLegacyConfigNames are the older root names, which Hugo still reads after
// the ones above.
//
// They are not evidence on their own. config.json and config.yaml are names any
// project might use for something else, so one of them counts only beside a
// directory that Hugo alone puts in a project root.
var hugoLegacyConfigNames = []string{"config.toml", "config.yaml", "config.yml", "config.json"}

// hugoDirectories are the directories a Hugo site keeps beside its
// configuration, and one of them has to be there for a legacy name to count.
var hugoDirectories = []string{"content", "layouts", "archetypes"}

// jekyllConfigNames are the root configuration files that identify a Jekyll
// site, in the order Jekyll resolves them in lib/jekyll/configuration.rb.
var jekyllConfigNames = []string{"_config.yml", "_config.yaml", "_config.toml"}

// eleventyConfigNames are the file names Eleventy accepts for its configuration.
var eleventyConfigNames = []string{".eleventy.js", "eleventy.config.js", "eleventy.config.mjs", "eleventy.config.cjs"}

// eleventyDependency is the package that identifies an Eleventy site where no
// configuration file is present, which Eleventy permits.
const eleventyDependency = "@11ty/eleventy"

// jekyllGemPattern matches the jekyll gem in a Gemfile, in either quoting style
// and with or without a version constraint.
var jekyllGemPattern = regexp.MustCompile(`(?m)^\s*gem\s+['"]jekyll['"]`)

// detectHugo recognises a Hugo site.
//
// Hugo reads nothing from the environment, so both the port and the interface
// go on the command line. --bind is stated even though 127.0.0.1 is already the
// default, because the command is what a person reads to know where the site is
// reachable, and a default is not visible there.
func detectHugo(root string) ([]Service, []Unresolved) {
	marker := firstExisting(root, hugoConfigNames)
	if marker == "" {
		marker = hugoLegacyMarker(root)
	}
	if marker == "" {
		return nil, nil
	}
	return []Service{service("frontend", "hugo server --port $PORT --bind 127.0.0.1", marker)}, nil
}

// hugoLegacyMarker returns the older configuration name in use, but only where
// the project also holds a directory that belongs to Hugo.
func hugoLegacyMarker(root string) string {
	marker := firstExisting(root, hugoLegacyConfigNames)
	if marker == "" {
		return ""
	}
	for _, name := range hugoDirectories {
		if directoryExists(join(root, name)) {
			return marker
		}
	}
	return ""
}

// detectJekyll recognises a Jekyll site.
//
// The configuration file alone is not enough. _config.yml is a name other tools
// use as well, and the command runs through bundler, so what settles it is the
// Gemfile declaring the jekyll gem. A site whose configuration is there without
// the gem is reported rather than proposed, because bundle exec would fail.
//
// The port and the host are both given: Jekyll defaults to 4000 on 127.0.0.1
// and reads no environment variable. --detach is deliberately left out, because
// a server that backgrounds itself leaves grat nothing to watch.
func detectJekyll(root string) ([]Service, []Unresolved) {
	marker := firstExisting(root, jekyllConfigNames)
	if marker == "" {
		return nil, nil
	}
	data, ok := readBounded(join(root, "Gemfile"))
	if !ok || !jekyllGemPattern.Match(data) {
		return nil, []Unresolved{{
			Marker: marker,
			Reason: "the site is configured for Jekyll but no Gemfile declares the jekyll gem, so bundle exec has nothing to run",
		}}
	}

	command := "bundle exec jekyll serve --port $PORT --host 127.0.0.1"
	return []Service{service("frontend", command, marker)}, nil
}

// detectEleventy recognises an Eleventy site and reports why grat cannot manage
// one, which is the whole of what this detector does.
//
// Its development server fails both halves of what every other command grat
// writes guarantees, and neither is settleable from the command line. It calls
// Node's listen with a port and no host, so it accepts connections on every
// interface, and Eleventy's own argument parser rejects --host outright, which
// leaves no way to pin it to loopback. It then increments the port and retries
// up to ten times when the one it was given is taken, and there is no equivalent
// of Vite's --strictPort to stop it, so the server can answer somewhere grat
// never waits.
//
// Reporting that is the point. A person who reads this knows to run Eleventy
// themselves, whilst a proposed command would look like every other one in the
// configuration and quietly be the only one on the network.
func detectEleventy(root string, value manifest) []Unresolved {
	marker := eleventyMarker(root, value)
	if marker == "" {
		return nil
	}
	return []Unresolved{{
		Marker: marker,
		Reason: "Eleventy's development server binds every interface with no flag to pin it to 127.0.0.1, " +
			"and moves to another port when the assigned one is taken, so grat would wait on a port it does not serve",
	}}
}

// eleventyMarker names what identified the site, preferring the configuration
// file because a workspace holds one before its packages are installed.
func eleventyMarker(root string, value manifest) string {
	if marker := firstExisting(root, eleventyConfigNames); marker != "" {
		return marker
	}
	if value.declares(eleventyDependency) {
		return "package.json"
	}
	return ""
}

// firstExisting returns the first of names that is a regular file below root,
// which is how a tool with several accepted configuration names is identified
// by the one it actually uses.
func firstExisting(root string, names []string) string {
	for _, name := range names {
		if fileExists(join(root, name)) {
			return name
		}
	}
	return ""
}

// directoryExists reports whether path is a directory.
func directoryExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
