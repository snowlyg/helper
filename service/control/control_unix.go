//go:build linux || darwin || solaris || aix || freebsd
// +build linux darwin solaris aix freebsd

package control

import (
	"github.com/snowlyg/helper/service"
)

func (p *program) Start(s service.Service) error {
	// do some work
	return nil
}

func (p *program) Stop(s service.Service) error {
	//stop
	return nil
}

type program struct{}

func Install(srvName, execPath, displayName, systemName, pwd string) error {
	s, err := service.NewService(&program{}, &service.Config{Name: srvName})
	if err != nil {
		return err
	}
	return s.Install()
}

func Status(srvName string) (service.Status, error) {
	s, err := service.NewService(&program{}, &service.Config{Name: srvName})
	if err != nil {
		return service.StatusUnknown, err
	}
	return s.Status()
}

func ProcessId(srvName string) (uint32, error) {
	return 0, nil
}

func Stop(srvName string) error {
	s, err := service.NewService(&program{}, &service.Config{Name: srvName})
	if err != nil {
		return err
	}
	return s.Stop()
}

func Start(srvName string) error {
	s, err := service.NewService(&program{}, &service.Config{Name: srvName})
	if err != nil {
		return err
	}
	return s.Start()
}

func Uninstall(srvName string) error {
	s, err := service.NewService(&program{}, &service.Config{Name: srvName})
	if err != nil {
		return err
	}
	return s.Uninstall()
}
