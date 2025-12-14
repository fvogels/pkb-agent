package header

import (
	"pkb-agent/ui/layout"
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Layout[T any] struct {
	header string
	child  layout.Layout[T]
	size   util.Size
}

func New[T any](header string, child layout.Layout[T]) *Layout[T] {
	return &Layout[T]{
		header: header,
		child:  child,
	}
}

func (layout *Layout[T]) LayoutResize(parent *T, size util.Size) tea.Cmd {
	layout.size = size

	command := layout.child.LayoutResize(
		parent,
		util.Size{
			Width:  size.Width,
			Height: size.Height - 1,
		},
	)

	return command
}

func (layout *Layout[T]) LayoutView(model *T) string {
	header := lipgloss.NewStyle().Width(layout.size.Width).Render(layout.header)
	child := layout.child.LayoutView(model)

	return lipgloss.JoinVertical(0, header, child)
}
