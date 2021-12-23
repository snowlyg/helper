package control

import (
	"fmt"

	"github.com/snowlyg/helper/service"
)

func Install(srvName, execPath, displayName, systemName, pwd string) error {
	status, err := service.ServiceStatus(srvName)
	if err != nil {
		return fmt.Errorf("get error msg %w", err)
	}

	if status == service.StatusRunning {
		return fmt.Errorf("%s is running", srvName)
	}

	if status == service.StatusUninstall {
		return service.ServiceInstall(srvName, execPath, displayName, systemName, pwd)
	}

	return nil
}

func Status(srvName string) (service.Status, error) {
	status, err := service.ServiceStatus(srvName)
	if err != nil {
		return service.StatusUnknown, fmt.Errorf("get service status  %w", err)
	}
	return status, nil
}

func Stop(srvName string) error {
	status, err := service.ServiceStatus(srvName)
	if err != nil {
		return fmt.Errorf("get error msg %w", err)
	}

	if status != service.StatusRunning {
		return nil
	}

	return service.ServiceStop(srvName)
}

func Start(srvName string) error {
	status, err := service.ServiceStatus(srvName)
	if err != nil {
		return fmt.Errorf("get service status  %w", err)
	}

	if status == service.StatusRunning {
		return nil
	}

	if status == service.StatusUninstall {
		return fmt.Errorf("service uninstall")
	}

	if status == service.StatusStopped {
		return service.ServiceStart(srvName)
	}

	return nil
}

func Uninstall(srvName string) error {
	status, err := service.ServiceStatus(srvName)
	if err != nil {
		return fmt.Errorf("service status get error %w", err)
	}

	if status == service.StatusUninstall {
		return nil
	}

	return service.ServiceUninstall(srvName)
}
