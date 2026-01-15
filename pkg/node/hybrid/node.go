package hybrid

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"pkb-agent/pkg/node"
	"pkb-agent/pkg/node/hybrid/actions/www"
	"pkb-agent/pkg/node/hybrid/page"
	markdownpage "pkb-agent/pkg/node/hybrid/page/markdown"
	snippetpage "pkb-agent/pkg/node/hybrid/page/snippet"
	"pkb-agent/tui"
	"pkb-agent/util"
	"pkb-agent/util/multifile"
	"pkb-agent/util/pathlib"
	"strings"
	"weak"
)

const TypeID uint32 = 4

type RawNode struct {
	name  string
	links []string
	path  pathlib.Path
	data  weak.Pointer[nodeData]
}

type nodeData struct {
	pages   []page.Page
	actions []node.Action
}

func (rawNode *RawNode) GetName() string {
	return rawNode.name
}

func (rawNode *RawNode) GetSearchStrings() []string {
	return util.Words(strings.ToLower(util.RemoveAccents(rawNode.name)))
}

func (rawNode *RawNode) GetLinks() []string {
	return rawNode.links
}

func (rawNode *RawNode) getData() (*nodeData, error) {
	data := rawNode.data.Value()

	if data == nil {
		file, err := multifile.Load(rawNode.path)
		if err != nil {
			return nil, err
		}

		metadataSegment := file.FindSegmentOfType("metadata")
		metadata, err := parseMetadata(metadataSegment.Contents)
		if err != nil {
			slog.Debug("Error parsing hybrid node's metadata; should have been caught earlier")
			return nil, err
		}

		actions, err := rawNode.parseActions(metadata.Actions)
		if err != nil {
			slog.Debug("Error parsing hybrid node's actions")
			return nil, fmt.Errorf("failed to parse actions: %w", err)
		}

		data = &nodeData{
			pages:   rawNode.loadPages(file),
			actions: actions,
		}

		rawNode.data = weak.Make(data)
	}

	return data, nil
}

func (rawNode *RawNode) CreateViewer(messageQueue tui.MessageQueue) tui.Component {
	data, err := rawNode.getData()
	if err != nil {
		panic("error loading data")
	}

	return NewViewer(messageQueue, rawNode, data)
}

// func (rawNode *RawNode) GetViewer() tea.Model {
// data, err := rawNode.getData()
// if err != nil {
// 	panic("error loading data")
// }

// return NewViewer(rawNode, data)
// }

func (rawNode *RawNode) Serialize(writer io.Writer) error {
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

func (rawNode *RawNode) loadPages(file *multifile.MultiFile) []page.Page {
	pages := []page.Page{}

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

		case "snippet":
			caption, foundCaption := segment.Attributes["caption"]
			if !foundCaption {
				caption = "untitled"
			}

			language, foundLanguage := segment.Attributes["language"]
			if !foundLanguage {
				slog.Warn("Snippet page is missing language attribute, using plaintext", slog.String("filename", file.Path.String()))
				language = "plaintext"
			}
			source := strings.Join(segment.Contents, "\n")
			page := snippetpage.New(caption, source, language)

			pages = append(pages, page)
		}
	}

	return pages
}

func (rawNode *RawNode) parseActions(rawActions []map[string]string) ([]node.Action, error) {
	result := []node.Action{}
	errs := []error{}

	for _, rawAction := range rawActions {
		action, err := rawNode.parseAction(rawAction)

		if err != nil {
			errs = append(errs, err)
		} else {
			result = append(result, action)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return result, nil
}

func (rawNode *RawNode) parseAction(rawAction map[string]string) (node.Action, error) {
	actionType, ok := rawAction["type"]
	if !ok {
		return nil, ErrMissingActionType
	}

	switch actionType {
	case "www":
		return www.Parse(rawAction)

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownActionType, actionType)
	}
}
