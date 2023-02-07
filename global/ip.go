package global

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/snowlyg/helper/arr"
)

func GetMacAddr() string {
	macAddr := ""
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return macAddr
	}

	for _, netInterface := range netInterfaces {
		flags := strings.Split(netInterface.Flags.String(), "|")
		flagsCheck := arr.NewCheckArrayType(len(flags))
		for _, flag := range flags {
			flagsCheck.Add(flag)
		}
		if !flagsCheck.Check(net.FlagUp.String()) {
			continue
		}
		addr := netInterface.HardwareAddr.String()
		if len(addr) > 0 {
			addr = strings.ReplaceAll(addr, ":", "")
			addr = strings.ToUpper(addr)
			return addr
		}

	}
	return macAddr
}

func LocalIP() string {
	ip := ""
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !ipnet.IP.IsMulticast() && !ipnet.IP.IsLinkLocalUnicast() && !ipnet.IP.IsLinkLocalMulticast() && ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				if len(ip) > 0 {
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
