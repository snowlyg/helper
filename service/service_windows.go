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

type Interface interface {
	Start() error
	Stop() error
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
