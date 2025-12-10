package session

import (
	"net/http"
	"reflect"
	"regexp"
	"time"

	"github.com/helays/utils/v2/dataType"
	"github.com/helays/utils/v2/net/http/cookiekit"
	"github.com/helays/utils/v2/tools"
)

// c0449773432e4a478d157a8a923199ac
// 用于校验session id 值是否合规
var sessionRegexp = regexp.MustCompile("^[0-9a-f]{32}$")

const prefix = "vsclub_"

func (m *Manager) cookieName() string {
	return prefix + m.options.Cookie.Name
}

// 获取sessionId
func (m *Manager) getSessionId(w http.ResponseWriter, r *http.Request, sessionId ...string) (string, error) {
	var sid string
	switch m.options.Carrier {
	case CookieCarrierCookie, "":
		cookie, err := r.Cookie(m.options.Cookie.Name)
		if err != nil || !sessionRegexp.MatchString(cookie.Value) {
			if len(sessionId) > 0 && sessionId[0] != "" {
				sid = sessionId[0]
			} else {
				sid = GenerateSessionID()
			}
			m.setSessionId(w, sid)
		} else {
			sid = cookie.Value
		}

	case CookieCarrierHeader:
		sid = r.Header.Get(m.cookieName())
		if sid == "" || !sessionRegexp.MatchString(sid) {
			if len(sessionId) > 0 && sessionId[0] != "" {
				sid = sessionId[0]
			} else {
				sid = GenerateSessionID()
			}
			m.setSessionId(w, sid)
		}
		return sid, nil
	default:
		return "", ErrUnSupport
	}
	return sid, nil
}

func (m *Manager) setSessionId(w http.ResponseWriter, sid string) {
	expire := time.Time{}
	if m.options.Cookie.MaxAge > 0 {
		expire = time.Now().Add(time.Duration(m.options.Cookie.MaxAge) * time.Second)
	}
	switch m.options.Carrier {
	case CookieCarrierCookie, "":
		cfg := m.options.Cookie.Clone()
		cfg.Value = sid
		cfg.Expires = expire
		cookiekit.SetCookie(w, cfg)
	case CookieCarrierHeader:
		w.Header().Set(m.cookieName(), sid)

	}
}

func (m *Manager) deleteSessionId(w http.ResponseWriter) {
	switch m.options.Carrier {
	case CookieCarrierCookie, "":
		cfg := m.options.Cookie.Clone()
		cookiekit.DelCookie(w, cfg)
	case CookieCarrierHeader:
		w.Header().Del(m.cookieName())
	}
}

// GetSessionId 获取sessionId
func (m *Manager) GetSessionId(w http.ResponseWriter, r *http.Request) (string, error) {
	return m.getSessionId(w, r)
}

func (m *Manager) getSession(sessionId string, name string) (*Session, error) {
	sv, err := m.storage.Get(sessionId, name)
	if err != nil {
		return nil, err
	}
	if time.Time(sv.ExpireTime).Before(time.Now()) {
		_ = m.storage.Delete(sessionId, name)
		return nil, ErrNotFound
	}
	return sv, nil
}

// UpdateSessionId 更新sessionId
func (m *Manager) UpdateSessionId(w http.ResponseWriter, r *http.Request, sessionId ...string) (string, error) {
	if err := m.Destroy(w, r); err != nil {
		return "", err
	}
	var sid string
	if len(sessionId) > 0 && sessionId[0] != "" {
		sid = sessionId[0]
	} else {
		sid = GenerateSessionID()
	}
	m.setSessionId(w, sid)
	return sid, nil
}

// Get 获取session
func (m *Manager) Get(w http.ResponseWriter, r *http.Request, name string, dst any) error {
	sessionId, err := m.getSessionId(w, r)
	if err != nil {
		return err
	}
	return m.GetWithSessionID(sessionId, name, dst)
}

func (m *Manager) GetWithSessionID(sessionID string, name string, dst any) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return ErrNotPointer
	}
	sv, err := m.getSession(sessionID, name)
	if err != nil {
		return err
	}
	v.Elem().Set(reflect.ValueOf(sv.Values.Val))
	return nil
}

// GetUp 获取session并更新有效期
func (m *Manager) GetUp(w http.ResponseWriter, r *http.Request, name string, dst any) error {
	sessionId, err := m.getSessionId(w, r)
	if err != nil {
		return err
	}
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return ErrNotPointer
	}
	sv, err := m.getSession(sessionId, name)
	if err != nil {
		return err
	}
	sv.ExpireTime = dataType.CustomTime(time.Now().Add(sv.Duration))
	if err = m.storage.Save(sv); err != nil {
		return err
	}
	v.Elem().Set(reflect.ValueOf(sv.Values.Val))
	return nil
}

