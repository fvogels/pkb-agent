package cli

import (
	"log/slog"
	"pkb-agent/cli/b2"
	"pkb-agent/cli/graph"

	"github.com/spf13/cobra"
)

func RunCLI() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Info("Verbose mode enabled")

	rootCommand := cobra.Command{}

	rootCommand.AddCommand(b2.NewB2Command())
	rootCommand.AddCommand(graph.NewGraphCommand())

	rootCommand.Execute()
}
