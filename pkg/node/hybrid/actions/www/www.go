package www

import (
	"fmt"
	"pkb-agent/extern"
)

type Action struct {
	description string
	url         string
	key         string
}

func New(description string, url string) *Action {
	return &Action{
		description: description,
		url:         url,
	}
}

func (action Action) GetDescription() string {
	return action.description
}

func (action Action) GetKey() string {
	return action.key
}

func (action Action) Perform() {
	extern.OpenURLInBrowser(action.url)
}

func Parse(rawAction map[string]string) (*Action, error) {
	description, err := getAttribute(rawAction, "description")
	if err != nil {
		return nil, err
	}

	url, err := getAttribute(rawAction, "url")
	if err != nil {
		return nil, err
	}

	key, err := getAttribute(rawAction, "key")
	if err != nil {
		return nil, err
	}

	action := Action{
		description: description,
		url:         url,
		key:         key,
	}

	return &action, nil
}

func getAttribute(rawAction map[string]string, fieldName string) (string, error) {
	value, ok := rawAction[fieldName]

	if !ok {
		return "", fmt.Errorf("%w: %s", ErrMissingAttribute, fieldName)
	}

	return value, nil
}
