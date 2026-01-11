package tui

import (
	"github.com/gdamore/tcell/v3"
)

func NewEmptyGrid(size Size) Grid {
	style := tcell.StyleDefault

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
