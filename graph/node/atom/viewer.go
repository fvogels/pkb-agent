package atom

import "pkb-agent/tui"

type Component struct {
	size tui.Size
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)
	}
}

func (component *Component) Render() tui.Grid {
	return nil
}

func (component *Component) onResize(message tui.MsgResize) {
	component.size = message.Size
}

// type Model struct {
// 	size util.Size
// }

// func NewViewer() Model {
// 	return Model{}
// }

// func (model Model) Init() tea.Cmd {
// 	return nil
// }

// func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
// 	return model.TypedUpdate(message)
// }

// func (model Model) TypedUpdate(message tea.Msg) (Model, tea.Cmd) {
// 	switch message := message.(type) {
// 	case tea.WindowSizeMsg:
// 		return model.onResize(message)

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

// func (model Model) View() string {
// 	return "atom"
// }
