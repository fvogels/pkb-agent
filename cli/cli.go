package cli

import (
	"fmt"
	"log/slog"
	"os"
	"pkb-agent/cli/b2"
	"pkb-agent/cli/graph"

	"github.com/spf13/cobra"
)

func RunCLI() {
	verbose := false

	rootCommand := cobra.Command{}
	rootCommand.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	cobra.OnInitialize(func() {
		if verbose {
			slog.SetLogLoggerLevel(slog.LevelDebug)
			slog.Info("Verbose mode enabled")
		}
	})

	rootCommand.AddCommand(b2.NewB2Command())
	rootCommand.AddCommand(graph.NewGraphCommand())
	rootCommand.AddCommand(NewStartUserInterfaceCommand())
	rootCommand.AddCommand(NewListHighlightableLanguagesCommand())

	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
