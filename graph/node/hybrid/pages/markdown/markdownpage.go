package markdownpage

import (
	"pkb-agent/graph/node"
	"pkb-agent/tui"
	"pkb-agent/tui/component/label"
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

func (page *Page) CreateViewer() tui.Component {
	return label.New("markdownviewer", data.NewConstant("markdown page"))
}

func (page *Page) GetCaption() string {
	return page.caption
}

func (page *Page) GetActions() []node.Action {
	return nil
}
