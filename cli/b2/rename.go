package b2

import (
	"context"
	"fmt"
	"pkb-agent/backblaze"

	"github.com/spf13/cobra"
)

type RenameCommand struct {
	cobra.Command
	bucketName  string
	oldFilename string
	newFilename string
}

func NewRenameCommand() *cobra.Command {
	var command *RenameCommand

	command = &RenameCommand{
		Command: cobra.Command{
			Use:   "rename",
			Short: "Renames remote file",
			RunE: func(cmd *cobra.Command, args []string) error {
				command.bucketName = args[0]
				command.oldFilename = args[1]
				command.newFilename = args[2]

				return command.execute()
			},
			Args: cobra.ExactArgs(3),
		},
	}

	return &command.Command
}

func (c *RenameCommand) execute() error {
	ctx := context.Background()

	client, err := backblaze.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("failed to create b2 client: %w", err)
	}

	bucket, err := client.Bucket(ctx, c.bucketName)
	if err != nil {
		return fmt.Errorf("failed to fetch b2 bucket: %w", err)
	}

	object := bucket.Object(c.oldFilename)
	if object == nil {
		return fmt.Errorf("nil object")
	}
	if _, err := object.Attrs(ctx); err != nil {
		return err
	}

	fileId := object.ID()

	fmt.Printf(`b2 file copy-by-id %s %s %s && b2 rm b2://%s/%s`, fileId, c.bucketName, c.newFilename, c.bucketName, c.oldFilename)

	return nil
}
