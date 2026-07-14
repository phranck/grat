// Command grat manages configured local development processes from any project.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/phranck/grat/internal/cli"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	os.Exit(cli.Run(ctx, os.Args[1:], mustGetwd(), os.Stdout, os.Stderr))
}

func mustGetwd() string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return "."
	}
	return workingDirectory
}
