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

func GetListeningAddr(listenAddr string) (string, error) {
	// 解析监听地址（如 ":10001"）
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return "", fmt.Errorf("解析地址失败: %w", err)
	}

	// 如果配置中已明确指定IP，直接返回
	if addr.IP != nil && !addr.IP.IsUnspecified() {
		return fmt.Sprintf("%s:%d", addr.IP, addr.Port), nil
	}

	// 单次遍历网络接口
	var firstIP net.IP
	ifaces, _ := net.Interfaces() // 错误可忽略
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, a := range addrs {
			ipnet, ok := a.(*net.IPNet)
			if !ok {
				continue
			}
			ip := ipnet.IP
			if ip.IsLoopback() {
				// 发现环回地址立即返回
				return fmt.Sprintf("127.0.0.1:%d", addr.Port), nil
			}

			// 记录第一个非环回IPv4地址
			if firstIP == nil && ip.To4() != nil {
				firstIP = ip
			}
		}
	}

	// 没有环回地址时返回第一个找到的IP
	if firstIP != nil {
		return fmt.Sprintf("%s:%d", firstIP, addr.Port), nil
	}

	// 保底返回环回地址（即使没有实际接口）
	return fmt.Sprintf("127.0.0.1:%d", addr.Port), nil
}
