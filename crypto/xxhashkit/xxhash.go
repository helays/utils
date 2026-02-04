package xxhashkit

import (
	"github.com/cespare/xxhash/v2"
	"helay.net/go/utils/v3/tools"
)

func XXHashString(s string) string {
	return tools.Any2string(xxhash.Sum64String(s))
}

func XXHashBytes(b []byte) string {
	return tools.Any2string(xxhash.Sum64(b))
}
