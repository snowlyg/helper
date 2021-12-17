package ping

import (
	"fmt"
	"time"

	"github.com/go-ping/ping"
)

// GetPingMsg ping 检查网络情况
// icmp 检查5次，每次300毫秒超时，一共1500毫秒超时
// 只要有一个次响应就成功
func GetPingMsg(devIp string) (bool, string) {
	if devIp == "" {
		return false, "设备ip为空，请检查设备是否绑定ip"
	}
	pinger, err := ping.NewPinger(devIp)
	if err != nil {
		return false, err.Error()
	}
	pinger.Count = 3
	pinger.Interval = time.Duration(500 * time.Millisecond)
	pinger.Timeout = time.Duration(1500 * time.Millisecond)
	pinger.SetPrivileged(true)
	err = pinger.Run()
	if err != nil {
		return false, err.Error()
	}
	stats := pinger.Statistics()
	if stats.PacketsRecv >= 1 {
		return true, fmt.Sprintf("设备ip(%s)可以访问", devIp)
	}
	return false, fmt.Sprintf("设备(%s)可能已离线或者网络不稳定", devIp)
}
