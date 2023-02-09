package cmd

import (
	"log"

	"github.com/kardianos/service"
)

func (p *program) Start(service.Service) error {
	// Start should not block. Do the actual work async.
	log.Println("Starting HustWebAuth service...")
	go p.run()
	return nil
}

func (p *program) run() {
	runCycle()
}

func (p *program) Stop(service.Service) error {
	log.Println("Stoping HustWebAuth service...")
	return nil
}
