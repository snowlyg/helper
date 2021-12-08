package control

import (
	"fmt"

	"github.com/snowlyg/helper/service"
)

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
