package carrierRedis

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/helays/utils/v2/dataType"
	session2 "github.com/helays/utils/v2/net/http/session"
	"github.com/redis/go-redis/v9"
)

//
// ━━━━━━神兽出没━━━━━━
// 　　 ┏┓     ┏┓
// 　　┏┛┻━━━━━┛┻┓
// 　　┃　　　　　 ┃
// 　　┃　　━　　　┃
// 　　┃　┳┛　┗┳  ┃
// 　　┃　　　　　 ┃
// 　　┃　　┻　　　┃
// 　　┃　　　　　 ┃
// 　　┗━┓　　　┏━┛　Code is far away from bug with the animal protecting
// 　　　 ┃　　　┃    神兽保佑,代码无bug
// 　　　　┃　　　┃
// 　　　　┃　　　┗━━━┓
// 　　　　┃　　　　　　┣┓
// 　　　　┃　　　　　　┏┛
// 　　　　┗┓┓┏━┳┓┏┛
// 　　　　 ┃┫┫ ┃┫┫
// 　　　　 ┗┻┛ ┗┻┛
//
// ━━━━━━感觉萌萌哒━━━━━━
//
//
// User helay
// Date: 2024/12/8 17:27
//

// redis 引擎，可以使用ttl作为超时自动回收
// 无需手动调用删除，session超时后会自动删除

// redis.UniversalClient

// Instance session 实例
type Instance struct {
	option *session2.Options
	rdb    redis.UniversalClient
}

// New 初始化 session 内存 实例
func New(rdb redis.UniversalClient) *Instance {
	ins := &Instance{
		rdb: rdb,
	}
	return ins
}

// Apply 应用配置
func (this *Instance) Apply(options *session2.Options) {
	this.option = options
}

// Register 注册结构定义
// 在使用文件作为session引擎的时候，需要将存储session值的结构注册进来。
func (this *Instance) Register(value ...any) {
	if len(value) < 1 {
		return
	}
	for _, v := range value {
		gob.Register(v)
	}
}

// Close 关闭 db
func (this *Instance) Close() error {
	return nil
}

// 从redis中读取session 数据
func (this *Instance) get(w http.ResponseWriter, r *http.Request, name string) (*session2.Session, string, error) {
	sessionId, err := session2.GetSessionId(w, r, this.option) // 这一步一般不会失败
	if err != nil {
		return nil, "", err // 从cookie中获取sessionId失败
	}
	_k := session2.GetSessionName(sessionId, name)

	val, err := this.rdb.Get(context.Background(), _k).Bytes()
	if err != nil {
		return nil, "", err
	}
	sessionVal := &session2.Session{}
	if err = gob.NewDecoder(bytes.NewReader(val)).Decode(sessionVal); err != nil {
		// session 数据解析失败，删除session文件
		this.del(sessionId, name)
		return nil, "", err
	}
	if time.Time(sessionVal.ExpireTime).Before(time.Now()) {
		this.del(sessionId, name)
		// session已过期
		return nil, "", session2.ErrNotFound
	}
	return sessionVal, _k, nil
}

// 设置session，将其通过gob 写入文件
func (this *Instance) set(w http.ResponseWriter, r *http.Request, dst session2.Session) error {
	_k := session2.GetSessionName(dst.Id, dst.Name)

	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(dst)
	if err != nil {
		return err
	}
	// 将key单独存储，方便后续删除
	this.rdb.HSet(context.Background(), dst.Id, _k, _k)
	return this.rdb.Set(context.Background(), _k, buf.Bytes(), dst.Duration).Err()
}

// 从redis删除 session 数据
func (this *Instance) del(sessionId, _k string) {
	this.rdb.Del(context.Background(), _k)
	this.rdb.HDel(context.Background(), sessionId, _k)
}

// Get 获取session
func (this *Instance) Get(w http.ResponseWriter, r *http.Request, name string, dst any) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("dst must be a pointer")
	}
	sessionVal, _, err := this.get(w, r, name)
	if err != nil {
		return err
	}
	v.Elem().Set(reflect.ValueOf(sessionVal.Values.Val))
	return nil
}

// GetUp 获取session并更新过期时间
func (this *Instance) GetUp(w http.ResponseWriter, r *http.Request, name string, dst any) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("dst must be a pointer")
	}
	sessionVal, _, err := this.get(w, r, name)
	if err != nil {
		return err
	}
	// 更新session过期时间
	sessionVal.ExpireTime = dataType.CustomTime(time.Now().Add(sessionVal.Duration))
	if err = this.set(w, r, *sessionVal); err != nil {
		return err
	}
	v.Elem().Set(reflect.ValueOf(sessionVal.Values.Val))
	return nil
}

