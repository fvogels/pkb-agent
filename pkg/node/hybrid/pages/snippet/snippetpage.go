package snippetpage

import (
	"pkb-agent/pkg/node"
	"pkb-agent/pkg/node/hybrid/actions/clipboard"
	"pkb-agent/tui"
	"pkb-agent/tui/component/snippetview"
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
	source := snippetview.Source{
		Contents: page.source,
		Language: page.language,
	}
	return snippetview.New(data.NewConstant(source))
}

func (page *Page) GetCaption() string {
	return page.caption
}

func (page *Page) GetActions() []node.Action {
	return page.actions
}
