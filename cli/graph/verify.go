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

	if c.lookForCycles(g) {
		return nil
	}

	// Only safe to do when we're sure there are no cycles in the graph
	c.lookForDuplicateLinks(g)

	return nil
}

func (c *verifyGraphCommand) lookForCycles(g *graph.Graph) bool {
	if graph.ContainsCycles(g) {
		fmt.Println("Error: cycle detected")
		return true
	}

	return false
}

func (c *verifyGraphCommand) lookForDuplicateLinks(g *graph.Graph) {
	for _, node := range g.ListNodes() {
		duplicates := g.FindRedundantLinks(node)

		if duplicates.Size() > 0 {
			fmt.Printf("Node \"%s\" has redundant links:\n", node.Name)

			for _, link := range duplicates.ToSlice() {
				fmt.Printf("  %s\n", g.FindNodeByIndex(link).Name)
			}
		}
	}
}
