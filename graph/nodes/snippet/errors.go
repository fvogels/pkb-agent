package snippet

type ErrMissingSnippet struct{}

func (err *ErrMissingSnippet) Error() string {
	return "only metadata in snippet"
}
