package border

import (
	"pkb-agent/ui/layout"
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Layout[T any] struct {
	child       layout.Layout[T]
	size        util.Size
	borderStyle lipgloss.Border
}

func New[T any](child layout.Layout[T]) *Layout[T] {
	return &Layout[T]{
		child:       child,
		borderStyle: lipgloss.DoubleBorder(),
	}
}

func (layout *Layout[T]) LayoutResize(parent *T, size util.Size) tea.Cmd {
	layout.size = size

	childSize := util.Size{
		Width:  size.Width - 2,
		Height: size.Height - 2,
	}

	return layout.child.LayoutResize(parent, childSize)
}

func (layout *Layout[T]) LayoutView(model *T) string {
	style := lipgloss.NewStyle().Border(layout.borderStyle).Width(layout.size.Width - 2).Height(layout.size.Height - 2)
	return style.Render(layout.child.LayoutView(model))
}
