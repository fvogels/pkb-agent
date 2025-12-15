package markdownview

type MsgSetSource struct {
	Source string
}

type msgRenderingDone struct {
	renderedMarkdown string
}
