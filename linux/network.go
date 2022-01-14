package linux

import (
	"net"

	"github.com/pingcap/errors"
	"github.com/romberli/go-util/constant"
)

// GetDefaultIP gets the default non-local ip, if there are more than one ips, it will return the first one
func GetDefaultIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return constant.EmptyString, errors.Trace(err)
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipStr := getAddrDefaultIP(addr)
			if len(ipStr) > constant.ZeroInt {
				return ipStr, nil
			}
		}
	}

	return constant.EmptyString, errors.New("no ip found")
}

// getAddrDefaultIP returns default IP of host
func getAddrDefaultIP(addr net.Addr) string {
	var ip net.IP

	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	default:
		return constant.EmptyString
	}
	if ip.IsUnspecified() || ip.IsLoopback() {
		return constant.EmptyString
	}

	ip = ip.To4()
	if ip == nil {
		return constant.EmptyString
	}

	return ip.String()
}
