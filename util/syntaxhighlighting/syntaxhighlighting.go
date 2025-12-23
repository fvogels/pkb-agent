package syntaxhighlighting

import (
	"bufio"
	"bytes"

	"github.com/alecthomas/chroma/v2/quick"
)

func Highlight(source string, language string) (string, error) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	if err := quick.Highlight(writer, source, language, "terminal16m", "monokai"); err != nil {
		return "", err
	}
	writer.Flush()

	return buffer.String(), nil
}
