package markdownpage

import (
	"pkb-agent/graph/node"

	tea "github.com/charmbracelet/bubbletea"
)

type Page struct {
	caption string
	source  string
}

func New(caption string, source string) *Page {
	return &Page{
		caption: caption,
		source:  source,
	}
}

func (page *Page) CreateViewer() tea.Model {
	return NewModel(page.source)
}

func (page *Page) GetCaption() string {
	return page.caption
}

func (page *Page) GetActions() []node.Action {
	return nil
}
