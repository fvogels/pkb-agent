package b2

import (
	"context"
	"fmt"
	"os"
	"pkb-agent/backblaze"

	"github.com/spf13/cobra"
)

type listBucketsCommand struct {
	cobra.Command
}

func NewListBucketsCommand() *cobra.Command {
	var command *listBucketsCommand

	command = &listBucketsCommand{
		Command: cobra.Command{
			Use:   "list-buckets",
			Short: "List buckets",
			RunE: func(cmd *cobra.Command, args []string) error {
				return command.execute()
			},
		},
	}

	return &command.Command
}

func (c *listBucketsCommand) execute() error {
	ctx := context.Background()
	application_key := os.Getenv("APPLICATION_KEY")
	application_key_id := os.Getenv("APPLICATION_KEY_ID")

	client, err := backblaze.New(ctx, application_key, application_key_id)

	if err != nil {
		return fmt.Errorf("failed to create b2 client: %w", err)
	}

	buckets, err := client.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch b2 bucket list: %w", err)
	}

	for _, bucket := range buckets {
		fmt.Println(bucket)
	}

	return nil
}
