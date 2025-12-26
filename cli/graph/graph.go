package graph

import "github.com/spf13/cobra"

func NewGraphCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "graph",
		Short: "Graph related commands",
	}

	command.AddCommand(newListNodesCommand())
	command.AddCommand(newSearchGraphCommand())
	command.AddCommand(newVerifyGraphCommand())

	return &command
}
