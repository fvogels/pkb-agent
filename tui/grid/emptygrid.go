package grid

import (
	"pkb-agent/tui/position"
	"pkb-agent/tui/size"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func NewEmptyGrid(size size.Size) FiniteGrid {
	style := tcell.StyleDefault.Foreground(color.Reset).Background(color.Reset)

	result := emptyGrid{
		size:  size,
		style: &style,
	}

	return &result
}

type emptyGrid struct {
	size  size.Size
	style *tcell.Style
}

func (g *emptyGrid) Size() size.Size {
	return g.size
}

func (g *emptyGrid) At(position.Position) Cell {
	cell := Cell{
		Contents: ' ',
		Style:    g.style,
	}

	return cell
}
