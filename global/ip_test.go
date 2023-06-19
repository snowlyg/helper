package global

import (
	"strings"
	"testing"
)

var network = "10.0.0.1/22"

func TestLocalIP(t *testing.T) {
	want := "172.18.236.240"
	t.Run("test get local ip", func(t *testing.T) {
		ip := LocalIP(network)
		if ip != want {
			t.Errorf("LocalIP() want get %s but get %s", want, ip)
		}
	})
}

func TestGetMacAddrs(t *testing.T) {
	want := "00155DDB2E65"
	t.Run("test get mac addr", func(t *testing.T) {
		mac := GetMacAddr()
		if mac != strings.ToUpper(want) {
			t.Errorf("GetMacAddr() want get %s but get %s", want, mac)
		}
	})
}

func TestGetMacAddrInterface(t *testing.T) {
	want := "eth0"
	t.Run("test get mac addr interface", func(t *testing.T) {
		mai := getMacAddrInterface()
		if mai == nil {
			t.Error("mac addr interface is nil")
			return
		}
		if mai.Name != want {
			t.Errorf("getMacAddrInterface() want get %s but get %s", want, mai.Name)
		}
	})
}
func TestCheck(t *testing.T) {
	t.Run("test ip check", func(t *testing.T) {
		if !check("10.0.1.1", network) {
			t.Error("ip check is fail")
			return
		}
		if check("192.168.0.1", network) {
			t.Error("ip check is fail")
			return
		}
	})
}

func TestIsPortInUse(t *testing.T) {
	t.Run("test is port in use", func(t *testing.T) {
		if !IsPortInUse("10.0.0.26", 9092) {
			t.Errorf("IsPortInUse(9092) must be return ture but return false")
		}
	})
}
