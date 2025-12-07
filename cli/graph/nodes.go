package graph

import (
	"fmt"
	"pkb-agent/graph"
	"pkb-agent/graph/loaders/atomloader"
	"pkb-agent/graph/loaders/metaloader"
	"pkb-agent/graph/loaders/snippetloader"
	pathlib "pkb-agent/util/path"

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
	loader.AddLoader("atom", atomloader.New())
	loader.AddLoader("snippet", snippetloader.New())

	path := pathlib.New(`F:\repos\pkb\pkb-data\root.yaml`)
	callback := func(entry graph.Node) error {
		fmt.Printf("%s\n", entry.GetName())
		return nil
	}

	if err := loader.Load(path, callback); err != nil {
		return err
	}

	return nil
}
