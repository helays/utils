package carrier_redis

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/net/http/session"
	"github.com/redis/go-redis/v9"
)

type Instance struct {
	rdb redis.UniversalClient
}

func New(rdb redis.UniversalClient) *Instance {
	return &Instance{
		rdb: rdb,
	}
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

// GC redis 引擎，可以使用ttl作为超时自动回收
func (i *Instance) GC(_ context.Context) error {
	return nil
}

func (i *Instance) uniqueId(id, name string) string {
	return id + "_" + name
}

func (i *Instance) Save(s *session.Session) error {
	rk := i.uniqueId(s.Id, s.Name)
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(s); err != nil {
		return err
	}
	// 需要将key单独存储，方便后续按sessionId删除
	ctx := context.Background()
	if err := i.rdb.HSet(ctx, s.Id, rk, rk).Err(); err != nil {
		return err
	}
	return i.rdb.Set(ctx, rk, buf.Bytes(), s.Duration).Err()
}

// 删除session
func (i *Instance) del(sessionId string, names ...string) {
	ctx := context.Background()
	if len(names) > 0 {
		for _, name := range names {
			rk := i.uniqueId(sessionId, name)
			i.rdb.HDel(ctx, sessionId, rk)
			i.rdb.Del(ctx, rk)
		}
		return
	}
	// 先获取到所有sessionId下面的key
	res := i.rdb.HGetAll(ctx, sessionId)
	m, err := res.Result()
	if err != nil {
		ulogs.Errorf("redis session 获取sessionId:%s下的所有key失败 %v", sessionId, err)
		return
	}
	i.rdb.Del(ctx, sessionId)
	for k, _ := range m {
		rk := i.uniqueId(sessionId, k)
		i.rdb.Del(ctx, rk)
	}
}

// Get 获取session
func (i *Instance) Get(sessionId, name string) (*session.Session, error) {
	rk := i.uniqueId(sessionId, name)
	ctx := context.Background()
	val, err := i.rdb.Get(ctx, rk).Bytes()
	if err != nil {
		return nil, err
	}
	s := &session.Session{}
	if err = gob.NewDecoder(bytes.NewReader(val)).Decode(s); err != nil {
		return nil, err
	}
	return s, nil
}

func (i *Instance) GetAll(sessionId string) ([]*session.Session, error) {
	ctx := context.Background()
	// 先获取到所有sessionId下面的key
	res := i.rdb.HGetAll(ctx, sessionId)
	m, err := res.Result()
	if err != nil {
		return nil, err
	}
	var sessions []*session.Session
	for k, _ := range m {
		if s, _err := i.Get(sessionId, k); _err != nil {
			ulogs.Errorf("redis session 获取sessionId:%s下的key:%s失败 %v", sessionId, k, _err)
		} else {
			sessions = append(sessions, s)
		}
	}
	return sessions, nil
}

func (i *Instance) Delete(sessionId, name string) error {
	i.del(sessionId, name)
	return nil
}

func (i *Instance) DeleteAll(sessionId string) error {
	i.del(sessionId)
	return nil
}
