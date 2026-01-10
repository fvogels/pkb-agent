package cli

import (
	"fmt"
	"log/slog"
	"os"
	"pkb-agent/cli/b2"
	"pkb-agent/cli/graph"
	"runtime/pprof"

	"github.com/spf13/cobra"
)

const (
	profilingFilename = "pkb-agent.prof"
)

func RunCLI() {
	verbose := false
	profile := false

	rootCommand := cobra.Command{}
	rootCommand.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCommand.PersistentFlags().BoolVarP(&profile, "profile", "p", false, "Enable profiling")

	cobra.OnInitialize(func() {
		if verbose {
			slog.SetLogLoggerLevel(slog.LevelDebug)
			slog.Info("Verbose mode enabled")
		}

		if profile {
			out, _ := os.Create(profilingFilename)
			pprof.StartCPUProfile(out)
		}
	})

	cobra.OnFinalize(func() {
		pprof.StopCPUProfile()
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
