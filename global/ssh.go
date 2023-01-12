package global

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"log"

	"github.com/shopspring/decimal"
	"golang.org/x/crypto/ssh"
)

var (
	ErrConnectFail    = errors.New("SSH 连接失败")
	ErrNewSessionFail = errors.New("SSH 新建会话失败")
	ErrRunCommandFail = errors.New("SSH 执行命令失败")
)

type Cli struct {
	IP         string      //IP地址
	Username   string      //用户名
	Password   string      //密码
	Port       int         //端口号
	client     *ssh.Client //ssh客户端
	LastResult string      //最近一次Run的结果
	Debug      bool
}

// 创建命令行对象
// @param ip IP地址
// @param username 用户名
// @param password 密码
// @param debug 是否调试
// @param port 端口号,默认22
func NewSSH(ip string, username string, password string, debug bool, port ...int) *Cli {
	cli := new(Cli)
	cli.IP = ip
	cli.Username = username
	cli.Password = password
	cli.Debug = debug
	if len(port) <= 0 {
		cli.Port = 22
	} else {
		cli.Port = port[0]
	}
	return cli
}

// 执行shell
// @param shell shell脚本命令
func (c Cli) Run(shell string) (string, error) {
	if c.client == nil {
		if err := c.connect(); err != nil {
			if c.Debug {
				log.Println(err.Error())
			}
			return "", ErrConnectFail
		}
	}
	defer c.client.Close()
	session, err := c.client.NewSession()
	if err != nil {
		if c.Debug {
			log.Println(err.Error())
		}
		return "", ErrNewSessionFail
	}
	defer session.Close()
	buf, err := session.CombinedOutput(shell)
	if err != nil {
		if c.Debug {
			log.Println(err.Error())
		}
		return "", ErrRunCommandFail
	}

	c.LastResult = string(buf)
	return c.LastResult, nil
}

// 连接
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
		if c.Debug {
			log.Println(err.Error())
		}
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
	var ct float64
	cmd := "acpi -t"
	cpuTemp, err := c.Run(cmd)
	if err != nil {
		return 0, fmt.Errorf("%s : %w", cmd, err)
	}
	cpuTemps := strings.Split(cpuTemp, "\n")

	for _, v := range cpuTemps {
		flysnowRegexp := regexp.MustCompile(`[1-9]\d*\.\d*|0\.\d*[1-9]\d*`)
		params := flysnowRegexp.FindStringSubmatch(v)
		if len(params) > 0 {
			f, err := strconv.ParseFloat(strings.Trim(params[0], "\n"), 64)
			if err != nil {
				return 0, fmt.Errorf("ParseFloat %w", err)
			}
			current := f / 100
			if current > ct {
				ct = current
			}
		}
	}
	return ct, nil
}
