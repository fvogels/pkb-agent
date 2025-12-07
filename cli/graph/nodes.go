package graph

import (
	"pkb-agent/graph"
	"pkb-agent/graph/metaloader"
	pathlib "pkb-agent/util/pathlib"

	"github.com/spf13/cobra"
)

type ListNodesCommand struct {
	cobra.Command
}

func NewListNodesCommand() *cobra.Command {
	var command *ListNodesCommand

	command = &ListNodesCommand{
		Command: cobra.Command{
			Use:   "nodes",
			Short: "Prints out all nodes",
			RunE: func(cmd *cobra.Command, args []string) error {
				return command.execute()
			},
		},
	}

	return &command.Command
}

func (c *ListNodesCommand) execute() error {
	loader := metaloader.New()
	path := pathlib.New(`F:\repos\pkb\pkb-data\root.yaml`)

	if err := graph.LoadGraph(path, loader); err != nil {
		return err
	}

	return nil
}
