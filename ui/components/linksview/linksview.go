package linksview

import (
	"log/slog"
	"pkb-agent/ui/components/listview"
	"pkb-agent/ui/debug"
	"pkb-agent/ui/layout"
	"pkb-agent/ui/layout/horizontal"
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	size util.Size

	links     List
	backlinks List

	linksView     listview.Model[string]
	backlinksView listview.Model[string]
	layout        layout.Layout[Model]
}

type List interface {
	At(index int) string
	Length() int
}

type emptyList struct{}

func (emptyList *emptyList) At(index int) string {
	panic("invalid operation")
}

func (emptyList *emptyList) Length() int {
	return 0
}

func New() Model {
	renderer := func(link string) string {
		return link
	}

	layoutRoot := horizontal.New[Model]()
	layoutRoot.Add(
		func(size util.Size) int { return size.Width - size.Width/2 },
		layout.Wrap(func(m *Model) *listview.Model[string] { return &m.linksView }),
	)
	layoutRoot.Add(
		func(size util.Size) int { return size.Width / 2 },
		layout.Wrap(func(m *Model) *listview.Model[string] { return &m.backlinksView }),
	)

	linksView := listview.New(renderer, false, wrapLinksViewMessage)
	linksView.SetEmptyListMessage("no links")
	backlinksView := listview.New(renderer, false, wrapBacklinksViewMessage)
	backlinksView.SetEmptyListMessage("no backlinks")

	return Model{
		links:         &emptyList{},
		backlinks:     &emptyList{},
		linksView:     linksView,
		backlinksView: backlinksView,
		layout:        &layoutRoot,
	}
}

func (model Model) Init() tea.Cmd {
	return tea.Batch(
		model.linksView.Init(),
		model.backlinksView.Init(),
	)
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return model.TypedUpdate(message)
}

func (model Model) TypedUpdate(message tea.Msg) (Model, tea.Cmd) {
	debug.ShowBubbleTeaMessage(message)

	switch message := message.(type) {
	case tea.WindowSizeMsg:
		return model.onResize(message)

	case MsgSetLinks:
		return model.onSetLinks(message)

	case msgLinksListWrapper:
		slog.Warn("swallowed message from links list")
		return model, nil

	case msgBacklinksListWrapper:
		slog.Warn("swallowed message from backlinks list")
		return model, nil

	default:
		commands := []tea.Cmd{}

		util.UpdateChild(&model.linksView, message, &commands)
		util.UpdateChild(&model.backlinksView, message, &commands)

		return model, tea.Batch(commands...)
	}
}

func (model Model) View() string {
	return model.layout.LayoutView(&model)
}

func (model Model) onResize(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	command := model.layout.LayoutResize(&model, model.size)

	return model, command
}

func wrapLinksViewMessage(message tea.Msg) tea.Msg {
	return msgLinksListWrapper{
		wrapped: message,
	}
}

func wrapBacklinksViewMessage(message tea.Msg) tea.Msg {
	return msgBacklinksListWrapper{
		wrapped: message,
	}
}

func (model Model) onSetLinks(message MsgSetLinks) (Model, tea.Cmd) {
	model.links = message.Links
	model.backlinks = message.Backlinks

	commands := []tea.Cmd{}
	util.UpdateChild(&model.linksView, listview.MsgSetItems[string]{
		Items: model.links,
	}, &commands)
	util.UpdateChild(&model.backlinksView, listview.MsgSetItems[string]{
		Items: model.backlinks,
	}, &commands)

	return model, tea.Batch(commands...)
}
