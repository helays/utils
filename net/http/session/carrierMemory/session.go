package carrierMemory

import (
	"context"
	"fmt"
	"github.com/helays/utils/v2/dataType"
	session2 "github.com/helays/utils/v2/net/http/session"
	"github.com/helays/utils/v2/tools"
	"net/http"
	"reflect"
	"sync"
	"time"
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
// Date: 2024/12/7 23:23
//
// vsclub:sessionId:key:value

// Instance session 实例
type Instance struct {
	option *session2.Options
	ctx    context.Context
	cancel context.CancelFunc
}

// New 初始化 session 内存 实例
func New() *Instance {
	return &Instance{}
}

func (this *Instance) Apply(options *session2.Options) {
	this.option = options
	this.ctx, this.cancel = context.WithCancel(context.Background())
	// 还需要自动删除
	tools.RunAsyncTickerProbabilityFunc(this.ctx, !this.option.DisableGc, this.option.CheckInterval, this.option.GcProbability, this.gc)
}

func (this *Instance) gc() {
	sessionStorage.Range(func(key, value any) bool {
		ss := value.(session2.Session)
		if time.Time(ss.ExpireTime).Before(time.Now()) {
			sessionStorage.Delete(key)
		}
		return true
	})
}

func (this *Instance) Close() error {
	this.cancel() // 关闭定时器
	return nil
}

var (
	sessionStorage sync.Map // 存储session
)

func (this *Instance) get(w http.ResponseWriter, r *http.Request, name string) (*session2.Session, string, error) {
	sessionId, err := session2.GetSessionId(w, r, this.option) // 这一步一般不会失败
	if err != nil {
		return nil, "", err // 从cookie中获取sessionId失败
	}
	_k := session2.GetSessionName(sessionId, name)
	val, ok := sessionStorage.Load(_k)
	if !ok {
		return nil, "", session2.ErrNotFound
	}
	sessionVal := val.(session2.Session)
	if time.Time(sessionVal.ExpireTime).Before(time.Now()) {
		sessionStorage.Delete(_k)
		// session已过期
		return nil, "", session2.ErrNotFound
	}
	return &sessionVal, _k, nil
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
	sessionVal, _k, err := this.get(w, r, name)
	if err != nil {
		return err
	}
	// 更新session过期时间
	sessionVal.ExpireTime = dataType.CustomTime(time.Now().Add(sessionVal.Duration))
	sessionStorage.Store(_k, *sessionVal)
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
	sessionVal, _k, err := this.get(w, r, name)
	if err != nil {
		return err
	}
	v.Elem().Set(reflect.ValueOf(sessionVal.Values.Val))
	// 判断 距离过期时间小于等于duration 的时候，更新session的过期时间
	if time.Time(sessionVal.ExpireTime).Sub(time.Now()) <= duration {
		sessionVal.ExpireTime = dataType.CustomTime(time.Now().Add(sessionVal.Duration))
		sessionStorage.Store(_k, *sessionVal)
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
	sessionVal, _k, err := this.get(w, r, name)
	if err != nil {
		return err
	}
	v.Elem().Set(reflect.ValueOf(sessionVal.Values.Val))
	// 判断距离过期时间少了duration的时候，就延长duration
	if time.Time(sessionVal.ExpireTime).Sub(time.Now()) <= (sessionVal.Duration - duration) {
		sessionVal.ExpireTime = dataType.CustomTime(time.Now().Add(sessionVal.Duration))
		sessionStorage.Store(_k, *sessionVal)
	}
	return nil
}

// Flashes 获取并删除session
func (this *Instance) Flashes(w http.ResponseWriter, r *http.Request, name string, dst any) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("dst must be a pointer")
	}
	sessionVal, _k, err := this.get(w, r, name)
	sessionStorage.Delete(_k)
	if err != nil {
		return err
	}
	v.Elem().Set(reflect.ValueOf(sessionVal.Values.Val))
	return nil
}

// Set 设置session
// w
// r
// name  session 名称
// value session 值
// duration session 过期时间，默认为24小时
// 如果session 从不过期，可能会导致 session数据逃逸
func (this *Instance) Set(w http.ResponseWriter, r *http.Request, name string, value any, duration ...time.Duration) error {
	sessionId, _ := session2.GetSessionId(w, r, this.option)
	_k := session2.GetSessionName(sessionId, name)
	now := time.Now()
	sessionVal := session2.Session{
		Id:         sessionId,
		Name:       name,
		Values:     session2.SessionValue{Val: value},
		CreateTime: dataType.CustomTime(now),
		ExpireTime: dataType.CustomTime{},
		Duration:   session2.ExpireTime,
	}
	if len(duration) > 0 {
		sessionVal.Duration = duration[0]
	}
	sessionVal.ExpireTime = dataType.CustomTime(now.Add(sessionVal.Duration)) // 设置过期时间
	sessionStorage.Store(_k, sessionVal)
	return nil
}

// Del 删除session
func (this *Instance) Del(w http.ResponseWriter, r *http.Request, name string) error {
	sessionId, _ := session2.GetSessionId(w, r, this.option)
	_k := session2.GetSessionName(sessionId, name)
	sessionStorage.Delete(_k)
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
	// 删除所有以 sessionId 为前缀的 key
	sessionStorage.Range(func(key, value any) bool {
		sessVal := value.(session2.Session)
		if sessVal.Id == sessionId {
			sessionStorage.Delete(key)
		}
		return true
	})
	return nil
}
