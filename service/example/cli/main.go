// +build windows
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/snowlyg/helper/service/control"
)

func usage(errmsg string) {
	fmt.Fprintf(os.Stderr,
		"%s\n\n"+
			"usage: %s <command> <servicename>\n"+
			"       where <command> is one of\n"+
			"       install, remove, start, stop, status .\n"+
			"       and use install : .\n"+
			"       install <service name> <exec path> <display name> <system name> <password>  \n",
		errmsg, os.Args[0])
	os.Exit(2)
}

func main() {
	if len(os.Args) < 2 {
		usage("no command specified")
	}
	srvName := strings.ToLower(os.Args[2])
	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "start":
		err := control.Start(srvName)
		if err != nil {
			fmt.Printf("%v \n", err)
		}
		println("start success")
	case "install":

		if len(os.Args) != 7 {
			usage("no command specified")
		}
		err := control.Install(srvName, os.Args[3], os.Args[4], os.Args[5], os.Args[6])
		if err != nil {
			fmt.Printf("%v \n", err)
		}
		println("install success")
	case "stop":
		err := control.Stop(srvName)
		if err != nil {
			fmt.Printf("%v \n", err)
		}
		println("stop success")
	case "remove":
		err := control.Uninstall(srvName)
		if err != nil {
			fmt.Printf("%v \n", err)
		}
		println("remove success")
	case "status":
		println(control.Status(srvName))
	default:
		println("invaild command")
	}
}
