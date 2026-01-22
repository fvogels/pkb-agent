package keyview

import (
	"log/slog"
	"pkb-agent/persistent/list"
	"pkb-agent/tui"
	"pkb-agent/tui/data"
	"pkb-agent/tui/debug"
	"pkb-agent/util/uid"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Component struct {
	tui.ComponentBase
	keyBindings      data.View[list.List[tui.KeyBinding]]
	cachedGrid       tui.Grid
	keyStyle         *tui.Style
	descriptionStyle *tui.Style
	emptyStyle       *tui.Style
}

func New(messageQueue tui.MessageQueue, name string, keyBindings data.Value[list.List[tui.KeyBinding]]) *Component {
	keyStyle := tcell.StyleDefault.Background(color.NewHexColor(0xAAAAFF))
	descriptionStyle := tcell.StyleDefault.Background(color.NewHexColor(0x8888FF))
	emptyStyle := tcell.StyleDefault.Background(color.Reset)

	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         name,
			MessageQueue: messageQueue,
		},
		keyBindings:      data.NewView(keyBindings),
		cachedGrid:       nil,
		keyStyle:         &keyStyle,
		descriptionStyle: &descriptionStyle,
		emptyStyle:       &emptyStyle,
	}

	return &component
}

func (component *Component) Handle(message tui.Message) {
	debug.LogMessage(message)

	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)
	}
}

func (component *Component) Render() tui.Grid {
	slog.Debug("Rendering keyview")

	if component.keyBindings.Updated() || component.cachedGrid == nil {
		component.cachedGrid = component.renderKeyBindings()
	}

	return component.cachedGrid
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size
	component.cachedGrid = nil
}

func (component *Component) renderKeyBindings() tui.Grid {
	cell := tui.Cell{
		Contents: ' ',
		Style:    component.emptyStyle,
	}
	result := tui.NewMaterializedGrid(component.Size, func(tui.Position) tui.Cell { return cell })

	x := 0
	write := func(contents rune, style *tui.Style) {
		if x < result.Size().Width {
			result.Set(tui.Position{X: x, Y: 0}, tui.Cell{Contents: contents, Style: style})
			x++
		}
	}
	list.ForEach(component.keyBindings.Get(), func(index int, keyBinding tui.KeyBinding) {
		write(' ', component.keyStyle)
		for _, r := range []rune(keyBinding.Key) {
			write(r, component.keyStyle)
		}
		write(' ', component.keyStyle)

		write(' ', component.descriptionStyle)
		for _, r := range []rune(keyBinding.Description) {
			write(r, component.descriptionStyle)
		}
		write(' ', component.descriptionStyle)
		write(' ', component.emptyStyle)
	})

	return result
}
