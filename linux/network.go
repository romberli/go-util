package linux

import (
	"github.com/pkg/errors"
	"net"
)

// GetDefaultIP get a default non local ip, err is not nil, ip return 127.0.0.1
func GetDefaultIP() (ip string, err error) {
	ip = "127.0.0.1"

	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, i := range ifaces {
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

	err = errors.New("no ip found")
	return
}

func getAddrDefaultIP(addr net.Addr) string {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
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

// IsValidHost judge the host is a valid host or not.
func IsValidHost(host string) bool {
	if len(host) == 0 {
		return false
	}

	if ip := net.ParseIP(host); ip != nil {
		if ip.IsLoopback() {
			return false
		}
	}

	return host != "localhost"
}
