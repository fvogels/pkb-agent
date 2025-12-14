package sectionedfile

import (
	"os"
	"pkb-agent/util/pathlib"
	"strings"
)

type SectionedFile struct {
	Sections []Section
}

type Section struct {
	Lines []string
}

func LoadSectionedFile(path pathlib.Path, delimiterPredicate func(string) bool) (SectionedFile, error) {
	sections := []Section{}
	currentSection := []string{}

	bytes, err := os.ReadFile(path.String())
	if err != nil {
		return SectionedFile{}, err
	}

	lineGenerator := strings.Lines(string(bytes))

	lineGenerator(func(line string) bool {
		if delimiterPredicate(line) {
			sections = append(sections, Section{Lines: currentSection})
			currentSection = []string{}
		}

		return true
	})

	sections = append(sections, Section{Lines: currentSection})

	return SectionedFile{Sections: sections}, nil
}
