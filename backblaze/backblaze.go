package backblaze

import (
	"context"
	"io"
	"os"

	"github.com/Backblaze/blazer/b2"
)

func NewClient(ctx context.Context) (*b2.Client, error) {
	application_key_id := os.Getenv("APPLICATION_KEY_ID")
	application_key := os.Getenv("APPLICATION_KEY")

	return b2.NewClient(ctx, application_key_id, application_key)
}

func Download(ctx context.Context, client *b2.Client, bucketName string, remoteFilename string, writer io.Writer, concurrentDownloads int) error {
	bucket, err := client.Bucket(ctx, bucketName)
	if err != nil {
		return err
	}

	object := bucket.Object(remoteFilename)
	remoteReader := object.NewReader(ctx)
	defer remoteReader.Close()

	remoteReader.ConcurrentDownloads = concurrentDownloads

	if _, err := io.Copy(writer, remoteReader); err != nil {
		return err
	}

	return nil
}

func DownloadToFile(ctx context.Context, client *b2.Client, bucketName string, remoteFilename string, localFilename string, concurrentDownloads int) error {
	localFile, err := os.Create(localFilename)
	if err != nil {
		return err
	}
	defer localFile.Close()

	if err := Download(ctx, client, bucketName, remoteFilename, localFile, concurrentDownloads); err != nil {
		return err
	}

	return nil
}

func Delete(ctx context.Context, client *b2.Client, bucketName string, remoteFilename string) error {
	bucket, err := client.Bucket(ctx, bucketName)
	if err != nil {
		return err
	}

	object := bucket.Object(remoteFilename)
	if err := object.Delete(ctx); err != nil {
		return err
	}

	return nil
}

func Upload(ctx context.Context, client *b2.Client, bucketName string, reader io.Reader, remoteFilename string) error {
	bucket, err := client.Bucket(ctx, bucketName)
	if err != nil {
		return err
	}

	object := bucket.Object(remoteFilename)
	writer := object.NewWriter(ctx)
	defer writer.Close()

	if _, err := io.Copy(writer, reader); err != nil {
		return err
	}

	return nil
}
