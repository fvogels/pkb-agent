package snippet

import (
	"bufio"
	"bytes"
	"log/slog"
	"os"
	pathlib "pkb-agent/util/pathlib"
	"slices"
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
)

type Info struct {
	Path                    pathlib.Path
	LanguageForHighlighting string
}

func (data *Info) GetSource() (string, error) {
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

func (data *Info) GetHighlightedSource() (string, string, error) {
	rawSource, err := data.GetSource()
	if err != nil {
		return "", "", err
	}

	if !slices.Contains(lexers.Names(true), data.LanguageForHighlighting) {
		slog.Error(
			"Unsupported language for highlighting",
			slog.String("language", data.LanguageForHighlighting),
		)
	}

	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	if err := quick.Highlight(writer, rawSource, data.LanguageForHighlighting, "terminal16m", "monokai"); err != nil {
		return "", "", err
	}
	writer.Flush()

	return rawSource, buffer.String(), nil
}
