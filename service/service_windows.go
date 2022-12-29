//go:build windows
// +build windows

// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

var version = "windows-service"

var interactive = false

// IsWindowsService 不能修改成 IsWindowsService 否则无法正常启动应用
func init() {
	ChooseSystem(windowsSystem{})
	var err error
	interactive, err = svc.IsWindowsService()
	if err != nil {
		panic(err)
	}
}

type windowsSystem struct{}

func (windowsSystem) String() string {
	return version
}

func (windowsSystem) Detect() bool {
	return true
}
func (windowsSystem) Interactive() bool {
	return interactive
}
func (windowsSystem) New(i Interface, c *Config) (Service, error) {
	ws := &WindowsService{
		i: i,
		c: c,
	}
	return ws, nil
}

type WindowsService struct {
	c *Config
	i Interface
}

func NewWindowsService(i Interface, c *Config) (*WindowsService, error) {
	ws := &WindowsService{
		i: i,
		c: c,
	}
	return ws, nil
}

func (ws *WindowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (bool, uint32) {
	elog, err := eventlog.Open(ws.c.Name)
	if err != nil {
		return true, 1
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("%s service execute", ws.c.Name))

	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}

	if err := ws.i.Start(nil); err != nil {
		elog.Info(1, fmt.Sprintf("%s service start failed: %v", ws.c.Name, err))
		return true, 1
	}

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		c := <-r
		switch c.Cmd {
		case svc.Interrogate:
			elog.Info(1, fmt.Sprintf("%s service current status", ws.c.Name))
			changes <- c.CurrentStatus
		case svc.Stop:
			elog.Info(1, fmt.Sprintf("%s service stop", ws.c.Name))
			changes <- svc.Status{State: svc.StopPending}
			if err := ws.i.Stop(nil); err != nil {
				elog.Info(1, fmt.Sprintf("%s service stop failed: %v", ws.c.Name, err))
				return true, 2
			}
			break loop
		case svc.Shutdown:
			elog.Info(1, fmt.Sprintf("%s service shutdown", ws.c.Name))
			changes <- svc.Status{State: svc.StopPending}
			err := ws.i.Stop(nil)
			if err != nil {
				elog.Info(1, fmt.Sprintf("%s service shutdown failed: %v", ws.c.Name, err))
				return true, 2
			}
			break loop
		default:
			continue loop
		}
	}

	return false, 0
}

func (ws *WindowsService) Run() error {
	if interactive {

		elog, err := eventlog.Open(ws.c.Name)
		if err != nil {
			return err
		}
		defer elog.Close()

		elog.Info(1, fmt.Sprintf("starting %s service", ws.c.Name))

		run := svc.Run
		err = run(ws.c.Name, ws)
		if err != nil {
			elog.Error(1, fmt.Sprintf("%s service run failed: %v", ws.c.Name, err))
			return err
		}

	} else {
		err := ws.i.Start(nil)
		if err != nil {
			return err
		}

		sigChan := make(chan os.Signal, 1)

		signal.Notify(sigChan, os.Interrupt)

		<-sigChan

		return ws.i.Stop(nil)
	}
	return nil
}

func ServiceInstall(svcName, execPath, dispalyName, serviceStartName, pwd string, args ...string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(svcName)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", svcName)
	}
	s, err = m.CreateService(svcName, execPath, mgr.Config{
		DisplayName:      dispalyName,
		StartType:        mgr.StartAutomatic,
		ServiceStartName: serviceStartName,
		Password:         pwd,
	}, args...)
	if err != nil {
		return fmt.Errorf("create service failed: %s", err)
	}
	defer s.Close()
	err = eventlog.InstallAsEventCreate(svcName, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("install as event create failed: %s", err)
	}
	return nil
}

// Start signals to the OS service manager the given service should start.
func (ws *WindowsService) Start() error {
	return ServiceStart(ws.c.Name)
}

// Stop signals to the OS service manager the given service should stop.
func (ws *WindowsService) Stop() error {
	return ServiceStop(ws.c.Name)
}

// Restart signals to the OS service manager the given service should stop then start.
func (ws *WindowsService) Restart() error {
	if err := ServiceStop(ws.c.Name); err != nil {
		return err
	}
	return ServiceStart(ws.c.Name)
}

// Install setups up the given service in the OS service manager. This may require
// greater rights. Will return an error if it is already installed.
func (ws *WindowsService) Install() error {
	return nil
}

// Uninstall removes the given service from the OS service manager. This may require
// greater rights. Will return an error if the service is not present.
func (ws *WindowsService) Uninstall() error {
	elog, err := eventlog.Open(ws.c.Name)
	if err != nil {
		return err
	}
	defer elog.Close()

	if err := ServiceUninstall(ws.c.Name); err != nil {
		elog.Error(1, fmt.Sprintf("%s service uninstall failed: %v", ws.c.Name, err))
		return err
	}
	elog.Info(1, fmt.Sprintf("uninstall %s service", ws.c.Name))
	return nil
}

