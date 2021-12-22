package control

import (
	"fmt"

	"github.com/snowlyg/helper/service"
)

func Status(srvName string) (service.Status, error) {
	status, err := service.ServiceStatus(srvName)
	if err != nil {
		return service.StatusUnknown, fmt.Errorf("get service status  %w", err)
	}
	return status, nil
}
