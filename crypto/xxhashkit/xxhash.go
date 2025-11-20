package xxhashkit

import (
	"github.com/cespare/xxhash/v2"
	"github.com/helays/utils/v2/tools"
)

func XXHashString(s string) string {
	return tools.Any2string(xxhash.Sum64String(s))
}

func XXHashBytes(b []byte) string {
	return tools.Any2string(xxhash.Sum64(b))
}
