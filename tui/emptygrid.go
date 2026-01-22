package tui

import (
	"pkb-agent/tui/position"

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

func (grid *emptyGrid) Size() Size {
	return grid.size
}

func (grid *emptyGrid) At(position.Position) Cell {
	cell := Cell{
		Contents: ' ',
		Style:    grid.style,
	}

	return cell
}
