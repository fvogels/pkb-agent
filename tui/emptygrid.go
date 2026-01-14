package tui

import (
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func NewEmptyGrid(size Size) Grid {
	style := tcell.StyleDefault.Foreground(color.Reset).Background(color.Reset)

	result := emptyGrid{
		size:  size,
		style: &style,
	}

	return &result
}

type emptyGrid struct {
	size  Size
	style *Style
}

func (grid *emptyGrid) GetSize() Size {
	return grid.size
}

func (grid *emptyGrid) Get(Position) Cell {
	cell := Cell{
		Contents: ' ',
		Style:    grid.style,
	}

	return cell
}
