package hybrid

import (
	"io"
	markdownpage "pkb-agent/graph/node/hybrid/pages/markdown"
	"pkb-agent/util"
	"pkb-agent/util/multifile"
	"pkb-agent/util/pathlib"
	"strings"
	"weak"

	tea "github.com/charmbracelet/bubbletea"
)

const TypeID uint32 = 4

type RawNode struct {
	name  string
	links []string
	path  pathlib.Path
	data  weak.Pointer[nodeData]
}

type nodeData struct {
	pages []Page
}

type Page interface {
	GetCaption() string
	CreateViewer() tea.Model
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

func (node *RawNode) getData() (*nodeData, error) {
	data := node.data.Value()

	if data == nil {
		file, err := multifile.Load(node.path)
		if err != nil {
			return nil, err
		}

		data = &nodeData{
			pages: node.loadPages(file),
		}

		node.data = weak.Make(data)
	}

	return data, nil
}

func (node *RawNode) GetViewer() tea.Model {
	data, err := node.getData()
	if err != nil {
		panic("error loading data")
	}

	return NewViewer(node, data)
}

func (node *RawNode) Serialize(writer io.Writer) error {
	// bufferSize := 0
	// bufferSize += 4              // type id
	// bufferSize += 4              // len(name)
	// bufferSize += len(node.name) // name
	// bufferSize += 4              // len(url)
	// bufferSize += len(node.url)  // url
	// buffer := make([]byte, 0, bufferSize)

	// // Type ID
	// binary.LittleEndian.AppendUint32(buffer, TypeID)

	// // Name
	// binary.LittleEndian.AppendUint32(buffer, uint32(len(node.name)))
	// buffer = append(buffer, ([]byte)(node.name)...)

	// // URL
	// binary.LittleEndian.AppendUint32(buffer, uint32(len(node.url)))
	// buffer = append(buffer, ([]byte)(node.url)...)

	// if _, err := writer.Write(buffer); err != nil {
	// 	return err
	// }

	// TODO

	return nil
}

func (node *RawNode) loadPages(file *multifile.MultiFile) []Page {
	pages := []Page{}

	for _, segment := range file.Segments {
		switch segment.Type {
		case "markdown":
			caption, foundCaption := segment.Attributes["caption"]
			if !foundCaption {
				caption = "untitled"
			}
			source := strings.Join(segment.Contents, "\n")
			page := markdownpage.New(caption, source)

			pages = append(pages, page)
		}
	}

	return pages
}
