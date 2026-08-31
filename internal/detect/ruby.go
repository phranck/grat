package detect

import "regexp"

// railsGemPattern matches the rails gem in a Gemfile, in either quoting style
// and with or without a version constraint.
var railsGemPattern = regexp.MustCompile(`(?m)^\s*gem\s+['"]rails['"]`)

// detectRails recognises a Rails application.
//
// Two things have to be present. The Gemfile has to declare the rails gem,
// which distinguishes an application from any other Ruby project, and bin/rails
// has to exist, which Rails generates and which is what the command runs.
//
// The command passes the port explicitly even though Rails reads PORT on its
// own, because a project carrying its own config/puma.rb can bind elsewhere and
// the flag is what settles it. Source for both: the Rails command line guide at
// guides.rubyonrails.org/command_line.html and the port resolution in
// railties/lib/rails/commands/server/server_command.rb.
func detectRails(root string) ([]Service, []Unresolved) {
	data, ok := readBounded(join(root, "Gemfile"))
	if !ok {
		return nil, nil
	}
	if !railsGemPattern.Match(data) {
		// A Ruby project that is not Rails has no server this can start.
		return nil, nil
	}
	if !fileExists(join(root, "bin", "rails")) {
		return nil, []Unresolved{{
			Marker: "Gemfile",
			Reason: "the Gemfile declares rails but bin/rails is missing, so the application is not generated yet",
		}}
	}

	return []Service{service("backend", "bin/rails server -b 127.0.0.1 -p $PORT", "bin/rails")}, nil
}
