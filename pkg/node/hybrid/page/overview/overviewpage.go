package overview

import (
	"pkb-agent/pkg/node"
	"pkb-agent/pkg/node/hybrid/page"
	"pkb-agent/tui"
)

type Page struct {
	pages []page.Page
}

func New(pages []page.Page) *Page {
	return &Page{
		pages: pages,
	}
}

func (page *Page) CreateViewer(messageQueue tui.MessageQueue) tui.Component {
	return NewPageComponent(messageQueue, page)
}

func (page *Page) GetCaption() string {
	return "Overview"
}

func (page *Page) GetActions() []node.Action {
	return nil
}
