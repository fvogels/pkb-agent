package multifile

import (
	"fmt"
	"pkb-agent/util/pathlib"
	"strings"
)

func Load(path pathlib.Path) (*MultiFile, error) {
	lines, err := path.ReadLines()
	if err != nil {
		return nil, err
	}

	return Parse(lines)
}

func Parse(contents []string) (*MultiFile, error) {
	parser := parser{}

	return parser.parse(contents)
}

type parser struct {
	segmentDelimiter         string // line prefix that denotes the start of a new segment; always contains a space at the end
	segments                 []*Segment
	currentSegmentType       string
	currentSegmentAttributes map[string]string
	currentContents          []string
}

func (parser *parser) parse(lines []string) (*MultiFile, error) {
	for index, line := range lines {
		if index == 0 {
			if err := parser.parseFirstLine(line); err != nil {
				return nil, err
			}
		} else if remainder, ok := parser.isDelimiterLine(line); ok {
			parser.finishSegment()

			if err := parser.parseDelimiterLine(remainder); err != nil {
				return nil, err
			}
		} else {
			parser.currentContents = append(parser.currentContents, line)
		}
	}

	parser.finishSegment()

	result := MultiFile{
		Segments: parser.segments,
	}

	return &result, nil
}

func (parser *parser) parseFirstLine(line string) error {
	index := strings.Index(line, " ")

	if index == -1 {
		return &ErrInvalidFirstLine{}
	}

	if index == 0 {
		return &ErrInvalidDelimiter{}
	}

	parser.segmentDelimiter = line[:index+1]

	return parser.parseDelimiterLine(line[index+1:])
}

func (parser *parser) isDelimiterLine(line string) (string, bool) {
	return strings.CutPrefix(line, parser.segmentDelimiter)
}

// parseDelimiterLine expects to receive a string from which the delimiter has already been stripped.
func (parser *parser) parseDelimiterLine(str string) error {
	parts := strings.Split(str, " ")

	if len(parts) == 0 {
		return fmt.Errorf("no segment type specified: %w", &ErrInvalidDelimiterLine{})
	}

	parser.currentSegmentType = parts[0]

	for _, part := range parts[1:] {
		if err := parser.parseAttribute(part); err != nil {
			return err
		}
	}

	return nil
}

func (parser *parser) parseAttribute(str string) error {
	parts := strings.SplitN(str, "=", 2)

	if len(parts) != 2 {
		return &ErrInvalidAttribute{}
	}

	key := parts[0]
	value := parts[1]

	if _, found := parser.currentSegmentAttributes[key]; found {
		return &ErrDuplicateAttributeKeys{}
	}

	parser.currentSegmentAttributes[key] = value

	return nil
}

func (parser *parser) finishSegment() {
	segment := Segment{
		Type:       parser.currentSegmentType,
		Attributes: parser.currentSegmentAttributes,
		Contents:   parser.currentContents,
	}

	parser.segments = append(parser.segments, &segment)
}
