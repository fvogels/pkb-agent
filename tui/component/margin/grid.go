package margin

import (
	"fmt"
	tuigrid "pkb-agent/tui/grid"
	"pkb-agent/tui/position"
	"pkb-agent/tui/size"
)

type grid struct {
	parent    *Component
	childGrid tuigrid.FiniteGrid
}

func newGrid(parent *Component) tuigrid.FiniteGrid {
	grid := grid{
		parent:    parent,
		childGrid: parent.child.Render(),
	}

	return &grid
}

func (grid *grid) Size() size.Size {
	childSize := grid.childGrid.Size()

	return size.Size{
		Width:  childSize.Width + 2,
		Height: childSize.Height + 2,
	}
}

func (g *grid) At(pos position.Position) tuigrid.Cell {
	if !g.isValidPosition(pos) {
		size := g.Size()
		panic(fmt.Sprintf("invalid position (%d, %d), size %dx%d in component %s", pos.X, pos.Y, size.Width, size.Height, g.parent.Name))
	}

	x := pos.X
	y := pos.Y
	marginSize := g.parent.marginSize

	if g.inMargin(pos) {
		return tuigrid.Cell{
			Contents: ' ',
			Style:    g.parent.style,
			OnClick:  nil,
		}
	} else {
		return g.childGrid.At(
			position.Position{
				X: x - marginSize,
				Y: y - marginSize,
			},
		)
	}
}

func (g *grid) isValidPosition(position position.Position) bool {
	return g.Size().IsValidPosition(position)
}

func (g *grid) inMargin(pos position.Position) bool {
	x := pos.X
	y := pos.Y
	width := g.Size().Width
	height := g.Size().Height
	size := g.parent.marginSize

	return x < size || y < size || x >= width-size || y >= height-size
}
