package vertical

import (
	"pkb-agent/ui/layout"
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Layout[T any] struct {
	children []child[T]
	size     util.Size
}

type child[T any] struct {
	determineHeight func(size util.Size) int
	computedHeight  int
	component       layout.Layout[T]
}

func New[T any]() Layout[T] {
	return Layout[T]{
		children: nil,
	}
}

func (layout *Layout[T]) LayoutResize(parent *T, size util.Size) tea.Cmd {
	layout.size = size
	commands := []tea.Cmd{}

	for index := range layout.children {
		child := &layout.children[index]
		child.computedHeight = child.determineHeight(size)

		command := child.component.LayoutResize(parent, util.Size{
			Width:  size.Width,
			Height: child.computedHeight,
		})

		commands = append(commands, command)
	}

	return tea.Batch(commands...)
}

func (layout *Layout[T]) LayoutView(model *T) string {
	parts := []string{}

	for _, child := range layout.children {
		style := lipgloss.NewStyle().Width(layout.size.Width).Height(child.computedHeight)
		part := style.Render(child.component.LayoutView(model))
		parts = append(parts, part)
	}

	return lipgloss.JoinVertical(0, parts...)
}

func (layout *Layout[T]) Add(determineHeight func(size util.Size) int, component layout.Layout[T]) {
	layout.children = append(layout.children, child[T]{
		determineHeight: determineHeight,
		component:       component,
	})
}
