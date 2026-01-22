package tui

import (
	"pkb-agent/tui/position"
	"pkb-agent/tui/size"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func NewEmptyGrid(size size.Size) Grid {
	style := tcell.StyleDefault.Foreground(color.Reset).Background(color.Reset)

	result := emptyGrid{
		size:  size,
		style: &style,
	}

	return &result
}

type emptyGrid struct {
	size  size.Size
	style *Style
}

func (grid *emptyGrid) Size() size.Size {
	return grid.size
}

func (grid *emptyGrid) At(position.Position) Cell {
	cell := Cell{
		Contents: ' ',
		Style:    grid.style,
	}

	return cell
}
