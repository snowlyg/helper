// +build windows
package control

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/snowlyg/helper/dir"
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

func ProcessId(srvName string) (uint32, error) {
	processId, err := service.ServiceProcessId(srvName)
	if err != nil {
		return processId, err
	}
	return processId, nil
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
		go func() {
			err := service.ServiceStop(srvName)
			if err != nil {
				fmt.Println(err)
			}
		}()
		fmt.Println(restop)
		time.Sleep(1 * time.Second)
		restop--
	}

	status, err = service.ServiceStatus(srvName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if status != service.StatusStopped {
		return errors.New("服务未停止")
	}

	pid, err := dir.ReadString("./.pid")
	if err != nil {
		return err
	}

	ppid, _ := strconv.Atoi(pid)
	if b, _ := process.PidExists(int32(ppid)); b {
		ps, _ := process.Processes()
		if len(ps) > 0 {
			for _, p := range ps {
				if p.Pid != int32(ppid) {
					continue
				}
				err = p.Kill()
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	return nil
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
			processId, err := service.ServiceProcessId(srvName)
			if err == nil {
				dir.WriteString("./.pid", strconv.FormatUint(uint64(processId), 10))
			}
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
