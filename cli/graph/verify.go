package graph

import (
	"errors"
	"fmt"
	"pkb-agent/pkg"
	"pkb-agent/pkg/loaders/metaloader"
	pathlib "pkb-agent/util/pathlib"

	"github.com/spf13/cobra"
)

type verifyGraphCommand struct {
	cobra.Command
}

var ErrCycle = errors.New("found cycle")
var ErrRedundantLinks = errors.New("found redundant links")

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

	g, err := pkg.LoadGraph(path, loader)
	if err != nil {
		return err
	}

	if c.lookForCycles(g) {
		return ErrCycle
	}

	// Only safe to do when we're sure there are no cycles in the graph
	if c.lookForDuplicateLinks(g) {
		return ErrRedundantLinks
	}

	return nil
}

func (c *verifyGraphCommand) lookForCycles(graph *pkg.Graph) bool {
	if pkg.ContainsCycles(graph) {
		fmt.Println("Error: cycle detected")
		return true
	}

	return false
}

func (c *verifyGraphCommand) lookForDuplicateLinks(graph *pkg.Graph) bool {
	result := false

	for _, node := range graph.ListNodes() {
		duplicates := graph.FindRedundantLinks(node)

		if duplicates.Size() > 0 {
			fmt.Printf("Node \"%s\" has redundant links:\n", node.GetName())

			for _, link := range duplicates.ToSlice() {
				fmt.Printf("  %s\n", graph.FindNodeByIndex(link).GetName())
			}

			result = true
		}
	}

	return result
}
