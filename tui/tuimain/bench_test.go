//go:build test

package tuimain

import (
	"pkb-agent/tui"
	"pkb-agent/tui/component/border"
	"pkb-agent/tui/component/stringlist"
	"pkb-agent/tui/component/strlist"
	"pkb-agent/tui/data"
	"testing"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func BenchmarkStringList(b *testing.B) {
	style := tcell.StyleDefault.Background(color.Green).Foreground(color.Reset)
	selectedStyle := tcell.StyleDefault.Background(color.Gray).Foreground(color.Reset)
	items := data.NewSliceList([]string{"a", "bb", "ccc", "dddd", "eeeee", "ffffff", "ggggggg", "hhhhhhhh"})
	selectedItem := data.NewVariable(0)

	list := border.New(stringlist.New(items, selectedItem, &style, &selectedStyle), style)
	list.Handle(tui.MsgResize{
		Size: tui.Size{Width: 500, Height: 100},
	})

	for b.Loop() {
		grid := list.Render()
		gridSize := grid.GetSize()

		for y := range gridSize.Height {
			for x := range gridSize.Width {
				position := tui.Position{X: x, Y: y}
				grid.Get(position)
			}
		}
	}
}

func BenchmarkStrList(b *testing.B) {
	style := tcell.StyleDefault.Background(color.Green).Foreground(color.Reset)
	selectedStyle := tcell.StyleDefault.Background(color.Gray).Foreground(color.Reset)
	items := data.NewSliceList([]string{"a", "bb", "ccc", "dddd", "eeeee", "ffffff", "ggggggg", "hhhhhhhh"})
	selectedItem := data.NewVariable(0)

	list := border.New(strlist.New(items, selectedItem, &style, &selectedStyle), style)
	list.Handle(tui.MsgResize{
		Size: tui.Size{Width: 500, Height: 100},
	})

	for b.Loop() {
		grid := list.Render()
		gridSize := grid.GetSize()

		for y := range gridSize.Height {
			for x := range gridSize.Width {
				position := tui.Position{X: x, Y: y}
				grid.Get(position)
			}
		}
	}
}
