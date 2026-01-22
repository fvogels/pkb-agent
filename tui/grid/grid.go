package grid

import (
	"pkb-agent/tui/position"
	"pkb-agent/tui/size"

	"github.com/gdamore/tcell/v3"
)

type Grid interface {
	Size() size.Size
	At(position position.Position) Cell
}

type Cell struct {
	Contents rune
	Style    *tcell.Style
	OnClick  func()
}
