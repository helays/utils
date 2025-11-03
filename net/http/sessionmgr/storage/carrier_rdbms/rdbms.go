package carrier_rdbms

import (
	"context"
	"encoding/gob"
	"time"

	"github.com/helays/utils/v2/net/http/sessionmgr"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Instance struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Instance {
	ins := &Instance{
		db: db.Session(&gorm.Session{}),
	}
	return ins
}

// Register 注册结构定义
// 在使用文件作为session引擎的时候，需要将存储session值的结构注册进来。
func (i *Instance) Register(value ...any) {
	if len(value) < 1 {
		return
	}
	for _, v := range value {
		gob.Register(v)
	}
}

func (i *Instance) Close() error {
	return nil
}

// GC 自动gc
func (i *Instance) GC(ctx context.Context) error {
	tx := i.db.WithContext(ctx).Model(sessionmgr.Session{})
	tx.Where(clause.Lte{Column: "expire_time", Value: time.Now()})
	return tx.Delete(nil).Error
}

func (i *Instance) Save(s *sessionmgr.Session) error {
	return i.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}, {Name: "name"}},
		UpdateAll: true,
	}).Create(s).Error
}

func (i *Instance) Get(sessionId, name string) (*sessionmgr.Session, error) {
	s := &sessionmgr.Session{}
	err := i.db.Where(map[string]any{"id": sessionId, "name": name}).Take(s).Error
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (i *Instance) GetAll(sessionId string) ([]*sessionmgr.Session, error) {
	var sessions []*sessionmgr.Session
	err := i.db.Where(map[string]any{"id": sessionId}).Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (i *Instance) Delete(sessionId, name string) error {
	return i.db.Where(map[string]any{"id": sessionId, "name": name}).Delete(sessionmgr.Session{}).Error
}

func (i *Instance) DeleteAll(sessionId string) error {
	return i.db.Where(map[string]any{"id": sessionId}).Delete(sessionmgr.Session{}).Error
}
