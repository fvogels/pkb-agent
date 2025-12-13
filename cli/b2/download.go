package b2

import (
	"context"
	"fmt"
	"os"
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
	application_key := os.Getenv("APPLICATION_KEY")
	application_key_id := os.Getenv("APPLICATION_KEY_ID")

	client, err := backblaze.New(ctx, application_key, application_key_id)
	if err != nil {
		return err
	}

	channel := make(chan int)

	go func() {
		client.DownloadToFile(ctx, c.bucketName, c.remoteFilename, c.localFilename, c.concurrentDownloads, func(n int) {
			channel <- n
		})

		channel <- -1
	}()

	value := <-channel
	for value != -1 {
		fmt.Printf("Downloaded %d bytes\n", value)
		value = <-channel
	}

	close(channel)

	return nil
}
