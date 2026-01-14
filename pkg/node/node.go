package node

import (
	"io"
	"pkb-agent/tui"
	"pkb-agent/util/pathlib"

	"github.com/charmbracelet/bubbles/key"
)

type Deserializer interface {
	Deserialize(io.Reader) (RawNode, error)
}

type RawNode interface {
	GetName() string
	GetSearchStrings() []string
	GetLinks() []string
	GetViewer() tui.Component
	Serialize(io.Writer) error
}

type Action interface {
	GetDescription() string
	Perform()
}

type Loader interface {
	Load(pathlib.Path, any, func(node RawNode) error) error
}

type MsgUpdateNodeViewerBindings struct {
	KeyBindings []key.Binding
}
