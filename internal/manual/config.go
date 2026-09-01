package manual

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/runtime"
)

// ConfigPage renders the manual for the configuration file, for section 7 of
// the manual, which is where a file format belongs.
//
// Every figure in it comes from the code that enforces it: the runtime defaults
// from config.DefaultRuntime, the ranges from Role.PortRange, the roles from
// config.Roles, the funnel ports from config.FunnelPublicPorts, and the
// inherited variables from runtime.InheritedEnvironment. A field described here
// and absent there would be a promise nothing keeps.
func ConfigPage(version string, date string) string {
	page := &builder{}

	page.line(".TH GRAT.CONFIG 7 \"" + date + "\" " + quote("grat "+version) + " " + quote("File Formats"))

	page.section("NAME")
	page.line("grat.config \\- the declarative description of a project's services")

	page.section("DESCRIPTION")
	page.paragraphs(configDescription)

	page.section("EXAMPLE")
	writeExample(page)

	page.section("TOP LEVEL")
	writeFields(page, topLevelFields)

	page.section("THE PROJECT TABLE")
	writeFields(page, projectFields)

	page.section("THE RUNTIME TABLE")
	page.paragraphs(runtimeIntro)
	writeRuntimeFields(page)

	page.section("A SERVICE")
	page.paragraphs(serviceIntro)
	writeFields(page, serviceFields)

	page.section("THE EXPOSE TABLE")
	page.paragraphs(exposeIntro)
	writeExposeFields(page)

	page.section("ROLES AND PORT RANGES")
	page.paragraphs(rolesIntro)
	writeRoleRanges(page)

	page.section("THE ENVIRONMENT A COMMAND RECEIVES")
	page.paragraphs(environmentIntro)
	writeEnvironment(page)

	page.section("SEE ALSO")
	page.paragraphs(configSeeAlso)

	return page.String()
}

// field is one key of the file, as the reader meets it.
type field struct {
	name     string
	required string
	meaning  string
}

// writeFields renders a group of keys as a definition list.
func writeFields(page *builder, fields []field) {
	for _, entry := range fields {
		page.line(".TP")
		page.line(".B " + escape(entry.name))
		page.line(escape(entry.required + " " + entry.meaning))
	}
}

// writeRuntimeFields builds the timing table from the defaults themselves, so a
// default changed in the code changes the page with it.
func writeRuntimeFields(page *builder) {
	defaults := config.DefaultRuntime()
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
		page.line(".TP")
		page.line(".B " + escape(entry.name))
		page.line(escape("Optional. " + entry.meaning + " The default is " + entry.value + "."))
	}
}

// writeExposeFields names the two keys of the expose table, with the ports taken
// from the configuration rather than repeated.
func writeExposeFields(page *builder) {
	ports := make([]string, 0, len(config.FunnelPublicPorts()))
	for _, port := range config.FunnelPublicPorts() {
		ports = append(ports, strconv.Itoa(port))
	}

	page.line(".TP")
	page.line(".B path")
	page.line(escape("Required once the table exists. The only path that goes public, beginning with a slash. Everything else the service serves stays on the machine."))
	page.line(".TP")
	page.line(".B public_port")
	page.line(escape("Optional. One of " + strings.Join(ports, ", ") + ", which are the ports a Tailscale funnel listens on. The default is " + strconv.Itoa(config.DefaultPublicPort) + "."))
}

// writeRoleRanges builds the range table from the roles themselves.
func writeRoleRanges(page *builder) {
	for _, role := range config.Roles() {
		portRange, known := role.PortRange()
		page.line(".TP")
		page.line(".B " + escape(string(role)))
		switch {
		case !known:
			page.line("No range is allocated for this role.")
		case portRange.First == 0:
			page.line("No port. The service is watched as a process and is never probed over HTTP, so it takes port 0 and no health path.")
		default:
			page.line(strconv.Itoa(portRange.First) + " to " + strconv.Itoa(portRange.Last))
		}
	}
}

// writeEnvironment lists the baseline the runtime actually passes.
func writeEnvironment(page *builder) {
	page.line(".PP")
	page.line(escape(strings.Join(runtime.InheritedEnvironment(), ", ") + "."))
	page.paragraphs(environmentRest)
	// The variable is named from the constant the runtime sets it under, so the
	// page cannot describe a name the code no longer uses.
	page.line(".PP")
	page.paragraphs(fmt.Sprintf(environmentTailnetHost, runtime.TailnetHostVariable()))
}

// writeExample prints a complete file, indented as a literal block.
func writeExample(page *builder) {
	page.line(".nf")
	page.line(".RS 4")
	for _, line := range strings.Split(strings.TrimSpace(configExample), "\n") {
		if line == "" {
			page.line(".sp")
			continue
		}
		page.line(escape(line))
	}
	page.line(".RE")
	page.line(".fi")
}
