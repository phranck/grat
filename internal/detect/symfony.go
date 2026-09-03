package detect

import (
	"encoding/json"
)

// symfonyBundle is the package every Symfony application requires and no
// library does. Composer alone says a project is PHP; this says it serves.
const symfonyBundle = "symfony/framework-bundle"

// symfonyServeCommand is Symfony's own local server.
//
// It blocks the terminal by default, which is what grat needs; -d is the flag
// that would put it in the background and is deliberately absent. --no-tls
// keeps it on http, because grat probes an http health path, and it listens on
// 127.0.0.1 by default rather than on every interface.
const symfonyServeCommand = "symfony server:start --no-tls --port=$PORT"

// detectSymfony recognises a Symfony application by its own bundle.
//
// Only the Symfony CLI command is proposed. `php -S` would serve a Symfony
// application after a fashion and is the obvious fallback, but grat reads files
// and never runs anything, so it cannot tell whether the CLI is installed; a
// second command offered on a guess is a command that fails at the first start
// with nothing to point at. Somebody without the CLI writes one line instead.
func detectSymfony(root string) ([]Service, []Unresolved) {
	data, ok := readBounded(join(root, "composer.json"))
	if !ok {
		return nil, nil
	}

	var manifest struct {
		Require map[string]string `json:"require"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		// Laravel's detector reports an unreadable manifest already, and one
		// project has one manifest, so saying it twice would be saying it twice.
		return nil, nil
	}
	if _, required := manifest.Require[symfonyBundle]; !required {
		return nil, nil
	}
	return []Service{service("backend", symfonyServeCommand, "composer.json")}, nil
}
