package sessionmgr

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

// 创建session ID
func newSessionId() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}
