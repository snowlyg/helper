package service

import (
	"golang.org/x/sys/windows/svc/mgr"
)

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
