package sourceviewer

type MsgSetSource struct {
	Source   string
	Language string
}

type msgSourceFormatted struct {
	recipient       int
	formattedSource string
}
