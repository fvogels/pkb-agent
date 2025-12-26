package graph

import (
	"fmt"
	"pkb-agent/graph"
	"pkb-agent/graph/metaloader"
	pathlib "pkb-agent/util/pathlib"

	"github.com/spf13/cobra"
)

type verifyGraphCommand struct {
	cobra.Command
}

func newVerifyGraphCommand() *cobra.Command {
	var command *verifyGraphCommand

	command = &verifyGraphCommand{
		Command: cobra.Command{
			Use:   "verify",
			Short: "Checks the graph",
			RunE: func(cmd *cobra.Command, args []string) error {
				return command.execute()
			},
			Args: cobra.NoArgs,
		},
	}

	return &command.Command
}

func (c *verifyGraphCommand) execute() error {
	loader := metaloader.New()
	path := pathlib.New(`F:\repos\pkb\pkb-data\root.yaml`)

	g, err := graph.LoadGraph(path, loader)
	if err != nil {
		return err
	}

	if graph.ContainsCycles(g) {
		fmt.Println("Error: cycle detected")
		return nil
	}

	return nil
}
