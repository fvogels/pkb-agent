package overview

import (
	"pkb-agent/persistent/list"
	"pkb-agent/pkg/node/hybrid/page"
	"pkb-agent/tui"
	"pkb-agent/tui/component/stringsview"
	"pkb-agent/tui/data"
	"pkb-agent/tui/grid"
	"pkb-agent/util"
	"pkb-agent/util/uid"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type pageComponent struct {
	tui.ComponentBase
	parent *Page
	root   tui.Component
}

func NewPageComponent(messageQueue tui.MessageQueue, parent *Page) *pageComponent {
	component := pageComponent{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "nameless overview page",
			MessageQueue: messageQueue,
		},
		root: createRoot(messageQueue, parent.pages),
	}

	return &component
}

func createRoot(messageQueue tui.MessageQueue, pages []page.Page) tui.Component {
	// return label.New(messageQueue, "label", data.NewConstant(fmt.Sprintf("%d pages", len(pages))))

	captionStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	captionsAsSlice := util.Map(pages, func(page page.Page) stringsview.Item {
		caption := page.GetCaption()
		return stringsview.Item{
			Runes: []rune(caption),
			Style: &captionStyle,
		}
	})
	captionsAsList := list.FromSlice(captionsAsSlice)
	return stringsview.New(messageQueue, data.NewConstant(captionsAsList))
}

func (component *pageComponent) Render() grid.FiniteGrid {
	return component.root.Render()
}

func (component *pageComponent) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgKey:
		component.onKey(message)

	case tui.MsgStateUpdated:
		component.root.Handle(message)
		component.onStateUpdated()

	default:
		component.root.Handle(message)
	}
}

func (component *pageComponent) onKey(message tui.MsgKey) {
	// No key bindings for this page
}

func (component *pageComponent) onStateUpdated() {
	component.MessageQueue.Enqueue(page.MsgSetPageKeyBindings{
		Bindings: list.FromItems[tui.KeyBinding](),
	})
}
