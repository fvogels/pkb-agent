package snippetpage

import (
	"pkb-agent/graph/node"

	tea "github.com/charmbracelet/bubbletea"
)

type Page struct {
	caption  string
	source   string
	language string
}

func New(caption string, source string, language string) *Page {
	return &Page{
		caption:  caption,
		source:   source,
		language: language,
	}
}

func (page *Page) CreateViewer() tea.Model {
	return NewModel(page.source, page.language)
}

func (page *Page) GetCaption() string {
	return page.caption
}

func (page *Page) GetActions() []node.Action {
	return nil
}
