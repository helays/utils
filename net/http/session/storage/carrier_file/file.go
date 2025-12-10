package carrier_file

import (
	"context"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/helays/utils/v2/close/vclose"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/net/http/session"
	"github.com/helays/utils/v2/tools"
)

type Instance struct {
	path string
}

func New(path ...string) (*Instance, error) {
	ins := &Instance{
		path: "runtime/session",
	}
	if len(path) > 0 {
		ins.path = path[0]
	}
	ins.path = tools.Fileabs(ins.path)
	if err := tools.Mkdir(ins.path); err != nil {
		return nil, fmt.Errorf("创建session文件存放目录失败 %v", err)
	}
	return ins, nil
}

func (i *Instance) Close() error {
	return nil
}

func (i *Instance) GC(_ context.Context) error {
	sessionIdFiles, err := os.ReadDir(i.path)
	if err != nil {
		return err
	}
	for _, files := range sessionIdFiles {
		if !files.IsDir() {
			continue
		}
		sessionFiles, err := os.ReadDir(filepath.Join(i.path, files.Name()))
		if err != nil {
			ulogs.Errorf("session文件夹读取失败 %v", err)
			continue
		}
		for _, file := range sessionFiles {
			sessionFile := filepath.Join(i.path, files.Name(), file.Name())
			if file.IsDir() {
				i.del(sessionFile)
				continue
			}
			i.expireCheck(sessionFile)
		}
	}
	return err
}

// 检查session文件是否过期
func (i *Instance) expireCheck(sessionFile string) {
	s := &session.Session{}
	f, err := os.Open(sessionFile)
	defer vclose.Close(f)
	if err != nil {
		ulogs.Errorf("打开session文件[%s]失败 %v", sessionFile, err)
		return
	}
	if err = gob.NewDecoder(f).Decode(s); err != nil {
		vclose.Close(f)
		ulogs.Errorf("session文件[%s]解析失败 %v", sessionFile, err)
		return
	}
	if time.Time(s.ExpireTime).Before(time.Now()) {
		i.del(sessionFile)
	}
}

func (i *Instance) del(file string) {
	ulogs.CheckErrf(os.RemoveAll(file), "session文件删除失败 %v", file)
}

func (i *Instance) Save(s *session.Session) error {
	sessionFile := filepath.Join(i.path, s.Id, s.Name)
	f, err := os.OpenFile(sessionFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	defer vclose.Close(f)
	if err != nil {
		return err
	}
	return gob.NewEncoder(f).Encode(s)
}

func (i *Instance) Get(sessionId, name string) (*session.Session, error) {
	sessionFile := filepath.Join(i.path, sessionId, name)
	f, err := os.Open(sessionFile)
	defer vclose.Close(f)
	if err != nil {
		return nil, err
	}
	s := &session.Session{}
	if err = gob.NewDecoder(f).Decode(s); err != nil {
		return nil, err
	}
	if time.Time(s.ExpireTime).Before(time.Now()) {
		i.del(sessionFile)
		return nil, session.ErrNotFound
	}
	return s, nil
}

func (i *Instance) GetAll(sessionId string) ([]*session.Session, error) {
	sessionFile := filepath.Join(i.path, sessionId)
	sessionFiles, err := os.ReadDir(sessionFile)
	if err != nil {
		return nil, err
	}
	var sessions []*session.Session
	for _, file := range sessionFiles {
		if file.IsDir() {
			continue
		}
		if s, _e := i.Get(sessionId, file.Name()); _e == nil {
			sessions = append(sessions, s)
		} else {
			ulogs.Errorf("批量session文件[%s/%s]解析失败 %v", sessionId, file.Name(), _e)
		}
	}
	return sessions, nil
}

func (i *Instance) Delete(sessionId, name string) error {
	sessionFile := filepath.Join(i.path, sessionId, name)
	return os.RemoveAll(sessionFile)
}

func (i *Instance) DeleteAll(sessionId string) error {
	sessionFile := filepath.Join(i.path, sessionId)
	return os.RemoveAll(sessionFile)
}
