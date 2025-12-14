package snippet

import (
	"bufio"
	"bytes"
	"os"
	pathlib "pkb-agent/util/pathlib"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
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

func (data *Extra) GetHighlightedSource() (string, error) {
	rawSource, err := data.GetSource()
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	if err := quick.Highlight(writer, rawSource, "go", "terminal16m", "monokai"); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
