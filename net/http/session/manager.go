package session

import (
	"context"
	"encoding/gob"
	"sync"
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
	Id         string              `json:"id" gorm:"primaryKey;autoIncrement:false;type:varchar(64);not null;index;comment:Session ID"`
	Name       string              `json:"name" gorm:"primaryKey;autoIncrement:false;type:varchar(128);not null;index;comment:Session的名字"`
	Values     SessionValue        `json:"values" gorm:"comment:session数据"`
	CreateTime dataType.CustomTime `json:"create_time" gorm:"comment:session 创建时间"`
	ExpireTime dataType.CustomTime `json:"expire_time" gorm:"not null;index;comment:session 过期时间"`
	Duration   time.Duration       `json:"duration" gorm:"comment:session有效期"`
}

type Manager struct {
	options *Options
	storage StorageDriver
}

// StorageDriver 存储驱动接口 - 只负责存储，不处理业务逻辑
type StorageDriver interface {
	Save(session *Session) error
	Get(sessionId, name string) (*Session, error)
	GetAll(sessionId string) ([]*Session, error)
	Delete(sessionId, name string) error
	DeleteAll(sessionId string) error

	// GC 相关
	GC(ctx context.Context) error
	Close() error

	Register(value ...any) // 注册结构定义
}

// noinspection all
func New(ctx context.Context, storage StorageDriver, opt ...*Options) *Manager {
	options := &Options{}
	if len(opt) > 0 {
		options = opt[0]
	}
	options.Cookie.Name = tools.Ternary(options.Cookie.Name == "", CookieName, options.Cookie.Name)
	options.CheckInterval = tools.Ternary(options.CheckInterval <= 0, Interval, options.CheckInterval)
	options.Carrier = tools.Ternary(options.Carrier == "", "cookie", options.Carrier)
	options.GcProbability = tools.Ternary(options.GcProbability <= 0, 0.9, options.GcProbability) // 默认GC 90%

	manager := &Manager{
		options: options,
		storage: storage,
	}
	gob.Register(SessionValue{})
	if !options.DisableGc {
		go manager.startGC(ctx)
	}
	return manager
}

func (m *Manager) startGC(ctx context.Context) {
	tools.RunAsyncTickerProbabilityWithContext(ctx, true, m.options.CheckInterval, m.options.GcProbability, func(ctx context.Context) {
		ulogs.CheckErrf(m.storage.GC(ctx), "session gc 失败")
	})
}

var (
	session *Manager
	once    sync.Once
)

func StartSession(ctx context.Context, storage StorageDriver, opt ...*Options) {
	once.Do(func() {
		session = New(ctx, storage, opt...)
		ulogs.Infof("session 启动成功")
	})

}

func GetSession() *Manager {
	if session == nil {
		panic("session 未初始化")
	}
	return session
}
