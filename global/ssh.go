package global

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/crypto/ssh"
)

type Cli struct {
	IP         string      //IP地址
	Username   string      //用户名
	Password   string      //密码
	Port       int         //端口号
	client     *ssh.Client //ssh客户端
	LastResult string      //最近一次Run的结果
}

//创建命令行对象
//@param ip IP地址
//@param username 用户名
//@param password 密码
//@param port 端口号,默认22
func NewSSH(ip string, username string, password string, port ...int) *Cli {
	cli := new(Cli)
	cli.IP = ip
	cli.Username = username
	cli.Password = password
	if len(port) <= 0 {
		cli.Port = 22
	} else {
		cli.Port = port[0]
	}
	return cli
}

//执行shell
//@param shell shell脚本命令
func (c Cli) Run(shell string) (string, error) {
	if c.client == nil {
		if err := c.connect(); err != nil {
			return "", fmt.Errorf("连接失败 %w", err)
		}
	}
	defer c.client.Close()
	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("新建会话失败 %w", err)
	}
	defer session.Close()
	buf, err := session.CombinedOutput(shell)
	if err != nil {
		return "", fmt.Errorf("执行命令失败 %w", err)
	}

	c.LastResult = string(buf)
	return c.LastResult, nil
}

//连接
func (c *Cli) connect() error {
	config := ssh.ClientConfig{
		User: c.Username,
		Auth: []ssh.AuthMethod{ssh.Password(c.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", c.IP, c.Port)
	sshClient, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return err
	}
	c.client = sshClient
	return nil
}

type DeviceMem struct {
	Free      decimal.Decimal
	Total     decimal.Decimal
	FreeRound decimal.Decimal
}

// GetMem 内存信息
func (c *Cli) GetMem() (DeviceMem, error) {
	var deviceMen DeviceMem
	totalDec, err := c.getMenInfo("cat /proc/meminfo | grep -w MemTotal", "MemTotal:")
	if err != nil {
		return deviceMen, fmt.Errorf("memTotal %w", err)
	}
	freeDec, err := c.getMenInfo("cat /proc/meminfo | grep -w MemFree", "MemFree:")
	if err != nil {
		return deviceMen, fmt.Errorf("memFree %w", err)
	}
	buffersDec, err := c.getMenInfo("cat /proc/meminfo | grep -w Buffers", "Buffers:")
	if err != nil {
		return deviceMen, fmt.Errorf("buffers %w", err)
	}
	cachedDec, err := c.getMenInfo("cat /proc/meminfo | grep -w Cached", "Cached:")
	if err != nil {
		return deviceMen, fmt.Errorf("cached %w", err)
	}

	deviceMen.Total = totalDec
	deviceMen.Free = freeDec.Add(buffersDec).Add(cachedDec)
	deviceMen.FreeRound = deviceMen.Free.DivRound(deviceMen.Total, 2)
	return deviceMen, nil
}

// getMenInfo 内存信息
func (c *Cli) getMenInfo(cmd, key string) (decimal.Decimal, error) {
	info, err := c.Run(cmd)
	if err != nil {
		return decimal.Zero, fmt.Errorf("%s %w", cmd, err)
	}

	dec, err := decimal.NewFromString(strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(info, key, ""), "kB", "")))
	if err != nil {
		return decimal.Zero, fmt.Errorf("NewFromString %w", err)
	}
	return dec, err
}

// GetDf 硬盘率
func (c *Cli) GetDf() (string, error) {
	total, err := c.Run("df /sdcard -h")
	if err != nil {
		return "", fmt.Errorf("df /sdcard -h  %w", err)
	}
	flysnowRegexp := regexp.MustCompile(`(100|[1-9]?\d(\.\d\d?\d?)?)%`)
	params := flysnowRegexp.FindStringSubmatch(total)

	if len(params) > 0 {
		return params[0], nil
	}
	return "", nil
}

// GetSignal wifi信号
func (c *Cli) GetSignal() (string, error) {
	signal, err := c.Run("iw dev wlan0 link | grep -w signal:")
	if err != nil {
		return "", fmt.Errorf("iw dev wlan0 link | grep -w signal:  %w", err)
	}

	return strings.TrimSpace(strings.ReplaceAll(signal, "signal:", "")), nil
}

// GetDatetime 当前时间
func (c *Cli) GetDatetime() (string, error) {
	datetime, err := c.Run("date +'%Y/%m/%d %T %Z'")
	if err != nil {
		return "", fmt.Errorf("date %w", err)
	}

	return datetime, nil
}

// GetCpuTemp CPU温度
func (c *Cli) GetCpuTemp() (float64, error) {
	cpuTemp, err := c.Run("cat /sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		return 0, fmt.Errorf("cat /sys/class/thermal/thermal_zone0/temp %w", err)
	}
	f, err := strconv.ParseFloat(strings.Trim(cpuTemp, "\n"), 64)
	if err != nil {
		return 0, fmt.Errorf("ParseFloat %w", err)
	}
	return f / 1000, nil
}