// GetUpByRemainTime 根据剩余时间更新session
// 当session 的有效期小于duration，那么将session的有效期延长到 session.Duration-duration
// 比如：设置了15天有效期，duration设置一天，那么当检测到session的有效期 不大于一天的时候就更新session
func (m *Manager) GetUpByRemainTime(w http.ResponseWriter, r *http.Request, name string, dst any, duration time.Duration, callbacks ...Callback) error {
	sessionId, err := m.getSessionId(w, r)
	if err != nil {
		return err
	}
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return ErrNotPointer
	}
	sv, err := m.getSession(sessionId, name)
	if err != nil {
		return err
	}
	v.Elem().Set(reflect.ValueOf(sv.Values.Val))

	// 注意：当剩余时间 <= duration 时，会续期整个session.Duration
	// 例如：session=1小时，duration=1分钟
	// 剩余59秒时会触发，续期到1小时后
	if sv.ExpireTime.Sub(time.Now()) <= duration {
		return m.extendSession(sv, callbacks...)
	}
	return nil
}

// GetUpByDuration 根据duration
// 距离session 的过期时间少了duration那么长时间后，就延长 duration
// 比如：设置了15天的有效期，duration设置成1天，当有效期剩余不到 15-1 的时候延长duration
func (m *Manager) GetUpByDuration(w http.ResponseWriter, r *http.Request, name string, dst any, duration time.Duration, callbacks ...Callback) error {
	sessionId, err := m.getSessionId(w, r)
	if err != nil {
		return err
	}
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return ErrNotPointer
	}
	sv, err := m.getSession(sessionId, name)
	if err != nil {
		return err
	}
	v.Elem().Set(reflect.ValueOf(sv.Values.Val))

	// 注意，当有效期剩余时间少于 duration-sv.Duration 时，会触发
	// 例如：session=15天，duration=1天
	// 那么时间过了一天，就会触发，延长session有效期到15天
	if sv.ExpireTime.Sub(time.Now()) <= (sv.Duration - duration) {
		return m.extendSession(sv, callbacks...)
	}
	return nil
}

// 延长session有效期
func (m *Manager) extendSession(sv *Session, callbacks ...Callback) error {
	sv.ExpireTime = dataType.CustomTime(time.Now().Add(sv.Duration))
	for _, cb := range callbacks {
		if cb.BeforeRenew != nil {
			if err := cb.BeforeRenew(sv.ExpireTime, sv.Values.Val); err != nil {
				return err
			}
		}
	}
	if err := m.storage.Save(sv); err != nil {
		return err
	}
	for _, cb := range callbacks {
		if cb.AfterRenew != nil {
			if err := cb.AfterRenew(sv.ExpireTime, sv.Values.Val); err != nil {
				return err
			}
		}
	}
	return nil
}

// Flashes 获取并删除session
func (m *Manager) Flashes(w http.ResponseWriter, r *http.Request, name string, dst any) error {
	sessionId, err := m.getSessionId(w, r)
	if err != nil {
		return err
	}
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return ErrNotPointer
	}
	sv, err := m.getSession(sessionId, name)
	if err != nil {
		return err
	}
	v.Elem().Set(reflect.ValueOf(sv.Values.Val))
	return m.storage.Delete(sv.Id, sv.Name)
}

// Set 设置session
func (m *Manager) Set(w http.ResponseWriter, r *http.Request, value *Value) error {
	sessionId, err := m.getSessionId(w, r)
	if err != nil {
		return err
	}
	value.SessionID = sessionId
	return m.SetVal(value)
}

func (m *Manager) SetVal(value *Value) error {
	now := time.Now()
	sv := Session{
		Id:         value.SessionID,
		Name:       value.Field,
		Values:     NewSessionValue(value.Value),
		CreateTime: dataType.NewCustomTime(now),
		Duration:   tools.Ternary(value.TTL > 0, value.TTL, ExpireTime),
	}
	sv.ExpireTime = dataType.CustomTime(now.Add(sv.Duration)) // 设置过期时间
	return m.storage.Save(&sv)
}

func (m *Manager) SetWithNewSessionId(w http.ResponseWriter, r *http.Request, value *Value) error {
	sid, err := m.UpdateSessionId(w, r, value.SessionID)
	if err != nil {
		return err
	}
	value.SessionID = sid
	return m.SetVal(value)
}

func (m *Manager) Del(w http.ResponseWriter, r *http.Request, name string) error {
	sessionId, err := m.getSessionId(w, r)
	if err != nil {
		return err
	}
	return m.storage.Delete(sessionId, name)
}

func (m *Manager) Destroy(w http.ResponseWriter, r *http.Request) error {
	sessionId, err := m.getSessionId(w, r)
	if err != nil {
		return err
	}
	m.deleteSessionId(w)
	return m.storage.DeleteAll(sessionId)
}

func (m *Manager) Close() error {
	return m.storage.Close()
}
