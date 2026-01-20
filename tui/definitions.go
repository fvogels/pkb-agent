package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v3"
)

const (
	SafeMode = true
)

type Component interface {
	GetIdentifier() int
	Handle(Message)
	Render() Grid
}

type MeasurableComponent interface {
	Component
	Measure() Size
}

type ComponentBase struct {
	Identifier   int
	Name         string
	MessageQueue MessageQueue
	Size         Size
}

func (base *ComponentBase) GetIdentifier() int {
	return base.Identifier
}

type Grid interface {
	GetSize() Size
	Get(position Position) Cell
}

type Style = tcell.Style

type Cell struct {
	Contents rune
	Style    *Style
	OnClick  func()
}

type Position struct {
	X int
	Y int
}

func (position Position) String() string {
	return fmt.Sprintf("(%d, %d)", position.X, position.Y)
}

type Size struct {
	Width  int
	Height int
}

func (size Size) String() string {
	return fmt.Sprintf("%dx%d", size.Width, size.Height)
}
