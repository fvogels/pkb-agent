package atom

import (
	"encoding/binary"
	"io"
	"pkb-agent/tui"
	"pkb-agent/tui/component/label"
	"pkb-agent/tui/data"
	"pkb-agent/util"
	"strings"
)

const TypeID uint32 = 1

type RawNode struct {
	name  string
	links []string
}

func (node *RawNode) GetName() string {
	return node.name
}

func (node *RawNode) GetSearchStrings() []string {
	return util.Words(strings.ToLower(util.RemoveAccents(node.name)))
}

func (node *RawNode) GetLinks() []string {
	return node.links
}

func (node *RawNode) GetViewer(messageQueue tui.MessageQueue) tui.Component {
	return label.New(messageQueue, "atomviewer", data.NewConstant("atom"))
}

func (node *RawNode) Serialize(writer io.Writer) error {
	bufferSize := 0
	bufferSize += 4              // type id
	bufferSize += 4              // len(name)
	bufferSize += len(node.name) // name
	buffer := make([]byte, 0, bufferSize)

	// Type ID
	binary.LittleEndian.AppendUint32(buffer, TypeID)

	// Name
	binary.LittleEndian.AppendUint32(buffer, uint32(len(node.name)))
	buffer = append(buffer, ([]byte)(node.name)...)

	if _, err := writer.Write(buffer); err != nil {
		return err
	}

	return nil
}
