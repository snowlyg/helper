package global

import "testing"

func TestLocalIP(t *testing.T) {
	want := "10.0.1.171"
	t.Run("test get local ip", func(t *testing.T) {
		ip := LocalIP()
		if ip != want {
			t.Errorf("LocalIP() want get %s but get %s", want, ip)
		}
	})
}

func TestGetMacAddrs(t *testing.T) {
	want := "C85ACF0D1B14"
	t.Run("test get mac addr", func(t *testing.T) {
		mac := GetMacAddr()
		if mac != want {
			t.Errorf("GetMacAddr() want get %s but get %s", want, mac)
		}
	})
}
