package markdown

import (
	pathlib "pkb-agent/util/pathlib"
	"pkb-agent/util/sectionedfile"
	"strings"
)

type Extra struct {
	Path pathlib.Path
}

func (data *Extra) GetSource() (string, error) {
	file, err := sectionedfile.LoadSectionedFile(data.Path, isDelimiter)
	if err != nil {
		return "", err
	}

	if len(file.Sections) < 3 {
		return "", &ErrMalformed{path: data.Path}
	}

	if len(file.Sections[0].Lines) > 0 {
		return "", &ErrMalformed{path: data.Path}
	}

	return strings.Join(file.Sections[2].Lines, "\n"), nil
}
