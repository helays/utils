package gin_adapter

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/helays/utils/v2/net/http/session"
)

type GinManager struct {
	manager *session.Manager
}

func New(manager *session.Manager) *GinManager {
	return &GinManager{manager: manager}
}

func (m *GinManager) GetSessionId(c *gin.Context) (string, error) {
	return m.manager.GetSessionId(c.Writer, c.Request)
}

func (m *GinManager) Get(c *gin.Context, name string, dst any) error {
	return m.manager.Get(c.Writer, c.Request, name, dst)
}

func (m *GinManager) GetUp(c *gin.Context, name string, dst any) error {
	return m.manager.GetUp(c.Writer, c.Request, name, dst)
}

func (m *GinManager) GetUpByRemainTime(c *gin.Context, name string, dst any, duration time.Duration) error {
	return m.manager.GetUpByRemainTime(c.Writer, c.Request, name, dst, duration)
}

func (m *GinManager) GetUpByDuration(c *gin.Context, name string, dst any, duration time.Duration) error {
	return m.manager.GetUpByDuration(c.Writer, c.Request, name, dst, duration)
}

func (m *GinManager) Flashes(c *gin.Context, name string, dst any) error {
	return m.manager.Flashes(c.Writer, c.Request, name, dst)
}

func (m *GinManager) Set(c *gin.Context, value *session.Value) error {
	return m.manager.Set(c.Writer, c.Request, value)
}

func (m *GinManager) Del(c *gin.Context, name string) error {
	return m.manager.Del(c.Writer, c.Request, name)
}

func (m *GinManager) Destroy(c *gin.Context) error {
	return m.manager.Destroy(c.Writer, c.Request)
}

func (m *GinManager) Close() error {
	return m.manager.Close()
}
