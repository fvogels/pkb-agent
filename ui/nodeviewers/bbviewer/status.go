package bbviewer

import "fmt"

type status interface {
	view() string
	getChannel() chan status
}

type statusDownloading struct {
	bytesDownloaded int
	channel         chan status
}

type statusFinished struct{}

type statusErrorOccurred struct {
	err error
}

type statusUnzipping struct {
	channel chan status
}

func (status statusDownloading) view() string {
	return fmt.Sprintf("Downloaded %d bytes", status.bytesDownloaded)
}

func (status statusDownloading) getChannel() chan status {
	return status.channel
}

func (status statusFinished) view() string {
	return "Done"
}

func (status statusFinished) getChannel() chan status {
	return nil
}

func (status statusErrorOccurred) view() string {
	return fmt.Sprintf("An error occurred: %v", status.err)
}

func (status statusErrorOccurred) getChannel() chan status {
	return nil
}

func (status statusUnzipping) view() string {
	return "Unzipping"
}

func (status statusUnzipping) getChannel() chan status {
	return status.channel
}
