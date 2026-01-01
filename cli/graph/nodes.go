package graph

import (
	"fmt"
	"pkb-agent/graph"
	"pkb-agent/graph/loaders/metaloader"
	pathlib "pkb-agent/util/pathlib"

	"github.com/spf13/cobra"
)

type ListNodesCommand struct {
	cobra.Command
}

func newListNodesCommand() *cobra.Command {
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

	g, err := graph.LoadGraph(path, loader)
	if err != nil {
		return err
	}

	names := g.ListNodeNames()
	for _, name := range names {
		fmt.Println(name)
	}

	return nil
}
