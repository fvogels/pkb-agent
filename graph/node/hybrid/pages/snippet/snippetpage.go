package snippetpage

import (
	"pkb-agent/graph/node"
	"pkb-agent/graph/node/hybrid/actions/clipboard"
	"pkb-agent/tui"
	"pkb-agent/tui/component/label"
	"pkb-agent/tui/data"
)

type Page struct {
	caption  string
	source   string
	language string
	actions  []node.Action
}

func New(caption string, source string, language string) *Page {
	return &Page{
		caption:  caption,
		source:   source,
		language: language,
		actions: []node.Action{
			clipboard.New(source),
		},
	}
}

func (page *Page) CreateViewer() tui.Component {
	return label.New("snippetviewer", data.NewConstant("snippet page"))
}

func (page *Page) GetCaption() string {
	return page.caption
}

func (page *Page) GetActions() []node.Action {
	return page.actions
}
