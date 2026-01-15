package bookmark

import (
	"pkb-agent/persistent/list"
	"pkb-agent/tui"
	"pkb-agent/tui/application/messages"
	"pkb-agent/tui/component/label"
	"pkb-agent/tui/data"
	"pkb-agent/ui/uid"
)

type Component struct {
	tui.ComponentBase
	child tui.Component
}

func NewViewer(messageQueue tui.MessageQueue) *Component {
	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "unnamed bookmark viewer",
			MessageQueue: messageQueue,
		},
		child: label.New(messageQueue, "bookmark label", data.NewConstant("bookmark!")),
	}

	return &component
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgActivate:
		component.onActivate()

	case tui.MsgResize:
		component.onResize(message)
	}
}

func (component *Component) onActivate() {
	component.MessageQueue.Enqueue(messages.MsgSetNodeKeyBindings{
		Bindings: list.New[tui.KeyBinding](),
	})
}

func (component *Component) Render() tui.Grid {
	return component.child.Render()
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size
	component.child.Handle(message)
}

// var keyMap = struct {
// 	OpenInBrowser key.Binding
// }{
// 	OpenInBrowser: key.NewBinding(
// 		key.WithKeys("o"),
// 		key.WithHelp("o", "open"),
// 	),
// }

// type Model struct {
// 	size util.Size
// 	url  string
// }

// func NewViewer(url string) Model {
// 	return Model{
// 		url: url,
// 	}
// }

// func (model Model) Init() tea.Cmd {
// 	return model.signalKeybindingsUpdate()
// }

// func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
// 	return model.TypedUpdate(message)
// }

// func (model Model) TypedUpdate(message tea.Msg) (Model, tea.Cmd) {
// 	switch message := message.(type) {
// 	case tea.WindowSizeMsg:
// 		return model.onResize(message)

// 	case tea.KeyMsg:
// 		return model.onKeyPressed(message)

// 	default:
// 		return model, nil
// 	}
// }

// func (model Model) onResize(message tea.WindowSizeMsg) (Model, tea.Cmd) {
// 	model.size = util.Size{
// 		Width:  message.Width,
// 		Height: message.Height,
// 	}

// 	return model, nil
// }

// func (model Model) onKeyPressed(message tea.KeyMsg) (Model, tea.Cmd) {
// 	switch {
// 	case key.Matches(message, keyMap.OpenInBrowser):
// 		return model.onOpenInBrowser()

// 	default:
// 		return model, nil
// 	}
// }

// func (model Model) onOpenInBrowser() (Model, tea.Cmd) {
// 	if err := extern.OpenURLInBrowser(model.url); err != nil {
// 		panic("failed to open browser")
// 	}

// 	return model, nil
// }

// func (model Model) View() string {
// 	return model.url
// }

// func (model Model) signalKeybindingsUpdate() tea.Cmd {
// 	return func() tea.Msg {
// 		return node.MsgUpdateNodeViewerBindings{
// 			KeyBindings: []key.Binding{
// 				keyMap.OpenInBrowser,
// 			},
// 		}
// 	}
// }
