package hybrid

import (
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
	ExternalLinks  []ExternalLink
}

type ExternalLink struct {
	URL         string
	Description string
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

	data.ExternalLinks = metadata.ExternalLinks

	if markdownSegment := sectionedFile.FindSegmentOfType("markdown"); markdownSegment != nil {
		data.MarkdownSource = strings.Join(markdownSegment.Contents, "\n")
	}

	return &data, nil
}
