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
	MarkdownSource string
	Actions        []Action
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

	sectionedFile, err := multifile.Load(info.Path)
	if err != nil {
		return nil, err
	}

	metadataSegment := sectionedFile.FindSegmentOfType("metadata")
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

	actions, err := parseActions(metadata.Actions)
	if err != nil {
		return nil, err
	}
	data.Actions = actions

	if markdownSegment := sectionedFile.FindSegmentOfType("markdown"); markdownSegment != nil {
		data.MarkdownSource = strings.Join(markdownSegment.Contents, "\n")
	}

	return &data, nil
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
