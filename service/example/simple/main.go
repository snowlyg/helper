package main

import (
	"fmt"

	"github.com/snowlyg/helper/service"
)

func (p *program) Start() error {
	// do some work
	return nil
}

func (p *program) Stop() error {
	//stop
	return nil
}

type program struct{}

func main() {
	// new windows service
	s, err := service.NewService(&program{}, "service-name")
	if err != nil {
		fmt.Printf("new service get error %v \n", err)
	}
	s.Run()
}
