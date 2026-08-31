package detect

import "encoding/json"

// laravelPort is the command Laravel's own development server takes. It reads
// the port from the environment through $PORT, which is what grat sets.
const laravelServeCommand = "php artisan serve --host=127.0.0.1 --port=$PORT"

// detectLaravel recognises a Laravel application by the two things every one of
// them has: a Composer manifest and the `artisan` entry point beside it.
//
// Composer alone is not enough, because a PHP library has that too and has
// nothing to serve. The pair is what says this project can be started.
func detectLaravel(root string) ([]Service, []Unresolved) {
	manifestPath := join(root, "composer.json")
	data, ok := readBounded(manifestPath)
	if !ok {
		return nil, nil
	}
	if !fileExists(join(root, "artisan")) {
		// Composer without artisan is a library or a framework-less application.
		// Neither has a development server, so this is not a marker at all.
		return nil, nil
	}

	var manifest struct {
		Require map[string]string `json:"require"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, []Unresolved{{Marker: "composer.json", Reason: "the manifest is not readable JSON"}}
	}
	if _, exists := manifest.Require["laravel/framework"]; !exists {
		return nil, []Unresolved{{
			Marker: "composer.json",
			Reason: "an artisan entry point is present but the manifest does not require laravel/framework",
		}}
	}

	return []Service{service("backend", laravelServeCommand, "artisan")}, nil
}
