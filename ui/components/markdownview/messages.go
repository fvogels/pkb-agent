package markdownview

type MsgSetSource struct {
	Source string
}

type msgRenderingDone struct {
	recipient        int
	renderedMarkdown string
}
