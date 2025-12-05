package b2

import "github.com/spf13/cobra"

func NewB2Command() *cobra.Command {
	command := cobra.Command{
		Use:   "b2",
		Short: "Backblaze related commands",
	}

	command.AddCommand(NewListBucketsCommand())
	command.AddCommand(NewListFilesCommand())
	command.AddCommand(NewDownloadCommand())
	command.AddCommand(NewRenameCommand())

	return &command
}
