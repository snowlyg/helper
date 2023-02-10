package global

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/snowlyg/helper/arr"
)

func GetMacAddr() string {
	if getMacAddrInterface() == nil {
		return ""
	}
	addr := getMacAddrInterface().HardwareAddr.String()
	if len(addr) > 0 {
		addr = strings.ReplaceAll(addr, ":", "")
		addr = strings.ToUpper(addr)
		return addr
	}
	return ""
}

func getMacAddrInterface() *net.Interface {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	re, err := regexp.Compile(`^(ens|eth|waln|以太网|Ethernet)[0-9]*`)
	if err != nil {
		return nil
	}
	nameCheck := arr.NewCheckArrayType(0)
	nameCheck.AddMutil("eth0", "waln0")
	for _, netInterface := range netInterfaces {
		if nameCheck.Check(netInterface.Name) {
			return &netInterface
		}
		if re.MatchString(netInterface.Name) {
			return &netInterface
		}
	}
	return nil
}

func check(ip, network string) bool {
	parseIp, subnet, err := net.ParseCIDR(network)
	if err != nil {
		return false
	}
	if ip == parseIp.String() {
		return false
	}
	return subnet.Contains(net.ParseIP(ip))
}

func LocalIP(network string) string {
	ip := ""
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for i, addr := range addrs {
			i += 1
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !ipnet.IP.IsMulticast() && !ipnet.IP.IsLinkLocalUnicast() && !ipnet.IP.IsLinkLocalMulticast() && ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				if len(ip) > 0 && check(ip, network) {
					return ip
				}
			}
		}
	}
	return ip
}

func IsPortInUse(host string, port int64) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, fmt.Sprintf("%d", port)), time.Second*1)
	if err == nil {
		conn.Close()
		return true
	}

	return false
}
