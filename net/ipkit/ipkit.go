package ipkit

import (
	"encoding/binary"
	"net"
	"net/netip"
)

// Ip2Int 将IP地址转换为整数
func Ip2Int(ip net.IP) uint32 {
	return binary.BigEndian.Uint32(ip.To4())
}

// Int2IP converts a 32-bit integer back to an IPv4 address.
func Int2IP(ipInt uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, ipInt)
	return ip
}

func IPAddr2Int(addr netip.Addr) uint32 {
	if !addr.Is4() {
		return 0
	}
	bytes := addr.As4()
	return uint32(bytes[0])<<24 | uint32(bytes[1])<<16 | uint32(bytes[2])<<8 | uint32(bytes[3])
}
