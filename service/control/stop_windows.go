package control

import (
	"fmt"

	"github.com/snowlyg/helper/service"
)

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
