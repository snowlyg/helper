package control

import (
	"errors"
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

	restop := 3
	for restop > 0 {
		err := service.ServiceStop(srvName)
		if err != nil {
			fmt.Println(err)
		}
		status, _ = service.ServiceStatus(srvName)
		if status == service.StatusStopped {
			return nil
		}
		restop--
		fmt.Println("停止失败1次")
		continue
	}

	return errors.New("停止失败")
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

	restart := 3
	for restart > 0 {
		err := service.ServiceStart(srvName)
		if err != nil {
			fmt.Println(err)
		}
		status, _ = service.ServiceStatus(srvName)
		if status == service.StatusRunning {
			return nil
		}
		restart--
		fmt.Println("启动失败1次")
		continue
	}

	return errors.New("启动失败")
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
