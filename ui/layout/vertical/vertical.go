package vertical

import (
	"pkb-agent/ui/layout"
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type VerticalLayout[T any] struct {
	children []child[T]
}

type child[T any] struct {
	determineHeight func(size util.Size) int
	computedHeight  int
	component       layout.Child[T]
}

func New[T any]() VerticalLayout[T] {
	return VerticalLayout[T]{
		children: nil,
	}
}

func (layout *VerticalLayout[T]) Resize(model *T, size util.Size) tea.Cmd {
	commands := []tea.Cmd{}

	for index := range layout.children {
		child := &layout.children[index]
		child.computedHeight = child.determineHeight(size)

		command := child.component.LayoutUpdate(model, util.Size{
			Width:  size.Width,
			Height: child.computedHeight,
		})

		commands = append(commands, command)
	}

	return tea.Batch(commands...)
}

func (layout *VerticalLayout[T]) View(model *T) string {
	parts := []string{}

	for _, child := range layout.children {
		style := lipgloss.NewStyle().Height(child.computedHeight)
		part := style.Render(child.component.LayoutView(model))
		parts = append(parts, part)
	}

	return lipgloss.JoinVertical(0, parts...)
}

func (layout *VerticalLayout[T]) Add(determineHeight func(size util.Size) int, component layout.Child[T]) {
	layout.children = append(layout.children, child[T]{
		determineHeight: determineHeight,
		component:       component,
	})
}
