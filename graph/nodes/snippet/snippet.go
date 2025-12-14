package snippet

import (
	"os"
	pathlib "pkb-agent/util/pathlib"
	"strings"
)

type Extra struct {
	Path pathlib.Path
}

func (data *Extra) GetSource() (string, error) {
	buffer, err := os.ReadFile(data.Path.String())

	if err != nil {
		return "", err
	}

	lineGenerator := strings.Lines(string(buffer))
	lines := []string{}
	foundMetadataSeparator := false
	lineGenerator(func(line string) bool {
		if !foundMetadataSeparator {
			if strings.TrimSpace(line) == "---" {
				foundMetadataSeparator = true
			}
		} else {
			trimmedLine := strings.TrimRight(line, "\n\r ")
			lines = append(lines, trimmedLine)
		}

		return true
	})

	return strings.Join(lines, "\n"), nil
}
