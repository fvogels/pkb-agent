package multifile

import "pkb-agent/util/pathlib"

type MultiFile struct {
	Segments []*Segment
	Path     pathlib.Path
}

type Segment struct {
	Type       string
	Attributes map[string]string
	Contents   []string
}

func (file *MultiFile) FindSegmentOfType(segmentType string) *Segment {
	for _, segment := range file.Segments {
		if segment.Type == segmentType {
			return segment
		}
	}

	return nil
}
