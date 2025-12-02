package session

import (
	"encoding/hex"

	"github.com/google/uuid"
)

// 创建session ID
func newSessionId() string {
	u := uuid.New()
	return hex.EncodeToString(u[:])
}
