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
