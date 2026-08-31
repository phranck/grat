package cli

import (
	"fmt"
	"strings"

	"github.com/phranck/grat/internal/presentation"
)

func isHelp(argument string) bool {
	return argument == "help" || argument == "-h" || argument == "--help"
}

func hasHelpFlag(arguments []string) bool {
	for _, argument := range arguments {
		if argument == "-h" || argument == "--help" {
			return true
		}
	}
	return false
}

type globalOptions struct {
	color presentation.ColorMode
}

func parseGlobalOptions(args []string) (globalOptions, []string, error) {
	options := globalOptions{color: presentation.ColorAuto}
	for len(args) > 0 {
		switch {
		case args[0] == "--version":
			return options, []string{"version"}, nil
		case args[0] == "--no-color":
			options.color = presentation.ColorNever
			args = args[1:]
		case args[0] == "--color":
			if len(args) < 2 {
				return globalOptions{}, nil, fmt.Errorf("--color requires auto, always, or never")
			}
			mode, err := presentation.ParseColorMode(args[1])
			if err != nil {
				return globalOptions{}, nil, err
			}
			options.color = mode
			args = args[2:]
		case strings.HasPrefix(args[0], "--color="):
			mode, err := presentation.ParseColorMode(strings.TrimPrefix(args[0], "--color="))
			if err != nil {
				return globalOptions{}, nil, err
			}
			options.color = mode
			args = args[1:]
		default:
			return options, args, nil
		}
	}
	return options, args, nil
}
