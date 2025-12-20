package multifile

type MultiFile struct {
	Segments []*Segment
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
