package session

import (
	"encoding/hex"

	"github.com/google/uuid"
)

// GenerateSessionID 创建session ID
func GenerateSessionID() string {
	u := uuid.New()
	return hex.EncodeToString(u[:])
}
