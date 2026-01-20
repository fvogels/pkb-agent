package linksview

import (
	"pkb-agent/persistent/list"
	"pkb-agent/pkg"
	"pkb-agent/tui"
	"pkb-agent/tui/application/messages"
	"pkb-agent/tui/component/captioned"
	"pkb-agent/tui/component/stringsview"
	"pkb-agent/tui/component/vstack"
	"pkb-agent/tui/data"
	"pkb-agent/util/uid"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Component struct {
	tui.ComponentBase
	root tui.Component
	node *pkg.Node
}

func New(messageQueue tui.MessageQueue, node *pkg.Node) *Component {
	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "unnamed links view",
			MessageQueue: messageQueue,
		},
		root: createRoot(messageQueue, node),
		node: node,
	}

	return &component
}

func createRoot(messageQueue tui.MessageQueue, node *pkg.Node) tui.Component {
	style := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)

	links := data.NewConstant(getLinkItems(node, &style))
	linksCaption := data.NewConstant([]rune("Links"))

	backLinks := data.NewConstant(getBacklinkItems(node, &style))
	backlinksCaption := data.NewConstant([]rune("Backlinks"))

	linksRoot := captioned.NewMeasurable(
		messageQueue,
		linksCaption,
		stringsview.New(messageQueue, links),
	)

	backlinksRoot := captioned.NewMeasurable(
		messageQueue,
		backlinksCaption,
		stringsview.New(messageQueue, backLinks),
	)

	root := vstack.New(
		messageQueue,
		list.FromItems[tui.MeasurableComponent](linksRoot, backlinksRoot),
	)

	return root
}

func (component *Component) Render() tui.Grid {
	return component.root.Render()
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgStateUpdated:
		component.root.Handle(message)
		component.onStateUpdated()

	default:
		component.root.Handle(message)
	}
}

func (component *Component) onStateUpdated() {
	component.MessageQueue.Enqueue(messages.MsgSetNodeKeyBindings{
		Bindings: list.FromItems[tui.KeyBinding](),
	})
}

func getLinkItems(node *pkg.Node, style *tui.Style) list.List[stringsview.Item] {
	linkedNodesList := list.FromSlice(node.GetLinks())

	return list.MapList(
		linkedNodesList,
		func(linkedNode *pkg.Node) stringsview.Item {
			name := linkedNode.GetName()

			return stringsview.Item{
				Runes: []rune(name),
				Style: style,
			}
		},
	)
}

func getBacklinkItems(node *pkg.Node, style *tui.Style) list.List[stringsview.Item] {
	linkedNodesList := list.FromSlice(node.GetBacklinks())

	return list.MapList(
		linkedNodesList,
		func(linkedNode *pkg.Node) stringsview.Item {
			name := linkedNode.GetName()

			return stringsview.Item{
				Runes: []rune(name),
				Style: style,
			}
		},
	)
}
