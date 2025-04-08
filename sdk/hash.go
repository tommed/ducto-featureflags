package sdk

import (
	"crypto/sha256"
	"hash/fnv"
	"os"
	"sync"
)

var (
	cachedHostname string
	hostnameOnce   sync.Once
)

func getHostname() string {
	hostnameOnce.Do(func() {
		hostname, err := os.Hostname()
		if err == nil {
			cachedHostname = hostname
		}
	})
	return cachedHostname
}

func hashToPercent(value, algo string) int {
	switch algo {
	case "sha256":
		h := sha256.Sum256([]byte(value))
		return int(h[0]) % 100
	default: // fallback: FNV
		h := fnv.New32a()
		_, _ = h.Write([]byte(value))
		return int(h.Sum32() % 100)
	}
}
