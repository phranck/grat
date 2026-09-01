package cli

import (
	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/version"
)

func printUsage(output presentation.Renderer) {
	output.Help(version.Current(), helpCommandGroups())
}

func helpCommandGroups() []presentation.CommandGroup {
	return []presentation.CommandGroup{
		{
			Title: "Project setup",
			Commands: []presentation.Command{
				{Usage: "discover [PATH]", Description: "Write a grat.config for this project, or choose from the projects below PATH"},
			},
		},
		{
			Title: "Service lifecycle",
			Commands: []presentation.Command{
				{Usage: "start [name...]", Description: "Start services and wait for configured readiness"},
				{Usage: "stop [name...]", Description: "Gracefully stop managed service processes"},
				{Usage: "restart [name...]", Description: "Stop, start, and verify selected services"},
				{Usage: "recover [--yes] [name...]", Description: "Preview and recover legacy managed processes"},
				{Usage: "status", Description: "Show managed process and health status"},
				{Usage: "logs [--follow] NAME", Description: "Print or follow a service log"},
			},
		},
		{
			Title: "Public access",
			Commands: []presentation.Command{
				{Usage: "expose [--path P] NAME", Description: "Publish a service to the internet; --path narrows it to one path"},
				{Usage: "expose status [name...]", Description: "Show what is published, with the public address"},
				{Usage: "hide [--path P] NAME", Description: "Withdraw a published service or path"},
			},
		},
		{
			Title: "Ports",
			Commands: []presentation.Command{
				{Usage: "ports audit", Description: "Find configured port collisions and live listeners"},
				{Usage: "ports assign [name...]", Description: "Assign free role-compatible ports"},
				{Usage: "ports reassign", Description: "Stop managed services and globally reassign ports"},
			},
		},
		{
			Title: "Directories",
			Commands: []presentation.Command{
				{Usage: "directories add PATH", Description: "Add a directory to scan for grat.config files"},
				{Usage: "directories remove PATH", Description: "Stop scanning a configured directory"},
				{Usage: "directories list", Description: "List configured scan directories; dir is an alias"},
			},
		},
		{
			Title: "Maintenance",
			Commands: []presentation.Command{
				{Usage: "update", Description: "Update grat according to its installation method"},
				{Usage: "uninstall", Description: "Remove grat, and Tailscale where grat installed it; asks for your password"},
			},
		},
		{
			Title: "Global options",
			Commands: []presentation.Command{
				{Usage: "version, --version", Description: "Print the installed grat version"},
				{Usage: "--color=MODE", Description: "Use auto, always, or never for terminal color"},
				{Usage: "--no-color", Description: "Disable terminal color explicitly"},
				{Usage: "help, --help", Description: "Show this command reference"},
			},
		},
	}
}
