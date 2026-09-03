package manual

import (
	"strconv"
	"strings"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/runtime"
)

// ConfigPage renders the configuration manual as a man page for section 7,
// which is where a file format belongs.
func ConfigPage(version string, date string) string {
	return Roff(ConfigDocument(version, date))
}

// ConfigDocument builds grat.config(7).
//
// Every figure in it comes from the code that enforces it: the runtime defaults
// from config.DefaultRuntime, the ranges from Role.PortRange, the roles from
// config.Roles, and the
// inherited variables from runtime.InheritedEnvironment. A field described here
// and absent there would be a promise nothing keeps.
func ConfigDocument(version string, date string) Document {
	return Document{
		Name:          "grat.config",
		ManualSection: 7,
		Category:      "File Formats",
		Title:         "the declarative description of a project's services",
		Version:       version,
		Date:          date,
		Sections: []Section{
			{Title: "Name", Blocks: []Block{
				Prose("grat.config - the declarative description of a project's services"),
			}},
			{Title: "Description", Blocks: []Block{Prose(configDescription)}},
			{Title: "Example", Blocks: []Block{Literal(configExample)}},
			{Title: "Top level", Blocks: []Block{fieldList(topLevelFields)}},
			{Title: "The project table", Blocks: []Block{fieldList(projectFields)}},
			{Title: "The runtime table", Blocks: []Block{
				Prose(runtimeIntro),
				runtimeFields(),
			}},
			{Title: "A service", Blocks: []Block{
				Prose(serviceIntro),
				fieldList(serviceFields),
			}},
			{Title: "Roles and port ranges", Blocks: []Block{
				Prose(rolesIntro),
				roleRanges(),
			}},
			{Title: "The environment a command receives", Blocks: []Block{
				Prose(environmentIntro),
				Prose(strings.Join(runtime.InheritedEnvironment(), ", ") + "."),
				Prose(environmentRest),
			}},
			{Title: "See also", Blocks: []Block{Prose(configSeeAlso)}},
		},
	}
}

// field is one key of the file, as the reader meets it.
type field struct {
	name     string
	required string
	meaning  string
}

// fieldList renders a group of keys as a definition list.
func fieldList(fields []field) Block {
	items := make([]Item, 0, len(fields))
	for _, entry := range fields {
		items = append(items, Item{Term: entry.name, Detail: entry.required + " " + entry.meaning})
	}
	return Definitions(items...)
}

// runtimeFields builds the timing table from the defaults themselves, so a
// default changed in the code changes the page with it.
func runtimeFields() Block {
	defaults := config.DefaultRuntime()
	items := []Item{}
	for _, entry := range []struct {
		name    string
		value   string
		meaning string
	}{
		{"start_timeout", defaults.StartTimeout, "How long a selected service may take to reach readiness."},
		{"probe_interval", defaults.ProbeInterval, "The wait between one listener and health check and the next."},
		{"health_timeout", defaults.HealthTimeout, "How long one health request may take."},
		{"shutdown_timeout", defaults.ShutdownTimeout, "The grace after SIGTERM before SIGKILL follows."},
		{"log_tail_lines", strconv.Itoa(defaults.LogTailLines), "How many closing log lines a startup failure carries."},
	} {
		items = append(items, Item{
			Term:   entry.name,
			Detail: "Optional. " + entry.meaning + " The default is " + entry.value + ".",
		})
	}
	return Definitions(items...)
}

// roleRanges builds the range table from the roles themselves.
func roleRanges() Block {
	rows := [][]string{}
	for _, role := range config.Roles() {
		portRange, known := role.PortRange()
		switch {
		case !known:
			rows = append(rows, []string{string(role), "No range is allocated for this role."})
		case portRange.First == 0:
			rows = append(rows, []string{string(role), "No port. The service is watched as a process and is never probed over HTTP, so it takes port 0 and no health path."})
		default:
			rows = append(rows, []string{string(role), strconv.Itoa(portRange.First) + " to " + strconv.Itoa(portRange.Last)})
		}
	}
	return Rows([]string{"Role", "Port range"}, rows)
}
