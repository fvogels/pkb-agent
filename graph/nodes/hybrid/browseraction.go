package hybrid

import "fmt"

type BrowserAction struct {
	ActionBase
	URL string
}

func parseBrowserAction(metadata map[string]string) (*BrowserAction, error) {
	description, found := metadata["description"]
	if !found {
		return nil, fmt.Errorf("missing description: %w", ErrInvalidAction{})
	}

	url, found := metadata["url"]
	if !found {
		return nil, fmt.Errorf("missing url: %w", ErrInvalidAction{})
	}

	action := BrowserAction{
		ActionBase: ActionBase{
			Description: description,
		},
		URL: url,
	}

	return &action, nil
}
