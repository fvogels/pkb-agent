package linksview

import (
	"log/slog"
	"pkb-agent/ui/components/listview"
	"pkb-agent/ui/debug"
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	size util.Size

	links     List
	backlinks List

	linksView     listview.Model[string]
	backlinksView listview.Model[string]
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

	linksView := listview.New(renderer, false, wrapLinksViewMessage)
	linksView.SetEmptyListMessage("no links")
	backlinksView := listview.New(renderer, false, wrapBacklinksViewMessage)
	backlinksView.SetEmptyListMessage("no backlinks")

	return Model{
		links:         &emptyList{},
		backlinks:     &emptyList{},
		linksView:     linksView,
		backlinksView: backlinksView,
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
	return lipgloss.JoinVertical(
		0,
		model.linksView.View(),
		model.backlinksView.View(),
	)
}

func (model Model) onResize(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	commands := []tea.Cmd{}

	util.UpdateChild(&model.linksView, tea.WindowSizeMsg{
		Width:  model.size.Width,
		Height: message.Height / 2,
	}, &commands)

	util.UpdateChild(&model.backlinksView, tea.WindowSizeMsg{
		Width:  model.size.Width,
		Height: message.Height - message.Height/2,
	}, &commands)

	return model, tea.Batch(commands...)
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
