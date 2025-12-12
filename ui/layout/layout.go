package layout

import (
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
)

type Layout[T any] interface {
	Resize(model *T, size util.Size) tea.Cmd
	View(model *T) string
}

type Child[T any] interface {
	LayoutUpdate(parent *T, size util.Size) tea.Cmd
	LayoutView(parent *T) string
}

type Component[T any] interface {
	TypedUpdate(message tea.Msg) (T, tea.Cmd)
	View() string
}

func Wrap[M any, C Component[C]](get func(*M) *C) Child[M] {
	return wrapper[M, C]{
		get: get,
	}
}

type wrapper[M any, C Component[C]] struct {
	get func(*M) *C
}

func (w wrapper[M, C]) LayoutUpdate(parent *M, size util.Size) tea.Cmd {
	component := w.get(parent)
	message := tea.WindowSizeMsg{
		Width:  size.Width,
		Height: size.Height,
	}
	updatedComponent, command := (*component).TypedUpdate(message)
	*component = updatedComponent
	return command
}

func (w wrapper[M, C]) LayoutView(parent *M) string {
	return (*w.get(parent)).View()
}
