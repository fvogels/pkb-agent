package grid

import "pkb-agent/tui/position"

type positionedGrid struct {
	grid   FiniteGrid
	offset position.Position
}