// Opens and returns a system logger. If the user program is running
// interactively rather then as a service, the returned logger will write to
// os.Stderr. If errs is non-nil errors will be sent on errs as well as
// returned from Logger's functions.
func (ws *WindowsService) Logger(errs chan<- error) (Logger, error) {
	return nil, nil
}

// SystemLogger opens and returns a system logger. If errs is non-nil errors
// will be sent on errs as well as returned from Logger's functions.
func (ws *WindowsService) SystemLogger(errs chan<- error) (Logger, error) {
	return nil, nil
}

// String displays the name of the service. The display name if present,
// otherwise the name.
func (ws *WindowsService) String() string {
	return version
}

// Platform displays the name of the system that manages the service.
// In most cases this will be the same as service.Platform().
func (ws *WindowsService) Platform() string {
	return ws.String()
}

// Status returns the current service status.
func (ws *WindowsService) Status() (Status, error) {
	return ServiceStatus(ws.c.Name)
}

func ServiceUninstall(srcName string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(srcName)
	if err != nil {
		return fmt.Errorf("service %s is not installed", srcName)
	}
	defer s.Close()
	err = s.Delete()
	if err != nil {
		return err
	}
	err = eventlog.Remove(srcName)
	if err != nil {
		return fmt.Errorf("RemoveEventLogSource() failed: %s", err)
	}

	return nil
}

func ServiceStop(srcName string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(srcName)
	if err != nil {
		return err
	}
	defer s.Close()

	return stopWait(s)
}

func stopWait(s *mgr.Service) error {
	// First stop the service. Then wait for the service to
	// actually stop before starting it.
	status, err := s.Control(svc.Stop)
	if err != nil {
		return err
	}

	timeDuration := time.Millisecond * 50

	timeout := time.After(getStopTimeout() + (timeDuration * 2))
	tick := time.NewTicker(timeDuration)
	defer tick.Stop()

	for status.State != svc.Stopped {
		select {
		case <-tick.C:
			status, err = s.Query()
			if err != nil {
				return err
			}
		case <-timeout:
			break
		}
	}
	return nil
}

// getStopTimeout fetches the time before windows will kill the service.
func getStopTimeout() time.Duration {
	// For default and paths see https://support.microsoft.com/en-us/kb/146092
	defaultTimeout := time.Millisecond * 20000
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control`, registry.READ)
	if err != nil {
		return defaultTimeout
	}
	sv, _, err := key.GetStringValue("WaitToKillServiceTimeout")
	if err != nil {
		return defaultTimeout
	}
	v, err := strconv.Atoi(sv)
	if err != nil {
		return defaultTimeout
	}
	return time.Millisecond * time.Duration(v)
}

// status
func ServiceStatus(srcName string) (Status, error) {
	m, err := mgr.Connect()
	if err != nil {
		return StatusUnknown, err
	}
	defer m.Disconnect()

	s, err := m.OpenService(srcName)
	if err != nil {
		if err.Error() == "The specified service does not exist as an installed service." {
			return StatusUninstall, nil
		}
		return StatusUnknown, err
	}
	defer s.Close()

	status, err := s.Query()
	if err != nil {
		return StatusUnknown, err
	}
	switch status.State {
	case svc.StartPending:
		fallthrough
	case svc.Running:
		return StatusRunning, nil
	case svc.PausePending:
		fallthrough
	case svc.Paused:
		fallthrough
	case svc.ContinuePending:
		fallthrough
	case svc.StopPending:
		fallthrough
	case svc.Stopped:
		return StatusStopped, nil
	default:
		return StatusUnknown, fmt.Errorf("unknown status %v", status)
	}
}

// status
func ServiceProcessId(srcName string) (uint32, error) {
	var processId uint32
	m, err := mgr.Connect()
	if err != nil {
		return processId, err
	}
	defer m.Disconnect()

	s, err := m.OpenService(srcName)
	if err != nil {
		if err.Error() == "The specified service does not exist as an installed service." {
			return processId, nil
		}
		return processId, err
	}
	defer s.Close()

	status, err := s.Query()
	if err != nil {
		return processId, err
	}

	processId = status.ProcessId

	switch status.State {
	case svc.StartPending:
		fallthrough
	case svc.Running:
		processId = status.ProcessId
		return processId, nil
	case svc.PausePending:
		fallthrough
	case svc.Paused:
		fallthrough
	case svc.ContinuePending:
		fallthrough
	case svc.StopPending:
		fallthrough
	case svc.Stopped:
		return processId, nil
	default:
		return processId, fmt.Errorf("unknown status %v", status)
	}
}

func ServiceStart(srcName string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(srcName)
	if err != nil {
		return err
	}
	defer s.Close()
	return s.Start()
}
