package multifile

type ErrInvalidFirstLine struct{}

func (err *ErrInvalidFirstLine) Error() string {
	return "invalid first line"
}

type ErrInvalidDelimiter struct{}

func (err *ErrInvalidDelimiter) Error() string {
	return "invalid delimiter"
}

type ErrInvalidDelimiterLine struct{}

func (err *ErrInvalidDelimiterLine) Error() string {
	return "invalid delimiter line"
}

type ErrInvalidAttribute struct{}

func (err *ErrInvalidAttribute) Error() string {
	return "invalid segment attribute"
}

type ErrDuplicateAttributeKeys struct{}

func (err *ErrDuplicateAttributeKeys) Error() string {
	return "multiple attributes with same key"
}
