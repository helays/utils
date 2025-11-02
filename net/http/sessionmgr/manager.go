package sessionmgr

import (
	"context"
	"time"

	"github.com/helays/utils/v2/dataType"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/tools"
)

const (
	Interval   = time.Hour // 默认检测频率
	CookieName = "vsclubId"
	ExpireTime = 1 * time.Hour // session默认24小时过期
)

// Session session 数据结构
type Session struct {
	Id         string              `json:"id" gorm:"primaryKey;autoIncrement:false;type:varchar(64);not null;index;comment:Session ID"`       // session id
	Name       string              `json:"name" gorm:"primaryKey;autoIncrement:false;type:varchar(128);not null;index;comment:Session的名字"` // session 的名字
	Values     SessionValue        `json:"values" gorm:"comment:session数据"`                                                                 // session 数据
	CreateTime dataType.CustomTime `json:"create_time" gorm:"comment:session 创建时间"`                                                       // 创建时间 ，没啥用，就看
	ExpireTime dataType.CustomTime `json:"expire_time" gorm:"not null;index;comment:session 过期时间"`                                        // 过期时间 ，用于自动回收的时候使用
	Duration   time.Duration       `json:"duration" gorm:"comment:session有效期"`                                                             // 有效期，主要是用于更新有效期的时候使用
}
type Manager struct {
	options *Options
	storage StorageDriver
}

// StorageDriver 存储驱动接口 - 只负责存储，不处理业务逻辑
type StorageDriver interface {
	Save(session *Session) error
	Get(sessionID, name string) (*Session, error)
	GetAll(sessionID string) ([]*Session, error)
	Delete(sessionID, name string) error
	DeleteAll(sessionID string) error

	// GC 相关
	GC() error
	Close() error
}

func New(ctx context.Context, storage StorageDriver, opt ...*Options) *Manager {
	options := &Options{}
	if len(opt) > 0 {
		options = opt[0]
	}
	options.CookieName = tools.Ternary(options.CookieName == "", CookieName, options.CookieName)
	options.CheckInterval = tools.Ternary(options.CheckInterval <= 0, Interval, options.CheckInterval)
	options.Carrier = tools.Ternary(options.Carrier == "", "cookie", options.Carrier)
	options.GcProbability = tools.Ternary(options.GcProbability <= 0, 0.9, options.GcProbability) // 默认GC 90%

	manager := &Manager{
		options: options,
		storage: storage,
	}
	if !options.DisableGc {
		go manager.startGC(ctx)
	}
	return manager
}

func (m *Manager) startGC(ctx context.Context) {
	tools.RunAsyncTickerProbabilityWithContext(ctx, true, m.options.CheckInterval, m.options.GcProbability, func(ctx context.Context) {
		ulogs.CheckErrf(m.storage.GC(), "session gc 失败")
	})
}
