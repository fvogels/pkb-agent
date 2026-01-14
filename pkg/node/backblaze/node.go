package backblaze

import (
	"encoding/binary"
	"io"
	"pkb-agent/tui"
	"pkb-agent/tui/component/label"
	"pkb-agent/tui/data"
	"pkb-agent/util"
	"strings"
)

const TypeID uint32 = 3

type RawNode struct {
	name     string
	links    []string
	bucket   string
	filename string
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

func (node *RawNode) GetViewer(tui.MessageQueue) tui.Component {
	return label.New("backblazeviewer", data.NewConstant("backblaze"))
}

func (node *RawNode) Serialize(writer io.Writer) error {
	bufferSize := 0
	bufferSize += 4              // type id
	bufferSize += 4              // len(name)
	bufferSize += len(node.name) // name
	bufferSize += 4              // len(bucket)
	bufferSize += len(node.name) // bucket
	bufferSize += 4              // len(bucket)
	bufferSize += len(node.name) // bucket
	buffer := make([]byte, 0, bufferSize)

	// Type ID
	binary.LittleEndian.AppendUint32(buffer, TypeID)

	// Name
	binary.LittleEndian.AppendUint32(buffer, uint32(len(node.name)))
	buffer = append(buffer, ([]byte)(node.name)...)

	// Bucket
	binary.LittleEndian.AppendUint32(buffer, uint32(len(node.bucket)))
	buffer = append(buffer, ([]byte)(node.bucket)...)

	// Filename
	binary.LittleEndian.AppendUint32(buffer, uint32(len(node.filename)))
	buffer = append(buffer, ([]byte)(node.filename)...)

	if _, err := writer.Write(buffer); err != nil {
		return err
	}

	return nil
}
