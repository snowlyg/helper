package service

import "errors"

// Status represents service status as an byte value
type Status byte

// Status of service represented as an byte
const (
	StatusUnknown Status = iota // Status is unable to be determined due to an error or it was not installed.
	StatusRunning
	StatusStopped
	StatusUninstall
)

var (
	// ErrNameFieldRequired is returned when Config.Name is empty.
	ErrNameFieldRequired = errors.New("缺失配置文件")
	// ErrNoServiceSystemDetected is returned when no system was detected.
	ErrNoServiceSystemDetected = errors.New("未检测到服务系统")
	// ErrNotInstalled is returned when the service is not installed.
	ErrNotInstalled = errors.New("服务未安装")
	ErrNotSuport    = errors.New("服务不支持")
)

type Interface interface {
	Start() error
	Stop() error
}
