package b2

import (
	"context"
	"fmt"
	"pkb-agent/backblaze"

	"github.com/spf13/cobra"
)

type ListBucketsCommand struct {
	cobra.Command
}

func NewListBucketsCommand() *cobra.Command {
	var command *ListBucketsCommand

	command = &ListBucketsCommand{
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

func (c *ListBucketsCommand) execute() error {
	ctx := context.Background()
	client, err := backblaze.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("failed to create b2 client: %w", err)
	}

	buckets, err := client.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch b2 bucket list: %w", err)
	}

	for _, bucket := range buckets {
		fmt.Println(bucket.Name())
	}

	return nil
}
