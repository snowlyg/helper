package service

import (
	"errors"
	"fmt"
)

// Status represents service status as an byte value
type Status byte

// Status of service represented as an byte
const (
	StatusUnknown Status = iota // Status is unable to be determined due to an error or it was not installed.
	StatusRunning
	StatusStopped
	StatusUninstall
)

const (
	optionKeepAlive            = "KeepAlive"
	optionKeepAliveDefault     = true
	optionRunAtLoad            = "RunAtLoad"
	optionRunAtLoadDefault     = false
	optionUserService          = "UserService"
	optionUserServiceDefault   = false
	optionSessionCreate        = "SessionCreate"
	optionSessionCreateDefault = false
	optionLogOutput            = "LogOutput"
	optionLogOutputDefault     = false
	optionPrefix               = "Prefix"
	optionPrefixDefault        = "application"

	optionRunWait            = "RunWait"
	optionReloadSignal       = "ReloadSignal"
	optionPIDFile            = "PIDFile"
	optionLimitNOFILE        = "LimitNOFILE"
	optionLimitNOFILEDefault = -1 // -1 = don't set in configuration
	optionRestart            = "Restart"

	optionSuccessExitStatus = "SuccessExitStatus"

	optionSystemdScript = "SystemdScript"
	optionSysvScript    = "SysvScript"
	optionRCSScript     = "RCSScript"
	optionUpstartScript = "UpstartScript"
	optionLaunchdConfig = "LaunchdConfig"
	optionOpenRCScript  = "OpenRCScript"

	optionLogDirectory = "LogDirectory"
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

var (
	system         System
	systemRegistry []System
)

type Interface interface {
	Start(s Service) error
	Stop(s Service) error
}

// NewService creates a new service based on a service interface and configuration.
func NewService(i Interface, c *Config) (Service, error) {
	if len(c.Name) == 0 {
		return nil, ErrNameFieldRequired
	}
	if system == nil {
		return nil, ErrNoServiceSystemDetected
	}
	return system.New(i, c)
}

// ChooseSystem chooses a system from the given system services.
// SystemServices are considered in the order they are suggested.
// Calling this may change what Interactive and Platform return.
func ChooseSystem(a ...System) {
	systemRegistry = a
	system = newSystem()
}

// ChosenSystem returns the system that service will use.
func ChosenSystem() System {
	return system
}

// AvailableSystems returns the list of system services considered
// when choosing the system service.
func AvailableSystems() []System {
	return systemRegistry
}

type System interface {
	// String returns a description of the system.
	String() string

	// Detect returns true if the system is available to use.
	Detect() bool

	// Interactive returns false if running under the system service manager
	// and true otherwise.
	Interactive() bool

	// New creates a new service for this system.
	New(i Interface, c *Config) (Service, error)
}

type Config struct {
	Name        string   // Required name of the service. No spaces suggested.
	DisplayName string   // Display name, spaces allowed.
	Description string   // Long description of service.
	UserName    string   // Run as username.
	Arguments   []string // Run with arguments.

	// Optional field to specify the executable for service.
	// If empty the current executable is used.
	Executable string

	// Array of service dependencies.
	// Not yet fully implemented on Linux or OS X:
	//  1. Support linux-systemd dependencies, just put each full line as the
	//     element of the string array, such as
	//     "After=network.target syslog.target"
	//     "Requires=syslog.target"
	//     Note, such lines will be directly appended into the [Unit] of
	//     the generated service config file, will not check their correctness.
	Dependencies []string

	// The following fields are not supported on Windows.
	WorkingDirectory string // Initial working directory.
	ChRoot           string

	// System specific options.
	Option KeyValue

	EnvVars map[string]string
}

type KeyValue map[string]interface{}

// bool returns the value of the given name, assuming the value is a boolean.
// If the value isn't found or is not of the type, the defaultValue is returned.
func (kv KeyValue) bool(name string, defaultValue bool) bool {
	if v, found := kv[name]; found {
		if castValue, is := v.(bool); is {
			return castValue
		}
	}
	return defaultValue
}

// int returns the value of the given name, assuming the value is an int.
// If the value isn't found or is not of the type, the defaultValue is returned.
func (kv KeyValue) int(name string, defaultValue int) int {
	if v, found := kv[name]; found {
		if castValue, is := v.(int); is {
			return castValue
		}
	}
	return defaultValue
}

// string returns the value of the given name, assuming the value is a string.
// If the value isn't found or is not of the type, the defaultValue is returned.
func (kv KeyValue) string(name string, defaultValue string) string {
	if v, found := kv[name]; found {
		if castValue, is := v.(string); is {
			return castValue
		}
	}
	return defaultValue
}

// float64 returns the value of the given name, assuming the value is a float64.
// If the value isn't found or is not of the type, the defaultValue is returned.
func (kv KeyValue) float64(name string, defaultValue float64) float64 {
	if v, found := kv[name]; found {
		if castValue, is := v.(float64); is {
			return castValue
		}
	}
	return defaultValue
}

// funcSingle returns the value of the given name, assuming the value is a func().
// If the value isn't found or is not of the type, the defaultValue is returned.
func (kv KeyValue) funcSingle(name string, defaultValue func()) func() {
	if v, found := kv[name]; found {
		if castValue, is := v.(func()); is {
			return castValue
		}
	}
	return defaultValue
}

// Platform returns a description of the system service.
func Platform() string {
	if system == nil {
		return ""
	}
	return system.String()
}

// Interactive returns false if running under the OS service manager
// and true otherwise.
func Interactive() bool {
	if system == nil {
		return true
	}
	return system.Interactive()
}

func newSystem() System {
	for _, choice := range systemRegistry {
		if choice.Detect() == false {
			continue
		}
		return choice
	}
	return nil
}

type Service interface {
	// Run should be called shortly after the program entry point.
	// After Interface.Stop has finished running, Run will stop blocking.
	// After Run stops blocking, the program must exit shortly after.
	Run() error

	// Start signals to the OS service manager the given service should start.
	Start() error

	// Stop signals to the OS service manager the given service should stop.
	Stop() error

	// Restart signals to the OS service manager the given service should stop then start.
	Restart() error

	// Install setups up the given service in the OS service manager. This may require
	// greater rights. Will return an error if it is already installed.
	Install() error

	// Uninstall removes the given service from the OS service manager. This may require
	// greater rights. Will return an error if the service is not present.
	Uninstall() error

	// Opens and returns a system logger. If the user program is running
	// interactively rather then as a service, the returned logger will write to
	// os.Stderr. If errs is non-nil errors will be sent on errs as well as
	// returned from Logger's functions.
	Logger(errs chan<- error) (Logger, error)

	// SystemLogger opens and returns a system logger. If errs is non-nil errors
	// will be sent on errs as well as returned from Logger's functions.
	SystemLogger(errs chan<- error) (Logger, error)

	// String displays the name of the service. The display name if present,
	// otherwise the name.
	String() string

	// Platform displays the name of the system that manages the service.
	// In most cases this will be the same as service.Platform().
	Platform() string

	// Status returns the current service status.
	Status() (Status, error)
}

// Logger writes to the system log.
type Logger interface {
	Error(v ...interface{}) error
	Warning(v ...interface{}) error
	Info(v ...interface{}) error

	Errorf(format string, a ...interface{}) error
	Warningf(format string, a ...interface{}) error
	Infof(format string, a ...interface{}) error
}

// ControlAction list valid string texts to use in Control.
var ControlAction = [5]string{"start", "stop", "restart", "install", "uninstall"}

// Control issues control functions to the service from a given action string.
func Control(s Service, action string) error {
	var err error
	switch action {
	case ControlAction[0]:
		err = s.Start()
	case ControlAction[1]:
		err = s.Stop()
	case ControlAction[2]:
		err = s.Restart()
	case ControlAction[3]:
		err = s.Install()
	case ControlAction[4]:
		err = s.Uninstall()
	default:
		err = fmt.Errorf("Unknown action %s", action)
	}
	if err != nil {
		return fmt.Errorf("Failed to %s %v: %v", action, s, err)
	}
	return nil
}
