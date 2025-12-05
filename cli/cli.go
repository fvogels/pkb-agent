package cli

import (
	"pkb-agent/cli/b2"

	"github.com/spf13/cobra"
)

func RunCLI() {
	rootCommand := cobra.Command{}

	rootCommand.AddCommand(b2.NewB2Command())

	rootCommand.Execute()
}
