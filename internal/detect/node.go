package detect

import (
	"encoding/json"
	"strings"
)

// nodeScripts maps a service name to the development scripts that would start
// it, in the order they are preferred. The first script that exists wins, so a
// project carrying both `dev:frontend` and `dev` yields one service rather than
// two.
var nodeScripts = []struct {
	name    string
	scripts []string
}{
	{name: "shared", scripts: []string{"dev:shared"}},
	{name: "backend", scripts: []string{"dev:backend"}},
	{name: "frontend", scripts: []string{"dev:frontend", "dev"}},
	{name: "developer", scripts: []string{"dev:developer"}},
	{name: "dashboard", scripts: []string{"dev:dashboard"}},
}

// detectNode reads package.json and derives one service per conventional
// development script it declares.
//
// The package manager follows the manifest rather than a guess: a
// `packageManager` field naming pnpm, or a pnpm lockfile beside it, selects
// pnpm, and everything else uses npm.
func detectNode(root string) ([]Service, []Unresolved) {
	data, ok := readBounded(join(root, "package.json"))
	if !ok {
		return nil, nil
	}

	var manifest struct {
		PackageManager string            `json:"packageManager"`
		Scripts        map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, []Unresolved{{Marker: "package.json", Reason: "the manifest is not readable JSON"}}
	}
	if len(manifest.Scripts) == 0 {
		return nil, []Unresolved{{Marker: "package.json", Reason: "the manifest declares no scripts"}}
	}

	runner := "npm run"
	if strings.HasPrefix(manifest.PackageManager, "pnpm@") || fileExists(join(root, "pnpm-lock.yaml")) {
		runner = "pnpm"
	}

	services := make([]Service, 0, len(nodeScripts))
	for _, candidate := range nodeScripts {
		for _, script := range candidate.scripts {
			if _, exists := manifest.Scripts[script]; !exists {
				continue
			}
			services = append(services, service(candidate.name, runner+" "+script, "package.json"))
			break
		}
	}
	if len(services) == 0 {
		return nil, []Unresolved{{
			Marker: "package.json",
			Reason: "none of the conventional development scripts are declared",
		}}
	}
	return services, nil
}
