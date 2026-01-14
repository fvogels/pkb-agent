package markdown

import "github.com/charmbracelet/glamour"

func Render(source string, lineSize int) (string, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(lineSize),
	)
	if err != nil {
		panic("failed to create markdown renderer")
	}
	renderedMarkdown, err := renderer.Render(source)
	if err != nil {
		panic("failed to render markdown file")
	}

	return renderedMarkdown, nil
}
