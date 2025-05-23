package gomutil

import (
	. "maragu.dev/gomponents"
)

func MergeHeads(nodes ...[]Node) []Node {
	var totalLen int
	for _, s := range nodes {
		totalLen += len(s)
	}

	merged := make([]Node, 0, totalLen)
	for _, s := range nodes {
		merged = append(merged, s...)
	}
	return merged
}
