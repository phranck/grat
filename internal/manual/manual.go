// Package manual builds grat's manual as a document and renders it.
//
// The manual is written once, as structure, and rendered twice: as a man page
// and as the Markdown a repository reads. Two documents describing one tool drift
// apart, which is what a generated one cannot do.
package manual

import (
	"strconv"
	"strings"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/runtime"
)

// Page renders the command manual as a man page.
func Page(version string, date string, groups []presentation.CommandGroup) string {
	return Roff(CommandDocument(version, date, groups))
}

// CommandDocument builds grat(1).
//
// Every figure in it comes from the code that enforces it: the roles and their
// ranges from config, the inherited environment from runtime, and the command
// list from the same reference grat help prints. A number written out here would
// be a second copy of one the code already holds.
func CommandDocument(version string, date string, groups []presentation.CommandGroup) Document {
	return Document{
		Name:          "grat",
		ManualSection: 1,
		Category:      "User Commands",
		Title:         "start, watch and stop the development services of a project",
		Version:       version,
		Date:          date,
		Sections: []Section{
			{Title: "Name", Blocks: []Block{
				Prose("grat - start, watch and stop the development services of a project"),
			}},
			{Title: "Synopsis", Blocks: []Block{
				Literal("grat [OPTION]... COMMAND [ARGUMENT]..."),
			}},
			{Title: "Description", Blocks: []Block{Prose(description)}},
			{Title: "Installation", Blocks: []Block{
				Prose(installBrew),
				Literal("brew install phranck/grat/grat"),
				Prose(installBinary),
				Literal("gh attestation verify ./grat_VERSION_OS_ARCH \\\n  --repo phranck/grat \\\n  --signer-workflow phranck/grat/.github/workflows/release.yml \\\n  --source-ref refs/tags/VERSION \\\n  --deny-self-hosted-runners"),
				Prose(installPages),
				Literal("sudo install -m 0644 grat.1 /usr/local/share/man/man1/grat.1\nsudo install -m 0644 grat.config.7 /usr/local/share/man/man7/grat.config.7"),
				Prose(installGenerated),
				Prose(installRuntime),
			}},
			commandSection(groups),
			{Title: "How grat decides what a project runs", Blocks: []Block{Prose(detection)}},
			roleSection(),
			environmentSection(),
			{Title: "Readiness and status", Blocks: []Block{
				Prose(readiness),
				Rows([]string{"State", "Meaning"}, [][]string{
					{"stopped", "No live managed process exists for the configured service."},
					{"running", "The managed process passes every check its role calls for."},
					{"unhealthy", "The process is alive whilst its identity, its listener ownership or its health check fails."},
				}),
				Prose(statusColumns),
			}},
			{Title: "Shutdown and restart", Blocks: []Block{Prose(shutdown)}},
			{Title: "Public access", Blocks: []Block{Prose(publicAccess)}},
			{Title: "Ports", Blocks: []Block{Prose(portAllocation)}},
			{Title: "Scan directories", Blocks: []Block{Prose(scanDirectories)}},
			{Title: "Maintenance", Blocks: []Block{Prose(maintenance)}},
			{Title: "Safety", Blocks: []Block{Prose(safety)}},
			{Title: "Files", Blocks: []Block{fileList()}},
			{Title: "Exit status", Blocks: []Block{exitStatusList()}},
			{Title: "See also", Blocks: []Block{Prose(seeAlso)}},
		},
	}
}

