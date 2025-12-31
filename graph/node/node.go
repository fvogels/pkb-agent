package node

import (
	"io"
	"pkb-agent/util/pathlib"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Deserializer interface {
	Deserialize(io.Reader) (RawNode, error)
}

type RawNode interface {
	GetName() string
	GetSearchStrings() []string
	GetLinks() []string
	GetViewer() tea.Model
	Serialize(io.Writer) error
}

type Action interface {
	GetDescription() string
	Perform()
}

type Loader interface {
	Load(path pathlib.Path, callback func(node RawNode) error) error
}

type MsgUpdateNodeViewerBindings struct {
	KeyBindings []key.Binding
}
