package carrierFile

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/helays/utils/close/vclose"
	"github.com/helays/utils/dataType"
	session2 "github.com/helays/utils/net/http/session"
	"github.com/helays/utils/tools"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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
// Date: 2024/12/8 1:50
//

//var (
//	db *badger.DB
//)

// Instance session 实例
type Instance struct {
	option *session2.Options
	Path   string `json:"path" yaml:"path" ini:"path"` // db路径
	ctx    context.Context
	cancel context.CancelFunc
}

// New 初始化 session 内存 实例
func New(opt ...Instance) (*Instance, error) {
	ins := &Instance{
		Path: "runtime/session",
	}
	if len(opt) > 0 {
		ins.Path = opt[0].Path
	}
	ins.Path = tools.Fileabs(ins.Path)
	err := tools.Mkdir(ins.Path)
	if err != nil {
		return nil, fmt.Errorf("创建session文件存放目录失败")
	}

	return ins, nil
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

// Apply 应用配置
func (this *Instance) Apply(options *session2.Options) {
	this.option = options
	this.ctx, this.cancel = context.WithCancel(context.Background())
	tools.RunAsyncTickerProbabilityFunc(this.ctx, !this.option.DisableGc, this.option.CheckInterval, this.option.GcProbability, this.gc)
}

// Close 关闭 db
func (this *Instance) Close() error {
	this.cancel()
	return nil
}

// gc 垃圾回收
func (this *Instance) gc() {
	files, err := os.ReadDir(this.Path)
	if err != nil {
		return
	}
	// 循环所有文件
	// 如果是文件夹，就直接删除
	// 如果文件打开失败，跳过处理
	// 如果解析失败，就删除
	// 判断是否过期，过期也直接删除
	for _, file := range files {
		// 读取所有session文件
		sessionPath := filepath.Join(this.Path, file.Name())
		if file.IsDir() {
			this.del(sessionPath)
			continue
		}
		sessionVal := &session2.Session{}

		f, err := os.Open(sessionPath)
		if err != nil {
			vclose.Close(f)
			continue
		}
		if err = gob.NewDecoder(f).Decode(sessionVal); err != nil {
			vclose.Close(f)
			this.del(sessionPath)
			continue
		}
		vclose.Close(f)
		if time.Time(sessionVal.ExpireTime).Before(time.Now()) {
			this.del(sessionPath)
		}
	}
}

// 从session 文件中读取session 数据
func (this *Instance) get(w http.ResponseWriter, r *http.Request, name string) (*session2.Session, string, error) {
	sessionId, err := session2.GetSessionId(w, r, this.option) // 这一步一般不会失败
	if err != nil {
		return nil, "", err // 从cookie中获取sessionId失败
	}
	_k := session2.GetSessionName(sessionId, name)
	// 从文件中读取数据
	sessionPath := filepath.Join(this.Path, _k)
	f, err := os.Open(sessionPath)
	defer vclose.Close(f)
	if err != nil {
		return nil, "", err
	}
	sessionVal := &session2.Session{}
	if err = gob.NewDecoder(f).Decode(sessionVal); err != nil {
		vclose.Close(f)
		// session 数据解析失败，删除session文件
		this.del(sessionPath)
		return nil, "", err
	}
	vclose.Close(f)
	if time.Time(sessionVal.ExpireTime).Before(time.Now()) {
		this.del(sessionPath)
		// session已过期
		return nil, "", session2.ErrNotFound
	}
	return sessionVal, _k, nil
}

// 设置session，将其通过gob 写入文件
func (this *Instance) set(w http.ResponseWriter, r *http.Request, dst session2.Session) error {
	_k := session2.GetSessionName(dst.Id, dst.Name)
	sessionPath := filepath.Join(this.Path, _k)
	f, err := os.OpenFile(sessionPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	defer vclose.Close(f)
	if err != nil {
		return err
	}
	return gob.NewEncoder(f).Encode(dst)
}

// 删除具体的session 文件
func (this *Instance) del(path string) {
	_ = os.RemoveAll(path)
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
	sessionVal, _k, err := this.get(w, r, name)
	this.del(filepath.Join(this.Path, _k))
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
func (this *Instance) Set(w http.ResponseWriter, r *http.Request, name string, value any, duration ...time.Duration) error {
	sessionId, _ := session2.GetSessionId(w, r, this.option)
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
	return this.set(w, r, sessionVal)
}

// Del 删除session
func (this *Instance) Del(w http.ResponseWriter, r *http.Request, name string) error {
	sessionId, _ := session2.GetSessionId(w, r, this.option)
	_k := session2.GetSessionName(sessionId, name)
	this.del(filepath.Join(this.Path, _k))
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
	files, err := os.ReadDir(this.Path)
	if err != nil {
		return err
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), sessionId) {
			filePath := filepath.Join(this.Path, file.Name())
			this.del(filePath)
		}
	}
	return nil
}
