package linux

import (
	"github.com/pkg/errors"
	"net"
)

// GetDefaultIP get a default non local ip, err is not nil, ip return 127.0.0.1
func GetDefaultIP() (ip string, err error) {
	ip = "127.0.0.1"

	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipStr := getAddrDefaultIP(addr); len(ipStr) > 0 {
				return ipStr, nil
			}
		}
	}

	return "", errors.New("no ip found.")
}

func getAddrDefaultIP(addr net.Addr) string {
	var ip net.IP

	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	default:
		return ""
	}
	if ip.IsUnspecified() || ip.IsLoopback() {
		return ""
	}

	ip = ip.To4()
	if ip == nil {
		return ""
	}

	return ip.String()
}
