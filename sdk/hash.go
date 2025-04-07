package sdk

import (
	"hash/fnv"
)

func hashToPercent(input string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(input))
	return int(h.Sum32() % 100)
}
