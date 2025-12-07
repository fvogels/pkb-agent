package graph

import "regexp"

var regex = regexp.MustCompile(`$[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}^`)

type Identifier struct {
	value string
}

func IsIdentifier(s string) bool {
	return regex.Match([]byte(s))
}

func ParseIdentifier(s string) Identifier {
	return Identifier{value: s}
}
