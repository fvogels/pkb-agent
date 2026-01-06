package cli

import (
	"pkb-agent/tui/tuimain"

	"github.com/spf13/cobra"
)

type startUserInterfaceCommand struct {
	CobraCommand cobra.Command
}

func NewStartUserInterfaceCommand() *cobra.Command {
	var command *startUserInterfaceCommand

	command = &startUserInterfaceCommand{
		CobraCommand: cobra.Command{
			Use:   "ui",
			Short: "Start ui",
			Long:  `Start user interface`,
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return command.execute()
			},
		},
	}

	return &command.CobraCommand
}

func (c *startUserInterfaceCommand) execute() error {
	verbose := false
	verboseFlag := c.CobraCommand.InheritedFlags().Lookup("verbose")

	if verboseFlag != nil {
		verbose = verboseFlag.Value.String() == "true"
	}

	return tuimain.Start(verbose)
}