// GetUpByTimeLeft 根据剩余时间更新session
// 当session 的有效期小于duration，那么将session的有效期延长到 session.Duration-duration
// 比如：设置了15天有效期，duration设置一天，那么当检测到session的有效期 不大于一天的时候就更新session
func (this *Instance) GetUpByTimeLeft(w http.ResponseWriter, r *http.Request, name string, dst any, duration time.Duration) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("dst must be a pointer")
	}
	sessionVal, _, err := this.get(w, r, name)
	if err != nil {
		return err
	}
	v.Elem().Set(reflect.ValueOf(sessionVal.Values.Val))
	// 判断 距离过期时间小于等于duration 的时候，更新session的过期时间
	if time.Time(sessionVal.ExpireTime).Sub(time.Now()) <= duration {
		sessionVal.ExpireTime = dataType.CustomTime(time.Now().Add(sessionVal.Duration))
		return this.set(w, r, *sessionVal)
	}
	return nil
}

// GetUpByDuration 根据duration 更新session
// 距离session 的过期时间少了duration那么长时间后，就延长 duration
// 比如：设置了15天的有效期，duration设置成1天，当有效期剩余不到 15-1 的时候延长duration
func (this *Instance) GetUpByDuration(w http.ResponseWriter, r *http.Request, name string, dst any, duration time.Duration) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("dst must be a pointer")
	}
	sessionVal, _, err := this.get(w, r, name)
	if err != nil {
		return err
	}
	v.Elem().Set(reflect.ValueOf(sessionVal.Values.Val))
	// 判断距离过期时间少了duration的时候，就延长duration
	if time.Time(sessionVal.ExpireTime).Sub(time.Now()) <= (sessionVal.Duration - duration) {
		sessionVal.ExpireTime = dataType.CustomTime(time.Now().Add(sessionVal.Duration))
		return this.set(w, r, *sessionVal)
	}
	return nil
}

// Flashes 获取并删除session
func (this *Instance) Flashes(w http.ResponseWriter, r *http.Request, name string, dst any) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("dst must be a pointer")
	}
	sessionVal, _, err := this.get(w, r, name)
	if err != nil {
		return err
	}
	this.del(sessionVal.Id, name)
	v.Elem().Set(reflect.ValueOf(sessionVal.Values.Val))
	return nil
}

// Set 设置session
// w
// r
// name  session 名称
// value session 值
// duration session 过期时间，默认为24小时
func (this *Instance) Set(w http.ResponseWriter, r *http.Request, name string, value any, duration ...time.Duration) error {
	sessionId, _ := session2.GetSessionId(w, r, this.option)
	sessionVal := session2.Session{
		Id:         sessionId,
		Name:       name,
		Values:     session2.SessionValue{Val: value},
		CreateTime: dataType.CustomTime{},
		ExpireTime: dataType.CustomTime{},
		Duration:   session2.ExpireTime,
	}
	if len(duration) > 0 {
		sessionVal.Duration = duration[0]
	}
	now := time.Now()
	sessionVal.Id = sessionId                                                 // 设置sessionId
	sessionVal.CreateTime = dataType.CustomTime(now)                          // 设置创建时间
	sessionVal.ExpireTime = dataType.CustomTime(now.Add(sessionVal.Duration)) // 设置过期时间
	return this.set(w, r, sessionVal)
}

// Del 删除session
func (this *Instance) Del(w http.ResponseWriter, r *http.Request, name string) error {
	sessionId, _ := session2.GetSessionId(w, r, this.option)
	_k := session2.GetSessionName(sessionId, name)
	this.del(sessionId, _k)
	return nil
}

// Destroy 销毁session
func (this *Instance) Destroy(w http.ResponseWriter, r *http.Request) error {
	sessionId, err := session2.GetSessionId(w, r, this.option)
	if err != nil {
		return err // 从cookie中获取sessionId失败
	}
	// 需要删除 cookie 或者 header
	session2.DeleteSessionId(w, this.option)
	res := this.rdb.HGetAll(context.Background(), sessionId)
	m, err := res.Result()
	if err != nil {
		return err
	}
	for k, _ := range m {
		this.rdb.Del(context.Background(), k)
	}
	this.rdb.Del(context.Background(), sessionId)
	// 删除当前 sessionid的所有key
	return nil
}
