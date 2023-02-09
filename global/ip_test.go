package global

import (
	"strings"
	"testing"
)

func TestLocalIP(t *testing.T) {
	want := "10.0.0.26"
	t.Run("test get local ip", func(t *testing.T) {
		ip := LocalIP()
		if ip != want {
			t.Errorf("LocalIP() want get %s but get %s", want, ip)
		}
	})
}

func TestGetMacAddrs(t *testing.T) {
	want := "000c2978391b"
	t.Run("test get mac addr", func(t *testing.T) {
		mac := GetMacAddr()
		if mac != strings.ToUpper(want) {
			t.Errorf("GetMacAddr() want get %s but get %s", want, mac)
		}
	})
}

func TestGetMacAddrInterface(t *testing.T) {
	want := "ens192"
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

func TestIsPortInUse(t *testing.T) {
	t.Run("test is port in use", func(t *testing.T) {
		if !IsPortInUse("127.0.0.1", 9092) {
			t.Errorf("IsPortInUse(9092) must be return ture but return false")
		}
	})
}
