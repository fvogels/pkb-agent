package ansigrid

import (
	"pkb-agent/tui"
	"pkb-agent/tui/position"
	"strings"

	"github.com/gdamore/tcell/v3"
	ansi "github.com/ktr0731/go-ansisgr"
)

type ansiGrid struct {
	size      tui.Size
	cells     [][]tui.Cell
	emptyCell tui.Cell
}

func Parse(str string, emptyStyle *tui.Style) tui.Grid {
	cells := [][]tui.Cell{}
	strings.Lines(str)(func(line string) bool {
		cells = append(cells, parseLine(line))
		return true
	})

	width := 0
	for _, row := range cells {
		if len(row) > width {
			width = len(row)
		}
	}

	height := len(cells)

	return &ansiGrid{
		size: tui.Size{
			Width:  width,
			Height: height,
		},
		cells: cells,
		emptyCell: tui.Cell{
			Contents: ' ',
			Style:    emptyStyle,
			OnClick:  nil,
		},
	}
}

func parseLine(line string) []tui.Cell {
	iterator := ansi.NewIterator(line)
	cells := []tui.Cell{}

	for {
		char, style, ok := iterator.Next()
		if !ok {
			return cells
		}

		cell := tui.Cell{
			Contents: char,
			Style:    translateStyle(style),
			OnClick:  nil,
		}

		cells = append(cells, cell)
	}
}

func translateStyle(ansiStyle ansi.Style) *tui.Style {
	tuiStyle := tcell.StyleDefault

	if color, ok := ansiStyle.Foreground(); ok {
		switch color.Mode() {
		case ansi.Mode16:
			tuiStyle = tuiStyle.Foreground(tcell.PaletteColor(color.Value() - 30))
		case ansi.Mode256:
			tuiStyle = tuiStyle.Foreground(tcell.PaletteColor(color.Value()))
		case ansi.ModeRGB:
			r, g, b := color.RGB()
			tuiStyle = tuiStyle.Foreground(tcell.NewRGBColor(int32(r), int32(g), int32(b)))
		}
	}

	if color, valid := ansiStyle.Background(); valid {
		switch color.Mode() {
		case ansi.Mode16:
			tuiStyle = tuiStyle.Background(tcell.PaletteColor(color.Value() - 40))
		case ansi.Mode256:
			tuiStyle = tuiStyle.Background(tcell.PaletteColor(color.Value()))
		case ansi.ModeRGB:
			r, g, b := color.RGB()
			tuiStyle = tuiStyle.Background(tcell.NewRGBColor(int32(r), int32(g), int32(b)))
		}
	}

	return &tuiStyle
}

func (grid *ansiGrid) Size() tui.Size {
	return grid.size
}

func (grid *ansiGrid) At(position position.Position) tui.Cell {
	x := position.X
	y := position.Y

	row := grid.cells[y]
	if x < len(row) {
		return row[x]
	} else {
		return grid.emptyCell
	}
}
