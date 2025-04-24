package checkIp

import (
	"fmt"
	"net"
)

// IsIP 判断输入是否是IP
func IsIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

func ParseIPAndPort(addr string) (net.IP, int, error) {
	// 尝试直接解析为IP（不带端口的情况）
	if ip := net.ParseIP(addr); ip != nil {
		return ip, 0, nil // 返回0表示没有端口
	}

	// 处理带端口的情况
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, 0, fmt.Errorf("地址格式错误，应为'ip:port'或'[ipv6]:port'：%v", err)
	}

	// 如果host被方括号包围（IPv6），去掉方括号
	if len(host) >= 2 && host[0] == '[' && host[len(host)-1] == ']' {
		host = host[1 : len(host)-1]
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return nil, 0, fmt.Errorf("无效的IP地址: %s", host)
	}

	portNum, err := parsePort(port)
	if err != nil {
		return nil, 0, err
	}

	return ip, portNum, nil
}

func parsePort(portStr string) (int, error) {
	port, err := net.LookupPort("tcp", portStr)
	if err != nil {
		return 0, fmt.Errorf("无效的端口号: %v", err)
	}
	return port, nil
}
