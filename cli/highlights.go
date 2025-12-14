package cli

import (
	"fmt"
	"slices"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/spf13/cobra"
)

type listHighlightableLanguagesCommand struct {
	CobraCommand cobra.Command
}

func NewListHighlightableLanguagesCommand() *cobra.Command {
	var command *listHighlightableLanguagesCommand

	command = &listHighlightableLanguagesCommand{
		CobraCommand: cobra.Command{
			Use:   "highlights",
			Short: "Supported highlighters",
			Long:  `Lists all languages for which there is support for syntax highlighting`,
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return command.execute()
			},
		},
	}

	return &command.CobraCommand
}

func (c *listHighlightableLanguagesCommand) execute() error {
	languages := lexers.Names(true)
	slices.Sort(languages)

	for _, language := range languages {
		fmt.Println(language)
	}

	return nil
}
