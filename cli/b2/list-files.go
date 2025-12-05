package b2

import (
	"context"
	"fmt"
	"pkb-agent/backblaze"

	"github.com/spf13/cobra"
)

type ListFilesCommand struct {
	cobra.Command
	bucketName string
}

func NewListFilesCommand() *cobra.Command {
	var command *ListFilesCommand

	command = &ListFilesCommand{
		Command: cobra.Command{
			Use:   "list-files",
			Short: "List files",
			RunE: func(cmd *cobra.Command, args []string) error {
				command.bucketName = args[0]

				return command.execute()
			},
			Args: cobra.ExactArgs(1),
		},
	}

	return &command.Command
}

func (c *ListFilesCommand) execute() error {
	ctx := context.Background()

	client, err := backblaze.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("failed to create b2 client: %w", err)
	}

	bucket, err := client.Bucket(ctx, c.bucketName)
	if err != nil {
		return fmt.Errorf("failed to fetch b2 bucket: %w", err)
	}

	iterator := bucket.List(ctx)

	for iterator.Next() {
		fmt.Println(iterator.Object().Name())
	}

	return nil
}
