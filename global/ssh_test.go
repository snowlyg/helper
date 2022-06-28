package global

import (
	_ "embed"
	"os"
	"testing"
)

func Test_NewSSH(t *testing.T) {
	sshPwd := os.Getenv("sshPwd")
	t.Run("测试新建ssh链接", func(t *testing.T) {
		ip := "10.0.0.26"
		name := "root"
		sshClient := NewSSH(ip, name, sshPwd, 2222)
		if sshClient == nil {
			t.Errorf("客户端为空")
			return
		}
	})
}
func Test_GetMem(t *testing.T) {
	sshPwd := os.Getenv("sshPwd")
	t.Run("测试获取设备内存", func(t *testing.T) {
		ip := "10.0.0.110"
		name := "root"
		sshClient := NewSSH(ip, name, sshPwd, 2222)
		if sshClient == nil {
			t.Errorf("客户端为空")
			return
		}
		_, err := sshClient.GetMem()
		if err != nil {
			t.Errorf(err.Error())
			return
		}
	})
}
func Test_GetDf(t *testing.T) {
	sshPwd := os.Getenv("sshPwd")
	t.Run("测试获取设备硬盘使用", func(t *testing.T) {
		ip := "10.0.0.110"
		name := "root"
		sshClient := NewSSH(ip, name, sshPwd, 2222)
		if sshClient == nil {
			t.Errorf("客户端为空")
			return
		}
		_, err := sshClient.GetDf()
		if err != nil {
			t.Errorf(err.Error())
			return
		}
	})
}
func Test_GetSignal(t *testing.T) {
	sshPwd := os.Getenv("sshPwd")
	t.Run("测试获取设备信号使用", func(t *testing.T) {
		ip := "10.0.0.110"
		name := "root"
		sshClient := NewSSH(ip, name, sshPwd, 2222)
		if sshClient == nil {
			t.Errorf("客户端为空")
			return
		}
		_, err := sshClient.GetSignal()
		if err != nil {
			t.Errorf(err.Error())
			return
		}
	})
}

func Test_GetDatetime(t *testing.T) {
	sshPwd := os.Getenv("sshPwd")
	t.Run("测试获取设备时间使用", func(t *testing.T) {
		ip := "10.0.0.110"
		name := "root"
		sshClient := NewSSH(ip, name, sshPwd, 2222)
		if sshClient == nil {
			t.Errorf("客户端为空")
			return
		}
		_, err := sshClient.GetDatetime()
		if err != nil {
			t.Errorf(err.Error())
			return
		}
	})
}
func Test_GetCpuTemp(t *testing.T) {
	sshPwd := os.Getenv("sshPwd")
	t.Run("测试获取设备时间使用", func(t *testing.T) {
		ip := "10.0.0.177"
		name := "root"
		sshClient := NewSSH(ip, name, sshPwd, 2222)
		if sshClient == nil {
			t.Errorf("客户端为空")
			return
		}
		cpuTemp, err := sshClient.GetCpuTemp()
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if cpuTemp <= 0 {
			t.Errorf("cpu temp is 0")
		}
	})
}
