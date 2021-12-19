package service

import (
	"errors"
	"fmt"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

type Status byte

const (
	StatusUnknown Status = iota // Status is unable to be determined due to an error or it was not installed.
	StatusRunning
	StatusStopped
	StatusUninstall
)

var (
	// ErrNameFieldRequired is returned when Config.Name is empty.
	ErrNameFieldRequired = errors.New("Config.Name field is required.")
	// ErrNoServiceSystemDetected is returned when no system was detected.
	ErrNoServiceSystemDetected = errors.New("No service system detected.")
	// ErrNotInstalled is returned when the service is not installed
	ErrNotInstalled = errors.New("the service is not installed")
)

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
