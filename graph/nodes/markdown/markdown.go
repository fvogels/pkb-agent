package markdown

import (
	pathlib "pkb-agent/util/pathlib"
	"pkb-agent/util/sectionedfile"
	"strings"
)

type Info struct {
	Path pathlib.Path
}

func (info *Info) GetSource() (string, error) {
	file, err := sectionedfile.LoadSectionedFile(info.Path, isDelimiter)
	if err != nil {
		return "", err
	}

	if len(file.Sections) < 3 {
		return "", &ErrMalformed{path: info.Path}
	}

	if len(file.Sections[0].Lines) > 0 {
		return "", &ErrMalformed{path: info.Path}
	}

	return strings.Join(file.Sections[2].Lines, "\n"), nil
}