// commandSection turns the command reference into one entry per command, with
// the detail behind it.
//
// The reference decides which commands exist and the details decide what is said
// about each. A command in the reference without a detail entry is a command
// that ships undocumented, which is what TestEveryCommandHasAnEntry prevents.
func commandSection(groups []presentation.CommandGroup) Section {
	// The groups open the section as one list, so a reader sees what kinds of
	// command exist before meeting them one at a time. Repeating each group as a
	// heading would put it on the same level as the commands under it, and a
	// reader scanning for a command name would have to tell the two apart.
	overview := make([]Item, 0, len(groups))
	blocks := []Block{Prose(commandsIntro)}
	for _, group := range groups {
		names := make([]string, 0, len(group.Commands))
		for _, command := range group.Commands {
			names = append(names, firstWords(command.Usage))
		}
		overview = append(overview, Item{Term: group.Title, Detail: strings.Join(names, ", ") + "."})
	}
	blocks = append(blocks, Definitions(overview...))

	for _, group := range groups {
		for _, command := range group.Commands {
			blocks = append(blocks, Subheading("grat "+command.Usage))
			blocks = append(blocks, Prose(command.Description+"."))
			detail, found := detailFor(command.Usage)
			if !found {
				continue
			}
			blocks = append(blocks, Prose(detail.detail))
			if len(detail.options) > 0 {
				items := make([]Item, 0, len(detail.options))
				for _, option := range detail.options {
					items = append(items, Item{Term: option.flag, Detail: option.meaning})
				}
				blocks = append(blocks, Definitions(items...))
			}
		}
	}
	return Section{Title: "Commands", Blocks: blocks}
}

// firstWords keeps the part of a usage string a reader would say aloud, so
// "logs [--follow] NAME" reads as "logs" in an overview.
func firstWords(usage string) string {
	words := []string{}
	for _, word := range strings.Fields(usage) {
		if strings.HasPrefix(word, "[") || strings.HasPrefix(word, "-") || strings.ToUpper(word) == word {
			break
		}
		// A usage naming two spellings, such as "version, --version", separates
		// them with a comma that is not part of either name.
		words = append(words, strings.TrimSuffix(word, ","))
	}
	if len(words) == 0 {
		return usage
	}
	return strings.Join(words, " ")
}

// detailFor finds the entry written for one command of the reference.
func detailFor(usage string) (commandDetail, bool) {
	for _, detail := range commandDetails {
		if detail.usage == usage {
			return detail, true
		}
	}
	return commandDetail{}, false
}

// roleSection builds the range table from the roles themselves, so a role added
// to the configuration appears here without anybody remembering to add it.
func roleSection() Section {
	rows := [][]string{}
	for _, role := range config.Roles() {
		portRange, known := role.PortRange()
		switch {
		case !known:
			rows = append(rows, []string{string(role), "no range", roleReadiness(role)})
		case portRange.First == 0:
			rows = append(rows, []string{string(role), "no port", roleReadiness(role)})
		default:
			rows = append(rows, []string{
				string(role),
				strconv.Itoa(portRange.First) + " to " + strconv.Itoa(portRange.Last),
				roleReadiness(role),
			})
		}
	}
	return Section{Title: "Roles and ports", Blocks: []Block{
		Prose(roles),
		Rows([]string{"Role", "Port range", "Readiness"}, rows),
		Prose(roleNaming),
	}}
}

// roleReadiness says what a role has to satisfy to count as ready.
func roleReadiness(role config.Role) string {
	if portRange, known := role.PortRange(); known && portRange.First == 0 {
		return "the managed process is alive"
	}
	return "managed process, owned listener, HTTP 2xx"
}

// environmentSection lists the baseline the runtime actually passes.
func environmentSection() Section {
	names := runtime.InheritedEnvironment()
	items := make([]Item, 0, len(names))
	for _, name := range names {
		items = append(items, Item{Term: name, Detail: "Passed through where the parent has it."})
	}
	return Section{Title: "The environment a command receives", Blocks: []Block{
		Prose(commandEnvironmentIntro),
		Definitions(items...),
		Prose(commandEnvironmentRest),
		Definitions(
			Item{Term: "PORT", Detail: "The port grat assigned, for a service that has one."},
			Item{Term: "BACKEND_URL", Detail: backendURLMeaning},
			Item{Term: runtime.TailnetHostVariable(), Detail: tailnetHostMeaning},
		),
	}}
}

// fileList names what grat reads and writes.
func fileList() Block {
	items := make([]Item, 0, len(files))
	for _, entry := range files {
		items = append(items, Item{Term: entry.path, Detail: entry.meaning, Emphasised: true})
	}
	return Definitions(items...)
}

// exitStatusList names what a status code means to a script calling grat.
func exitStatusList() Block {
	items := make([]Item, 0, len(exitStatus))
	for _, entry := range exitStatus {
		items = append(items, Item{Term: entry.code, Detail: entry.meaning})
	}
	return Definitions(items...)
}
