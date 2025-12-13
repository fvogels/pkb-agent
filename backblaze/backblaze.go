package backblaze

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"pkb-agent/util"

	"github.com/Backblaze/blazer/b2"
)

func NewClient(ctx context.Context) (*b2.Client, error) {
	application_key_id := os.Getenv("APPLICATION_KEY_ID")
	application_key := os.Getenv("APPLICATION_KEY")

	return b2.NewClient(ctx, application_key_id, application_key)
}

type BackblazeClient struct {
	b2client *b2.Client
}

func New(ctx context.Context, application_key string, application_key_id string) (*BackblazeClient, error) {
	b2client, err := b2.NewClient(ctx, application_key_id, application_key)
	if err != nil {
		return nil, fmt.Errorf("failed to create BackblazeClient: %w", err)
	}

	client := BackblazeClient{
		b2client: b2client,
	}

	return &client, nil
}

func (client *BackblazeClient) Download(ctx context.Context, bucketName string, remoteFilename string, writer io.Writer, concurrentDownloads int, callback func(progress int)) error {
	bucket, err := client.b2client.Bucket(ctx, bucketName)
	if err != nil {
		return err
	}

	object := bucket.Object(remoteFilename)
	remoteReader := object.NewReader(ctx)
	defer remoteReader.Close()

	remoteReader.ConcurrentDownloads = concurrentDownloads

	observableReader := observableReader{
		observedReader: remoteReader,
		totalBytesRead: 0,
		callback:       callback,
	}

	if _, err := io.Copy(writer, &observableReader); err != nil {
		return err
	}

	return nil
}

func (client *BackblazeClient) DownloadToBuffer(ctx context.Context, bucketName string, remoteFilename string, concurrentDownloads int, callback func(progress int)) ([]byte, error) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	if err := client.Download(ctx, bucketName, remoteFilename, writer, 2, callback); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

type observableReader struct {
	observedReader io.Reader
	totalBytesRead int
	callback       func(totalBytesRead int)
}

func (reader *observableReader) Read(buffer []byte) (int, error) {
	byte_count, err := reader.observedReader.Read(buffer)
	reader.totalBytesRead += byte_count
	reader.callback(reader.totalBytesRead)

	return byte_count, err
}

func (client *BackblazeClient) DownloadToFile(ctx context.Context, bucketName string, remoteFilename string, localFilename string, concurrentDownloads int, callback func(progress int)) error {
	localFile, err := os.Create(localFilename)
	if err != nil {
		return err
	}
	defer localFile.Close()

	if err := client.Download(ctx, bucketName, remoteFilename, localFile, concurrentDownloads, callback); err != nil {
		return err
	}

	return nil
}

func (client *BackblazeClient) Delete(ctx context.Context, bucketName string, remoteFilename string) error {
	bucket, err := client.b2client.Bucket(ctx, bucketName)
	if err != nil {
		return err
	}

	object := bucket.Object(remoteFilename)
	if err := object.Delete(ctx); err != nil {
		return err
	}

	return nil
}

func (client *BackblazeClient) Upload(ctx context.Context, bucketName string, reader io.Reader, remoteFilename string) error {
	bucket, err := client.b2client.Bucket(ctx, bucketName)
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

func (client *BackblazeClient) ListBuckets(ctx context.Context) ([]string, error) {
	buckets, err := client.b2client.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	bucketNames := util.Map(buckets, func(b *b2.Bucket) string { return b.Name() })

	return bucketNames, nil
}
