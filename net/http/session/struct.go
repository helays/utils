package session

import (
	"github.com/helays/utils/dataType"
	"io"
	"net/http"
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
// Date: 2024/12/7 22:16
//

// Dialector session 接口
type Dialector interface {
	Get(w http.ResponseWriter, r *http.Request, name string, dst any) error                                     // 获取session
	GetUp(w http.ResponseWriter, r *http.Request, name string, dst any) error                                   // 获取并刷新
	GetUpByTimeLeft(w http.ResponseWriter, r *http.Request, name string, dst any, duration time.Duration) error // 比如：设置了15天有效期，duration设置一天，那么当检测到session的有效期 不大于一天的时候就更新session
	GetUpByDuration(w http.ResponseWriter, r *http.Request, name string, dst any, duration time.Duration) error // 比如：设置了15天的有效期，duration设置成1天，当有效期剩余不到 15-1 的时候延长duration
	Flashes(w http.ResponseWriter, r *http.Request, name string, dst any) error                                 // 获取session然后销毁
	Set(w http.ResponseWriter, r *http.Request, name string, value any, duration ...time.Duration) error        // 设置session
	Del(w http.ResponseWriter, r *http.Request, name string) error                                              // 删除session
	Destroy(w http.ResponseWriter, r *http.Request) error                                                       // 销毁session

	Apply(*Options) // 设置
	io.Closer
}

// Session session 数据结构
type Session struct {
	Id         string              `json:"id" gorm:"primaryKey;autoIncrement:false;type:varchar(64);not null;index;comment:Session ID"`    // session id
	Name       string              `json:"name" gorm:"primaryKey;autoIncrement:false;type:varchar(128);not null;index;comment:Session的名字"` // session 的名字
	Values     SessionValue        `json:"values" gorm:"comment:session数据"`                                                                // session 数据
	CreateTime dataType.CustomTime `json:"create_time" gorm:"comment:session 创建时间"`                                                        // 创建时间 ，没啥用，就看
	ExpireTime dataType.CustomTime `json:"expire_time" gorm:"not null;index;comment:session 过期时间"`                                         // 过期时间 ，用于自动回收的时候使用
	Duration   time.Duration       `json:"duration" gorm:"comment:session有效期"`                                                             // 有效期，主要是用于更新有效期的时候使用
}

func (this Session) TableName() string {
	return "vslclub_session"
}

// Options session 配置
type Options struct {
	CookieName    string        `json:"cookie_name" yaml:"cookie_name" ini:"cookie_name"`          // 从cookie或者 header中读取 session的标识
	DisableGc     bool          `json:"disable_gc" yaml:"disable_gc" ini:"disable_gc"`             // 是否禁用session gc
	GcProbability float64       `json:"gc_probability" yaml:"gc_probability" ini:"gc_probability"` // session gc 概率
	CheckInterval time.Duration `json:"check_interval" yaml:"check_interval" ini:"check_interval"` // session 检测默认有效期
	Carrier       string        `json:"carrier" yaml:"carrier" ini:"carrier"`                      // session 载体，默认cookie
	// cookie相关配置
	Path   string `json:"path" yaml:"path" ini:"path"`
	Domain string `json:"domain" yaml:"domain" ini:"domain"`
	// MaxAge=0 means no Max-Age attribute specified and the cookie will be
	// deleted after the browser session ends.
	// MaxAge<0 means delete cookie immediately.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge   int           `json:"max_age" yaml:"max_age" ini:"max_age"`
	Secure   bool          `json:"secure" yaml:"secure" ini:"secure"`
	HttpOnly bool          `json:"http_only" yaml:"http_only" ini:"http_only"`
	SameSite http.SameSite `json:"same_site" yaml:"same_site" ini:"same_site"`
}
