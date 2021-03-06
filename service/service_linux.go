// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"errors"
	"os"
	"os/signal"
	"os/user"
	"syscall"
	"text/template"
)

type systemd struct {
	i    Interface
	Name string
}

func init() {
	// var err error
	// interactive, err = isInteractive()
	// if err != nil {
	// 	panic(err)
	// }
}

func NewService(i Interface, name string) (*systemd, error) {
	s := &systemd{
		i: i,
	}

	return s, nil
}

func (s *systemd) Run() error {
	err := s.i.Start()
	if err != nil {
		return err
	}

	func() {
		var sigChan = make(chan os.Signal, 3)
		signal.Notify(sigChan, syscall.SIGTERM, os.Interrupt)
		<-sigChan
	}()

	return s.i.Stop()
}

func (s *systemd) String() string {
	if len(s.Name) > 0 {
		return s.Name
	}
	return s.Name
}

func (s *systemd) Platform() string {
	return ""
}

func (s *systemd) getHomeDir() (string, error) {
	u, err := user.Current()
	if err == nil {
		return u.HomeDir, nil
	}

	// alternate methods
	homeDir := os.Getenv("HOME") // *nix
	if homeDir == "" {
		return "", errors.New("User home directory not found.")
	}
	return homeDir, nil
}

func (s *systemd) getServiceFilePath() (string, error) {
	// if s.userService {
	// 	homeDir, err := s.getHomeDir()
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	return homeDir + "/Library/LaunchAgents/" + s.Name + ".plist", nil
	// }
	return "/Library/LaunchDaemons/" + s.Name + ".plist", nil
}

func (s *systemd) template() *template.Template {
	functions := template.FuncMap{
		"bool": func(v bool) string {
			if v {
				return "true"
			}
			return "false"
		},
	}

	customConfig := "LaunchdConfig"

	if customConfig != "" {
		return template.Must(template.New("").Funcs(functions).Parse(customConfig))
	} else {
		return template.Must(template.New("").Funcs(functions).Parse(launchdConfig))
	}
}

func ServiceInstall(svcName, execPath, dispalyName, serviceStartName, pwd string, args ...string) error {
	return ErrNotSuport
}

func ServiceUninstall(srcName string) error {
	return ErrNotSuport
}

func ServiceStop(srcName string) error {
	return ErrNotSuport
}

// status
func ServiceStatus(srcName string) (Status, error) {
	return StatusUnknown, ErrNotSuport
}

func ServiceProcessId(srcName string) (uint32, error) {
	return 0, nil
}

func ServiceStart(srcName string) error {
	return ErrNotSuport
}

var launchdConfig = `<?xml version='1.0' encoding='UTF-8'?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN"
"http://www.apple.com/DTDs/PropertyList-1.0.dtd" >
<plist version='1.0'>
  <dict>
    <key>Label</key>
    <string>{{html .Name}}</string>
    <key>ProgramArguments</key>
    <array>
      <string>{{html .Path}}</string>
    {{range .Config.Arguments}}
      <string>{{html .}}</string>
    {{end}}
    </array>
    {{if .UserName}}<key>UserName</key>
    <string>{{html .UserName}}</string>{{end}}
    {{if .ChRoot}}<key>RootDirectory</key>
    <string>{{html .ChRoot}}</string>{{end}}
    {{if .WorkingDirectory}}<key>WorkingDirectory</key>
    <string>{{html .WorkingDirectory}}</string>{{end}}
    <key>SessionCreate</key>
    <{{bool .SessionCreate}}/>
    <key>KeepAlive</key>
    <{{bool .KeepAlive}}/>
    <key>RunAtLoad</key>
    <{{bool .RunAtLoad}}/>
    <key>Disabled</key>
    <false/>

    <key>StandardOutPath</key>
    <string>/usr/local/var/log/{{html .Name}}.out.log</string>
    <key>StandardErrorPath</key>
    <string>/usr/local/var/log/{{html .Name}}.err.log</string>

  </dict>
</plist>
`
