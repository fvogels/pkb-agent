package horizontal

import (
	"pkb-agent/ui/layout"
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Layout[T any] struct {
	children []child[T]
}

type child[T any] struct {
	determineWidth func(size util.Size) int
	computedWidth  int
	component      layout.Layout[T]
}

func New[T any]() Layout[T] {
	return Layout[T]{
		children: nil,
	}
}

func (layout *Layout[T]) LayoutResize(parent *T, size util.Size) tea.Cmd {
	commands := []tea.Cmd{}

	for index := range layout.children {
		child := &layout.children[index]
		child.computedWidth = child.determineWidth(size)

		command := child.component.LayoutResize(parent, util.Size{
			Width:  child.computedWidth,
			Height: size.Height,
		})

		commands = append(commands, command)
	}

	return tea.Batch(commands...)
}

func (layout *Layout[T]) LayoutView(model *T) string {
	parts := []string{}

	for _, child := range layout.children {
		style := lipgloss.NewStyle().Width(child.computedWidth)
		part := style.Render(child.component.LayoutView(model))
		parts = append(parts, part)
	}

	return lipgloss.JoinHorizontal(0, parts...)
}

func (layout *Layout[T]) Add(determineWidth func(size util.Size) int, component layout.Layout[T]) {
	layout.children = append(layout.children, child[T]{
		determineWidth: determineWidth,
		component:      component,
	})
}
