package ipAccess

import (
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/helays/utils/v2/net/ipkit"
)

type IPList struct {
	singleIPs map[string]struct{}
	subnetsV4 []*net.IPNet
	subnetsV6 []*net.IPNet
	rangesV4  []ipRange
	rangesV6  []ipRange
}

// ipRange 表示一个IP段
type ipRange struct {
	start net.IP
	end   net.IP
}

// NewIPList 创建一个新的IPList，支持可选参数（可变参数）
func NewIPList(ipList ...string) (*IPList, error) {
	list := &IPList{
		singleIPs: make(map[string]struct{}),
		subnetsV4: make([]*net.IPNet, 0),
		subnetsV6: make([]*net.IPNet, 0),
		rangesV4:  make([]ipRange, 0),
		rangesV6:  make([]ipRange, 0),
	}

	// 如果传入了IP列表，则添加到IPList中
	for _, item := range ipList {
		if err := list.Add(item); err != nil {
			return nil, fmt.Errorf("无法添加IP %s: %v", item, err)
		}
	}
	return list, nil
}

// Add 添加一个IP、子网或IP段到IPList中
func (l *IPList) Add(item string) error {
	// 如果是IP段模式（例如 "192.168.1.1-192.168.1.100"）
	if strings.Contains(item, "-") {
		return l.addIPRange(item)
	}

	// 如果是子网模式（例如 "192.168.1.0/24"）
	if strings.Contains(item, "/") {
		return l.addSubnet(item)
	}
	return l.addSingleIP(item)
}
func (l *IPList) addSingleIP(ip string) error {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return fmt.Errorf("无效的IP格式: %s", ip)
	}
	l.singleIPs[ip] = struct{}{}
	return nil
}

func (l *IPList) addSubnet(cidr string) error {
	_, subnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("无效的子网格式: %s", cidr)
	}

	if subnet.IP.To4() != nil {
		l.subnetsV4 = append(l.subnetsV4, subnet)
	} else {
		l.subnetsV6 = append(l.subnetsV6, subnet)
	}
	return nil
}
func (l *IPList) addIPRange(item string) error {
	rangeParts := strings.Split(item, "-")
	if len(rangeParts) != 2 {
		return fmt.Errorf("无效的IP段格式: %s", item)
	}

	startIP := net.ParseIP(rangeParts[0])
	endIP := net.ParseIP(rangeParts[1])
	if startIP == nil || endIP == nil {
		return fmt.Errorf("IP段中包含无效的IP地址: %s", item)
	}

	if isIPv4(startIP) != isIPv4(endIP) {
		return fmt.Errorf("IP段中的起始IP和结束IP版本不一致: %s", item)
	}

	if compareIPs(startIP, endIP) > 0 {
		return fmt.Errorf("起始IP必须小于或等于结束IP: %s", item)
	}

	if isIPv4(startIP) {
		l.insertRangeV4(startIP, endIP)
	} else {
		l.insertRangeV6(startIP, endIP)
	}
	return nil
}

func (l *IPList) insertRangeV4(start, end net.IP) {
	startInt := ipkit.Ip2Int(start)
	idx := sort.Search(len(l.rangesV4), func(i int) bool {
		return ipkit.Ip2Int(l.rangesV4[i].start) >= startInt
	})
	l.rangesV4 = append(l.rangesV4, ipRange{})
	copy(l.rangesV4[idx+1:], l.rangesV4[idx:])
	l.rangesV4[idx] = ipRange{start: start, end: end}
}

func (l *IPList) insertRangeV6(start, end net.IP) {
	idx := sort.Search(len(l.rangesV6), func(i int) bool {
		return compareIPs(l.rangesV6[i].start, start) > 0
	})

	l.rangesV6 = append(l.rangesV6, ipRange{})
	copy(l.rangesV6[idx+1:], l.rangesV6[idx:])
	l.rangesV6[idx] = ipRange{start: start, end: end}
}

// Contains 检查IP是否在IPList中
func (l *IPList) Contains(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// 检查单个IP
	if _, exists := l.singleIPs[ip]; exists {
		return true
	}

	if isIPv4(parsedIP) {
		return l.containsIPv4(parsedIP)
	}
	return l.containsIPv6(parsedIP)
}

func (l *IPList) containsIPv4(ip net.IP) bool {
	for _, subnet := range l.subnetsV4 {
		if subnet.Contains(ip) {
			return true
		}
	}

	target := ipkit.Ip2Int(ip)
	idx := sort.Search(len(l.rangesV4), func(i int) bool {
		return ipkit.Ip2Int(l.rangesV4[i].start) > target
	})
	if idx > 0 {
		r := l.rangesV4[idx-1]
		if target <= ipkit.Ip2Int(r.end) {
			return true
		}
	}
	return false
}

func (l *IPList) containsIPv6(ip net.IP) bool {
	ip = ip.To16()
	for _, subnet := range l.subnetsV6 {
		if subnet.Contains(ip) {
			return true
		}
	}
	idx := sort.Search(len(l.rangesV6), func(i int) bool {
		return compareIPs(l.rangesV6[i].start, ip) > 0
	})
	if idx > 0 {
		r := l.rangesV6[idx-1]
		if compareIPs(ip, r.end) <= 0 {
			return true
		}
	}
	return false
}

// isIPv4 判断IP是否为IPv4地址
func isIPv4(ip net.IP) bool {
	return ip.To4() != nil
}

// compareIPs 比较两个IP地址的大小
func compareIPs(ip1, ip2 net.IP) int {
	ip1 = ip1.To16()
	ip2 = ip2.To16()

	for i := 0; i < len(ip1); i++ {
		if ip1[i] < ip2[i] {
			return -1
		} else if ip1[i] > ip2[i] {
			return 1
		}
	}
	return 0
}
