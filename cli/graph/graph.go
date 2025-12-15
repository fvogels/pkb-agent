package graph

import "github.com/spf13/cobra"

func NewGraphCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "graph",
		Short: "Graph related commands",
	}

	command.AddCommand(NewListNodesCommand())
	command.AddCommand(NewSearchGraphCommand())

	return &command
}
