// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"fmt"
	"os"
	"os/signal"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

var elog debug.Log

var interactive = false

// IsAnInteractiveSession 不能修改成 IsWindowsService 否则无法正常启动应用
func init() {
	var err error
	interactive, err = svc.IsAnInteractiveSession()
	if err != nil {
		elog.Info(1, fmt.Sprintf("IsAnInteractiveSession failed: %v", err))
		panic(err)
	}
}

type WindowsService struct {
	Name string
	i    Interface
}

func NewService(i Interface, name string) (*WindowsService, error) {
	ws := &WindowsService{
		i:    i,
		Name: name,
	}
	return ws, nil
}

func (ws *WindowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (bool, uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}

	if err := ws.i.Start(); err != nil {
		elog.Info(1, fmt.Sprintf("%s service start failed: %v", ws.Name, err))
		return true, 1
	}

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		c := <-r
		switch c.Cmd {
		case svc.Interrogate:
			changes <- c.CurrentStatus
		case svc.Stop:
			changes <- svc.Status{State: svc.StopPending}
			if err := ws.i.Stop(); err != nil {
				elog.Info(1, fmt.Sprintf("%s service stop failed: %v", ws.Name, err))
				return true, 2
			}
			break loop
		case svc.Shutdown:
			changes <- svc.Status{State: svc.StopPending}
			err := ws.i.Stop()
			if err != nil {
				elog.Info(1, fmt.Sprintf("%s service shutdown failed: %v", ws.Name, err))
				return true, 2
			}
			break loop
		default:
			continue loop
		}
	}

	return false, 0
}

func (ws *WindowsService) Run(isDebug bool) error {
	var err error
	if !interactive {
		if isDebug {
			elog = debug.New(ws.Name)
		} else {
			elog, err = eventlog.Open(ws.Name)
			if err != nil {
				return err
			}
		}
		defer elog.Close()

		elog.Info(1, fmt.Sprintf("starting %s service", ws.Name))
		run := svc.Run
		if isDebug {
			run = debug.Run
		}
		err = run(ws.Name, ws)
		if err != nil {
			elog.Error(1, fmt.Sprintf("%s service failed: %v", ws.Name, err))
			return err
		}
		elog.Info(1, fmt.Sprintf("%s service stopped", ws.Name))
	}

	err = ws.i.Start()
	if err != nil {
		return err
	}

	sigChan := make(chan os.Signal)

	signal.Notify(sigChan, os.Interrupt)

	<-sigChan

	return ws.i.Stop()
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
		return fmt.Errorf("CreateService() failed: %s", err)
	}
	defer s.Close()
	err = eventlog.InstallAsEventCreate(svcName, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("InstallAsEventCreate() failed: %s", err)
	}
	return nil
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
