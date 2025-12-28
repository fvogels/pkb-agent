package multifile

import (
	"errors"
	"fmt"
	"log/slog"
	"pkb-agent/util"
	"pkb-agent/util/pathlib"
	"strings"
	"unicode"
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
	parser.startSegment()

	for index, line := range lines {
		if index == 0 {
			if err := parser.processFirstLine(line); err != nil {
				return nil, err
			}
		} else if remainder, ok := parser.isDelimiterLine(line); ok {
			parser.finishSegment()
			parser.startSegment()

			if err := parser.processDelimiterLine(remainder); err != nil {
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

func (parser *parser) processFirstLine(line string) error {
	index := strings.Index(line, " ")

	if index == -1 {
		return &ErrInvalidFirstLine{}
	}

	if index == 0 {
		return &ErrInvalidDelimiter{}
	}

	parser.segmentDelimiter = line[:index+1]

	return parser.processDelimiterLine(line[index+1:])
}

func (parser *parser) isDelimiterLine(line string) (string, bool) {
	return strings.CutPrefix(line, parser.segmentDelimiter)
}

// processDelimiterLine expects to receive a string from which the delimiter has already been stripped.
func (parser *parser) processDelimiterLine(str string) error {
	indexOfSpace := strings.Index(str, " ")
	if indexOfSpace == -1 {
		parser.currentSegmentType = str
	} else {
		parser.currentSegmentType = str[:indexOfSpace]

		table, err := parseAttributes(str[indexOfSpace+1:])
		if err != nil {
			return err
		}

		parser.currentSegmentAttributes = table
	}

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

func (parser *parser) startSegment() {
	parser.currentContents = nil
	parser.currentSegmentAttributes = make(map[string]string)
	parser.currentSegmentType = ""
}

var ErrInvalidAttributeString = errors.New("invalid attributes string")

func parseAttributes(str string) (map[string]string, error) {
	result := make(map[string]string)
	runes := util.ConvertToRunes(str)
	index := 0

	for {
		// Skip spaces
		for index < len(runes) && unicode.IsSpace(runes[index].Rune) {
			index++
		}

		// Check if we reached the end of the string
		if index == len(runes) {
			return result, nil
		}

		// Determine key
		keyStartIndex := runes[index].Index
		for runes[index].Rune != '=' {
			index++

			if index == len(runes) {
				slog.Error("Invalid attributes: could not find '='", slog.String("string", str))
				return nil, fmt.Errorf("missing =: %w", ErrInvalidAttributeString)
			}
		}
		key := str[keyStartIndex:runes[index].Index]

		// Check for duplicate keys
		if _, found := result[key]; found {
			slog.Error("Invalid attributes: duplicate key", slog.String("string", str), slog.String("duplicateKey", key))
			return nil, fmt.Errorf("duplicate key %s: %w", key, ErrInvalidAttributeString)
		}

		// Skip =
		index++

		// Something must come after the =
		if index == len(runes) {
			slog.Error("Invalid attributes", slog.String("string", str))
			return nil, fmt.Errorf("missing value: %w", ErrInvalidAttributeString)
		}

		var value string
		// Check if the value is delimited by quotes
		if runes[index].Rune == '"' {
			// Skip "
			index++
			valueStartIndex := runes[index].Index

			if index == len(runes) {
				slog.Error("Invalid attributes", slog.String("string", str))
				return nil, fmt.Errorf("missing value: %w", ErrInvalidAttributeString)
			}

			for runes[index].Rune != '"' {
				index++
				if index == len(runes) {
					slog.Error("Invalid attributes", slog.String("string", str))
					return nil, fmt.Errorf("unexpected end of value: %w", ErrInvalidAttributeString)
				}
			}

			value = str[valueStartIndex:runes[index].Index]

			// Skip "
			index++
		} else {
			valueStartIndex := runes[index].Index

			if index == len(runes) {
				slog.Error("Invalid attributes", slog.String("string", str))
				return nil, fmt.Errorf("missing value: %w", ErrInvalidAttributeString)
			}

			for index < len(runes) && runes[index].Rune != ' ' {
				index++
			}

			if index == len(runes) {
				value = str[valueStartIndex:]
			} else {
				value = str[valueStartIndex:runes[index].Index]
			}
		}

		// Store key/value pair
		result[key] = value
	}
}
