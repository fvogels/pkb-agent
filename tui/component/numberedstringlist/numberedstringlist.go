package numberedstringlist

import (
	"fmt"
	"pkb-agent/persistent/list"
	"pkb-agent/tui"
	"pkb-agent/tui/data"
	"pkb-agent/tui/debug"
	"pkb-agent/tui/grid"
	"pkb-agent/tui/position"
	"pkb-agent/util"
	"pkb-agent/util/uid"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Component struct {
	tui.ComponentBase
	items       data.Value[list.List[string]]
	emptyStyle  *tui.Style
	itemStyle   *tui.Style
	numberStyle *tui.Style
}

func New(messageQueue tui.MessageQueue, items data.Value[list.List[string]]) *Component {
	defaultEmptyStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	defaultItemStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	defaultNumberStyle := tcell.StyleDefault.Background(color.Blue).Foreground(color.Reset)

	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "unnamed numberedstringlist",
			MessageQueue: messageQueue,
		},
		items:       items,
		itemStyle:   &defaultItemStyle,
		emptyStyle:  &defaultEmptyStyle,
		numberStyle: &defaultNumberStyle,
	}

	return &component
}

func (component *Component) SetEmptyStyle(emptyStyle *tui.Style) {
	component.emptyStyle = emptyStyle
}

func (component *Component) SetItemStyle(itemStyle *tui.Style) {
	component.itemStyle = itemStyle
}

func (component *Component) SetNumberStyle(numberStyle *tui.Style) {
	component.numberStyle = numberStyle
}

func (component *Component) Handle(message tui.Message) {
	debug.LogMessage(message)

	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)

	case tui.MsgStateUpdated:
		component.onStateUpdated()
	}
}

func (component *Component) Render() grid.FiniteGrid {
	width := component.Size.Width
	height := component.Size.Height

	result := grid.NewMaterializedGrid(component.Size, func(pos position.Position) grid.Cell {
		return grid.Cell{
			Contents: ' ',
			Style:    component.emptyStyle,
		}
	})

	items := component.items.Get()

	y := 0
	yMax := util.MinInt(height, items.Size())
	for y < yMax {
		runes := []rune(fmt.Sprintf(" %s", items.At(y)))
		x := 0

		// Number
		for _, char := range fmt.Sprintf(" %d ", y+1) {
			result.Set(
				position.Position{X: x, Y: y},
				grid.Cell{Contents: char, Style: component.numberStyle},
			)

			x += 1
		}

		// Item label
		i := 0
		for i < len(runes) && x < width {
			result.Set(
				position.Position{X: x, Y: y},
				grid.Cell{Contents: runes[i], Style: component.itemStyle},
			)
			i += 1
			x += 1
		}

		// Remainder of the row
		for x < width {
			result.Set(
				position.Position{X: x, Y: y},
				grid.Cell{Contents: ' ', Style: component.itemStyle},
			)

			x += 1
		}

		y += 1
	}

	return result
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size
}

func (component *Component) onStateUpdated() {
	// NOP
}
