package carrier_memory

import (
	"strings"
	"time"

	"github.com/helays/utils/v2/map/syncMapWrapper"
	"github.com/helays/utils/v2/net/http/sessionmgr"
)

type Instance struct {
	storage syncMapWrapper.SyncMap[string, *sessionmgr.Session]
}

func New() *Instance {
	i := &Instance{}
	return i
}

func (i *Instance) Gc() error {
	i.storage.Range(func(key string, value *sessionmgr.Session) bool {
		if time.Time(value.ExpireTime).Before(time.Now()) {
			i.storage.Delete(key)
		}
		return true
	})
	return nil
}

func (i *Instance) uniqueId(id, name string) string {
	return id + "_" + name
}

func (i *Instance) Save(s *sessionmgr.Session) error {
	i.storage.Store(i.uniqueId(s.Id, s.Name), s)
	return nil
}

func (i *Instance) Get(sessionID, name string) (*sessionmgr.Session, error) {
	s, ok := i.storage.Load(i.uniqueId(sessionID, name))
	if !ok {
		return nil, sessionmgr.ErrNotFound
	}
	return s, nil
}

func (i *Instance) GetAll(sessionID string) ([]*sessionmgr.Session, error) {
	var sessions []*sessionmgr.Session
	i.storage.Range(func(key string, value *sessionmgr.Session) bool {
		if strings.HasPrefix(key, sessionID) {
			sessions = append(sessions, value)
		}
		return true
	})
	return sessions, nil
}

func (i *Instance) Delete(sessionID, name string) error {
	i.storage.Delete(i.uniqueId(sessionID, name))
	return nil
}

func (i *Instance) DeleteAll(sessionID string) error {
	i.storage.DeletePrefix(sessionID)
	return nil
}

func (i *Instance) Close() error {
	return nil
}
