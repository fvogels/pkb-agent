package www

import "pkb-agent/extern"

type Action struct {
	description string
	url         string
}

func (action Action) GetDescription() string {
	return action.description
}

func (action Action) Perform() {
	extern.OpenURLInBrowser(action.url)
}

func Parse(rawAction map[string]string) (*Action, error) {
	description, ok := rawAction["description"]
	if !ok {
		return nil, ErrMissingDescription
	}

	url, ok := rawAction["url"]
	if !ok {
		return nil, ErrMissingURL
	}

	action := Action{
		description: description,
		url:         url,
	}

	return &action, nil
}
