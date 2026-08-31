package detect

import (
	"path/filepath"
	"regexp"
)

// portFromEnvironment matches an application reading the port out of the
// environment, in the forms the three runtimes offer.
//
// This is the question that decides whether a Node server can be managed at
// all. None of NestJS, Express, Fastify or Deno reads PORT on its own, and none
// of them has a flag for it: their own documentation either fixes the port in
// the example or leaves the reading to the author. So the framework does not
// answer the question and the source does.
var portFromEnvironment = regexp.MustCompile(
	`process\.env\.PORT|Deno\.env\.get\(\s*['"]PORT['"]\s*\)|Bun\.env\.PORT`,
)

// nodeServers are the server frameworks recognised by their dependency.
var nodeServers = []struct {
	name       string
	dependency string
}{
	{name: "backend", dependency: "@nestjs/core"},
	{name: "backend", dependency: "fastify"},
	{name: "backend", dependency: "express"},
}

// startScripts are the scripts that run a server, in the order of preference.
// The one meant for development comes first, because it reloads.
var startScripts = []string{"start:dev", "dev", "start"}

// sourceDirectories are searched for a port being read, in the order a project
// is most likely to keep its entry point. Only these are read, and only one
// level deep, so a large repository is not walked for this question.
var sourceDirectories = []string{".", "src", "app", "server"}

// sourceExtensions are the files readsPortFromEnvironment opens.
var sourceExtensions = map[string]struct{}{
	".js": {}, ".mjs": {}, ".cjs": {}, ".ts": {}, ".mts": {}, ".cts": {},
}

// detectNodeServer recognises a Node server and decides whether grat can start
// it on an assigned port.
//
// Two things have to hold. A start script has to exist, because a hand-written
// server has no standard entry point and the manifest's main field describes a
// library rather than an application. And the source has to read the port from
// the environment, because grat communicates the port that way and nothing in
// these frameworks does it on their behalf.
func detectNodeServer(root string, value manifest) ([]Service, []Unresolved) {
	for _, candidate := range nodeServers {
		if !value.declares(candidate.dependency) {
			continue
		}

		script, exists := startScript(value)
		if !exists {
			return nil, []Unresolved{{
				Marker: "package.json",
				Reason: "a " + candidate.dependency + " server is declared but no dev or start script says how to run it",
			}}
		}
		if !readsPortFromEnvironment(root) {
			return nil, []Unresolved{{
				Marker: "package.json",
				Reason: "a " + candidate.dependency + " server is declared but nothing in the source reads process.env.PORT, so it would ignore the port grat assigns",
			}}
		}
		return []Service{service(candidate.name, value.scriptRunner(root)+" "+script, "package.json")}, nil
	}
	return nil, nil
}

// startScript returns the script that runs the server, preferring the one meant
// for development.
func startScript(value manifest) (string, bool) {
	for _, name := range startScripts {
		if _, exists := value.Scripts[name]; exists {
			return name, true
		}
	}
	return "", false
}

// readsPortFromEnvironment reports whether any source file near the top of the
// project reads the port out of the environment.
func readsPortFromEnvironment(root string) bool {
	for _, directory := range sourceDirectories {
		for _, entry := range entries(join(root, directory)) {
			if entry.IsDir() {
				continue
			}
			if _, wanted := sourceExtensions[filepath.Ext(entry.Name())]; !wanted {
				continue
			}
			data, ok := readBounded(filepath.Join(root, directory, entry.Name()))
			if !ok {
				continue
			}
			if portFromEnvironment.Match(data) {
				return true
			}
		}
	}
	return false
}
