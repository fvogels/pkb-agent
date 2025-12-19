package hybrid

import (
	"log/slog"
	"pkb-agent/util/pathlib"
	"pkb-agent/util/sectionedfile"
	"strings"
)

type Extra struct {
	Path pathlib.Path
}

type Data struct {
	MarkdownSource string
	URLs           []string
}

func (extra *Extra) GetData() (*Data, error) {
	var data Data

	sectionedFile, err := sectionedfile.LoadSectionedFile(extra.Path, isDelimiter)
	if err != nil {
		return nil, err
	}

	if len(sectionedFile.Sections) == 0 {
		slog.Debug(
			"Failed to load sectioned file",
			slog.String("path", extra.Path.String()),
		)
	}

	metadata, err := parseMetadata(sectionedFile.Sections[0].Lines)
	if err != nil {
		slog.Debug(
			"Failed to parse node metadata",
			slog.String("path", extra.Path.String()),
		)
		return nil, err
	}

	data.URLs = metadata.URLs

	if len(sectionedFile.Sections) > 1 {
		data.MarkdownSource = strings.Join(sectionedFile.Sections[1].Lines, "\n")
	}

	return &data, nil
}
