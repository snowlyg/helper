package ping

import (
	"testing"
)

func Test_GetPingMsg(t *testing.T) {
	ips := []string{"www.github.com", "www.google.com", "www.baidu.com", "www.chindeo.com", "10.0.0.113", "10.0.0.121"}
	for _, ip := range ips {
		t.Run("测试 ping 方法:"+ip, func(t *testing.T) {
			ok, msg := GetPingMsg(ip)
			if !ok {
				t.Errorf("%s ping is fault,get msg %s", ip, msg)
			}
		})
	}
}
