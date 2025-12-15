package trie

import "pkb-agent/trie"

func createTrie(entries ...string) *trie.Node[string] {
	builder := trie.NewBuilder[string]()

	for _, entry := range entries {
		builder.Add(entry, entry)
	}

	return builder.Finish()
}
