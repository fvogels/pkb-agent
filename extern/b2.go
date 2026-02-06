package extern

import (
	"context"
	"log/slog"
	"os"
	"pkb-agent/backblaze"
	"pkb-agent/util/pathlib"
	"pkb-agent/util/zipfile"
)

func BackblazeDownloadAndOpen(bucketName string, filename string) error {
	application_key := os.Getenv("APPLICATION_KEY")
	application_key_id := os.Getenv("APPLICATION_KEY_ID")
	ctx := context.Background()

	client, err := backblaze.New(ctx, application_key, application_key_id)
	if err != nil {
		slog.Error("Failed to create backblaze client")
		return err
	}

	buffer, err := client.DownloadToBuffer(ctx, bucketName, filename, 1, func(bytesDownloaded int) {
		slog.Debug(
			"Received bytes",
			slog.String("bucketName", bucketName),
			slog.String("filename", filename),
			slog.Int("bytesDownloaded", bytesDownloaded),
		)
	})
	if err != nil {
		slog.Error(
			"Failed to download file from Backblaze",
			slog.String("bucket", bucketName),
			slog.String("filename", filename),
		)
		return err
	}

	slog.Debug(
		"Unpacking zipfile",
		slog.String("bucket", bucketName),
		slog.String("filename", filename),
	)
	zippedFiles, err := zipfile.Unpack(buffer)
	if err != nil {
		slog.Error("Failed to unzip files")
		return err
	}

	slog.Debug(
		"Writing file to disk",
		slog.String("bucket", bucketName),
		slog.String("filename", filename),
	)
	zippedFile := zippedFiles[0]
	path, err := zippedFile.SaveToDirectory(pathlib.New("."))
	if err != nil {
		slog.Error("Failed to save downloaded file")
		return err
	}

	slog.Debug(
		"Opening file with default viewer",
		slog.String("bucket", bucketName),
		slog.String("filename", filename),
	)
	if err := OpenUsingDefaultViewer(path); err != nil {
		slog.Error("Failed to open downloaded file using default viewer")
		return err
	}

	return nil
}
