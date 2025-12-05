package b2

import (
	"context"
	"pkb-agent/backblaze"

	"github.com/spf13/cobra"
)

type DownloadCommand struct {
	cobra.Command
	bucketName          string
	remoteFilename      string
	localFilename       string
	concurrentDownloads int
}

func NewDownloadCommand() *cobra.Command {
	var command *DownloadCommand

	command = &DownloadCommand{
		Command: cobra.Command{
			Use:   "download",
			Short: "Download file from bucket",
			RunE: func(cmd *cobra.Command, args []string) error {
				command.bucketName = args[0]
				command.remoteFilename = args[1]

				if command.localFilename == "" {
					command.localFilename = command.remoteFilename
				}

				return command.execute()
			},
			Args: cobra.ExactArgs(2),
		},
	}

	command.Flags().StringVarP(&command.localFilename, "local", "o", "", "Local filename")
	command.Flags().IntVarP(&command.concurrentDownloads, "concurrent", "c", 2, "Number of concurrent downloads")

	return &command.Command
}

func (c *DownloadCommand) execute() error {
	ctx := context.Background()

	client, err := backblaze.NewClient(ctx)
	if err != nil {
		return err
	}

	backblaze.DownloadToFile(ctx, client, c.bucketName, c.remoteFilename, c.localFilename, c.concurrentDownloads)

	return nil
}
