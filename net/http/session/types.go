package session

import (
	"errors"
	"time"

	"github.com/helays/utils/v2/dataType"
)

//  这个需要移除 上级Session 已经实现了二进制序列化
//func (s SessionValue) GobEncode() ([]byte, error) {
//	return msgpack.Marshal(s.val)
//}
//
//func (s *SessionValue) GobDecode(data []byte) error {
//	return msgpack.Unmarshal(data, &s.val)
//}

type CookieCarrier string

func (c CookieCarrier) String() string {
	return string(c)
}

const (
	CookieCarrierCookie = "cookie"
	CookieCarrierHeader = "header"
)

var (
	ErrUnSupport  = errors.New("不支持的 session 载体")
	ErrNotFound   = errors.New("session 不存在")
	ErrNotPointer = errors.New("session 变量目标必须是指针")
)

// noinspection all
const SessionID = "session_id"

type Callback struct {
	BeforeRenew func(sessionID string, expire dataType.CustomTime, data any) error
	AfterRenew  func(sessionID string, expire dataType.CustomTime, data any) error
}

type Value struct {
	SessionID string        // 可自定义session id
	Field     string        // session 值字段
	Value     any           // session 值
	TTL       time.Duration // 有效期
}
