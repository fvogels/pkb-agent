package tui

import (
	"pkb-agent/tui/grid"
	"pkb-agent/tui/position"
	"pkb-agent/tui/size"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func NewEmptyGrid(size size.Size) grid.Grid {
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

func (g *emptyGrid) Size() size.Size {
	return g.size
}

func (g *emptyGrid) At(position.Position) grid.Cell {
	cell := grid.Cell{
		Contents: ' ',
		Style:    g.style,
	}

	return cell
}
