package main

import (
	"fmt"

	"github.com/snowlyg/helper/service"
)

func (p *program) Start(s service.Service) error {
	// do some work
	return nil
}

func (p *program) Stop(s service.Service) error {
	//stop
	return nil
}

type program struct{}

func main() {
	// new windows service
	s, err := service.NewService(&program{}, &service.Config{Name: "service-name"})
	if err != nil {
		fmt.Printf("new service get error %v \n", err)
	}
	s.Run()
}
