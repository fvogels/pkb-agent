package markdownpage

import (
	"pkb-agent/pkg/node"
	"pkb-agent/tui"
	"pkb-agent/tui/component/markdownview"
	"pkb-agent/tui/data"
)

type Page struct {
	caption string
	source  string
}

func New(caption string, source string) *Page {
	return &Page{
		caption: caption,
		source:  source,
	}
}

func (page *Page) CreateViewer(messageQueue tui.MessageQueue) tui.Component {
	return markdownview.New(messageQueue, data.NewConstant(page.source))
}

func (page *Page) GetCaption() string {
	return page.caption
}

func (page *Page) GetActions() []node.Action {
	return nil
}
