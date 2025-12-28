package hybrid

import (
	"errors"
	"fmt"
	"log/slog"
	"pkb-agent/util/multifile"
	"pkb-agent/util/pathlib"
	"strings"
)

type Info struct {
	Path pathlib.Path
}

type Data struct {
	Pages   []Page
	Actions []Action
}

type Page interface {
	GetAttributes() map[string]string
}

type PageBase struct {
	Type       string
	Attributes map[string]string
}

func (page *PageBase) GetAttributes() map[string]string {
	return page.Attributes
}

type MarkdownPage struct {
	PageBase
	Source string
}

type SnippetPage struct {
	PageBase
	Language string
	Source   string
}

type Action interface {
	GetDescription() string
}

type ActionBase struct {
	Description string
}

func (action *ActionBase) GetDescription() string {
	return action.Description
}

func (info *Info) GetData() (*Data, error) {
	var data Data

	multiFile, err := multifile.Load(info.Path)
	if err != nil {
		return nil, err
	}

	metadataSegment := multiFile.FindSegmentOfType("metadata")
	if metadataSegment == nil {
		panic("should not occur; this should have been caught earlier")
	}

	metadata, err := parseMetadata(metadataSegment.Contents)
	if err != nil {
		slog.Debug(
			"Failed to parse node metadata",
			slog.String("path", info.Path.String()),
		)
		return nil, err
	}

	pages, err := collectPages(multiFile)
	if err != nil {
		return nil, err
	}
	data.Pages = pages

	actions, err := parseActions(metadata.Actions)
	if err != nil {
		return nil, err
	}
	data.Actions = actions

	return &data, nil
}

func collectPages(multiFile *multifile.MultiFile) ([]Page, error) {
	pages := []Page{}
	errs := []error{}

	for _, segment := range multiFile.Segments {
		switch segment.Type {
		case "metadata":
			continue

		case "markdown":
			page := createMarkdownPage(segment)
			pages = append(pages, page)

		case "snippet":
			page := createSnippetPage(segment)
			pages = append(pages, page)

		default:
			slog.Error("Unknown section type", slog.String("segmentType", segment.Type))
			errs = append(errs, fmt.Errorf("Unknown segment %s: %w", segment.Type, ErrUnknownSegment))
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	} else {
		return pages, nil
	}
}

func createMarkdownPage(segment *multifile.Segment) *MarkdownPage {
	source := strings.Join(segment.Contents, "\n")
	page := MarkdownPage{
		PageBase: PageBase{
			Type:       segment.Type,
			Attributes: segment.Attributes,
		},
		Source: source,
	}

	return &page
}

func createSnippetPage(segment *multifile.Segment) *SnippetPage {
	source := strings.Join(segment.Contents, "\n")
	language, found := segment.Attributes["language"]
	if !found {
		slog.Error("Missing language in snippet")
		panic("missing language")
	}

	page := SnippetPage{
		PageBase: PageBase{
			Type:       segment.Type,
			Attributes: segment.Attributes,
		},
		Language: language,
		Source:   source,
	}

	return &page
}

func parseActions(metadata []map[string]string) ([]Action, error) {
	result := []Action{}
	errs := []error{}

	for _, entry := range metadata {
		action, err := parseAction(entry)
		if err != nil {
			errs = append(errs, err)
		} else {
			result = append(result, action)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	} else {
		return result, nil
	}
}

func parseAction(metadata map[string]string) (Action, error) {
	actionType, found := metadata["type"]
	if !found {
		return nil, fmt.Errorf("missing type: %w", &ErrInvalidAction{})
	}

	switch actionType {
	case "www":
		return parseBrowserAction(metadata)

	case "download":
		return parseDownloadAction(metadata)

	default:
		return nil, fmt.Errorf("unrecognized action type %s: %w", actionType, &ErrInvalidAction{})
	}
}
