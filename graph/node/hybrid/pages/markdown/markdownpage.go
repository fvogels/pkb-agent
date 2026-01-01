package markdownpage

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Page struct {
	source string
}

func New(source string) *Page {
	return &Page{
		source: source,
	}
}

func (page *Page) CreateViewer() tea.Model {
	return NewModel(page.source)
}
