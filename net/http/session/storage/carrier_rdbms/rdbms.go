package carrier_rdbms

import (
	"context"
	"time"

	"github.com/helays/utils/v2/db/userDb"
	"github.com/helays/utils/v2/net/http/session"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Instance struct {
	db *gorm.DB
}

// New 创建一个session存储实例
func New(db *gorm.DB) *Instance {
	ins := &Instance{
		db: db.Session(&gorm.Session{}),
	}
	userDb.AutoCreateTableWithStruct(ins.db, session.Session{}, "创建 session 表失败")
	return ins
}

func (i *Instance) Close() error {
	return nil
}

// GC 自动gc
func (i *Instance) GC(ctx context.Context) error {
	tx := i.db.WithContext(ctx).Model(session.Session{})
	tx.Where(clause.Lte{Column: "expire_time", Value: time.Now()})
	return tx.Delete(nil).Error
}

func (i *Instance) Save(s *session.Session) error {
	return i.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}, {Name: "name"}},
		UpdateAll: true,
	}).Create(s).Error
}

func (i *Instance) Get(sessionId, name string) (*session.Session, error) {
	s := &session.Session{}
	err := i.db.Where(map[string]any{"id": sessionId, "name": name}).Take(s).Error
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (i *Instance) GetAll(sessionId string) ([]*session.Session, error) {
	var sessions []*session.Session
	err := i.db.Where(map[string]any{"id": sessionId}).Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (i *Instance) Delete(sessionId, name string) error {
	return i.db.Where(map[string]any{"id": sessionId, "name": name}).Delete(session.Session{}).Error
}

func (i *Instance) DeleteAll(sessionId string) error {
	return i.db.Where(map[string]any{"id": sessionId}).Delete(session.Session{}).Error
}
