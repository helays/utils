package config

import "errors"

var (
	ErrProtocolInvalid = errors.New("协议无效")
	ErrInvalidParam    = errors.New("无效的参数")
)
