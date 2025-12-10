package carrier_memory

import (
	"context"
	"strings"
	"time"

	"github.com/helays/utils/v2/map/safettl"
	"github.com/helays/utils/v2/net/http/session"
)

type Instance struct {
	storage *safettl.PerKeyTTLMap[string, *session.Session]
}

func New() *Instance {
	i := &Instance{}
	i.storage = safettl.NewPerKeyTTLMapWithInterval[string, *session.Session](time.Minute) // session 回收策略不需要太高，1分钟
	return i
}

func (i *Instance) GC(_ context.Context) error {
	return nil
}

func (i *Instance) uniqueId(id, name string) string {
	return id + "_" + name
}

func (i *Instance) Save(s *session.Session) error {
	i.storage.StoreWithTTL(i.uniqueId(s.Id, s.Name), s, s.Duration)
	return nil
}

func (i *Instance) Get(sessionId, name string) (*session.Session, error) {
	s, ok := i.storage.Load(i.uniqueId(sessionId, name))
	if !ok {
		return nil, session.ErrNotFound
	}
	return s, nil
}

func (i *Instance) GetAll(sessionId string) ([]*session.Session, error) {
	var sessions []*session.Session
	i.storage.Range(func(key string, value *session.Session) bool {
		if strings.HasPrefix(key, sessionId) {
			sessions = append(sessions, value)
		}
		return true
	})
	return sessions, nil
}

func (i *Instance) Delete(sessionId, name string) error {
	i.storage.Delete(i.uniqueId(sessionId, name))
	return nil
}

func (i *Instance) DeleteAll(sessionId string) error {
	i.storage.DeletePrefix(sessionId)
	return nil
}

func (i *Instance) Close() error {
	return nil
}
