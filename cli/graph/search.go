package graph

import (
	"fmt"
	"pkb-agent/graph"
	"pkb-agent/graph/metaloader"
	pathlib "pkb-agent/util/pathlib"

	"github.com/spf13/cobra"
)

type SearchGraphCommand struct {
	cobra.Command
}

func NewSearchGraphCommand() *cobra.Command {
	var command *SearchGraphCommand

	command = &SearchGraphCommand{
		Command: cobra.Command{
			Use:   "search",
			Short: "Searches nodes matching a string",
			RunE: func(cmd *cobra.Command, args []string) error {
				str := args[0]
				return command.execute(str)
			},
			Args: cobra.ExactArgs(1),
		},
	}

	return &command.Command
}

func (c *SearchGraphCommand) execute(str string) error {
	loader := metaloader.New()
	path := pathlib.New(`F:\repos\pkb\pkb-data\root.yaml`)

	g, err := graph.LoadGraph(path, loader)
	if err != nil {
		return err
	}

	iterator := g.FindMatchingNodes(str)
	for iterator.Current() != nil {
		fmt.Println(iterator.Current().Name)
		iterator.Next()
	}

	return nil
}
