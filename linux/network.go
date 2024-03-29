package linux

import (
	"net"
	"strconv"
	"strings"

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

// IsValidIP checks if the ip is valid
func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// CompareIP compares two ip addresses,
// if ip1 is equal to ip2, it returns 0
// if ip1 is less than ip2, it returns -1
// if ip1 is greater than ip2, it returns 1
func CompareIP(ip1, ip2 string) (int, error) {
	if !IsValidIP(ip1) {
		return constant.ZeroInt, errors.Errorf("ip addr must be formatted as A.B.C.D, %s is not valid", ip1)
	}
	if !IsValidIP(ip2) {
		return constant.ZeroInt, errors.Errorf("ip addr must be formatted as A.B.C.D, %s is not valid", ip2)
	}

	ipList1 := strings.Split(ip1, constant.DotString)
	ipList2 := strings.Split(ip2, constant.DotString)

	for i := constant.ZeroInt; i < len(ipList1); i++ {
		p1, err := strconv.Atoi(ipList1[i])
		if err != nil {
			return constant.ZeroInt, errors.Trace(err)
		}
		p2, err := strconv.Atoi(ipList2[i])
		if err != nil {
			return constant.ZeroInt, errors.Trace(err)
		}

		if p1 < p2 {
			return -1, nil
		}
		if p1 > p2 {
			return 1, nil
		}
	}

	return constant.ZeroInt, nil
}

// CompareAddr compares two addresses,
// if addr1 is equal to addr2, it returns 0
// if addr1 is less than addr2, it returns -1
// if addr1 is greater than addr2, it returns 1
func CompareAddr(addr1, addr2 string) (int, error) {
	ip1, port1, err := net.SplitHostPort(addr1)
	if err != nil {
		return constant.ZeroInt, errors.Trace(err)
	}
	ip2, port2, err := net.SplitHostPort(addr2)
	if err != nil {
		return constant.ZeroInt, errors.Trace(err)
	}

	port1Int, err := strconv.Atoi(port1)
	if err != nil {
		return constant.ZeroInt, errors.Trace(err)
	}
	port2Int, err := strconv.Atoi(port2)
	if err != nil {
		return constant.ZeroInt, errors.Trace(err)
	}

	result, err := CompareIP(ip1, ip2)
	if err != nil {
		return constant.ZeroInt, err
	}
	if result == constant.ZeroInt {
		if port1Int < port2Int {
			return -1, nil
		}
		if port1Int > port2Int {
			return 1, nil
		}
	}

	return result, nil
}

// GetMinAddr gets the minimum address from the addr list
func GetMinAddr(addrs []string) (string, error) {
	if len(addrs) == constant.ZeroInt {
		return constant.EmptyString, errors.New("addr list is empty")
	}

	minAddr := addrs[constant.ZeroInt]
	for _, addr := range addrs {
		result, err := CompareAddr(addr, minAddr)
		if err != nil {
			return constant.EmptyString, err
		}
		if result == -1 {
			minAddr = addr
		}
	}

	return minAddr, nil
}

// SortAddrs sorts the addr list with quick sort algorithm
func SortAddrs(addrs []string) error {
	if len(addrs) == constant.ZeroInt {
		return nil
	}

	return quickSortAddrs(addrs, constant.ZeroInt, len(addrs)-1)
}

// quickSortAddrs sorts the addr list with quick sort algorithm
func quickSortAddrs(addrs []string, left, right int) error {
	if left >= right {
		return nil
	}

	var err error
	pivot := left
	index := pivot + 1
	i := index
	for i <= right {
		result, err := CompareAddr(addrs[i], addrs[pivot])
		if err != nil {
			return err
		}
		if result == -1 {
			addrs[i], addrs[index] = addrs[index], addrs[i]
			index++
		}
		i++
	}
	addrs[pivot], addrs[index-1] = addrs[index-1], addrs[pivot]
	pivot = index - 1

	err = quickSortAddrs(addrs, left, pivot-1)
	if err != nil {
		return err
	}
	err = quickSortAddrs(addrs, pivot+1, right)
	if err != nil {
		return err
	}

	return nil
}
