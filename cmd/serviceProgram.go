package cmd

import (
	"log"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
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

var (
	installCmd = &cobra.Command{
		Use:   "install",
		Short: "Install HustWebAuth service",
		Run: func(cmd *cobra.Command, args []string) {
			saveCfg = true

			svcConfig := newSVCConfig()

			s, err := newSVC(&program{}, svcConfig)
			if err != nil {
				log.Fatal(err)
				return
			}

			err = s.Install()
			if err != nil {
				log.Fatal(err)
				return
			}
			log.Println("HustWebAuth service has been installed")
		},
	}

	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start HustWebAuth service",
		Run: func(cmd *cobra.Command, args []string) {
			s, err := newSVC(&program{}, newSVCConfig())
			if err != nil {
				log.Fatal(err)
				return
			}

			err = s.Start()
			if err != nil {
				log.Fatal(err)
				return
			}
			log.Println("HustWebAuth service started.")
		},
	}

	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Get HustWebAuth service status",
		Run: func(cmd *cobra.Command, args []string) {
			s, err := newSVC(&program{}, newSVCConfig())
			if err != nil {
				log.Fatal(err)
				return
			}

			status, err := s.Status()
			if err != nil {
				log.Fatal(err)
				return
			}
			switch status {
			case service.StatusUnknown:
				log.Println("HustWebAuth service status is unable to be determined due to an error or it was not installed.")
			case service.StatusStopped:
				log.Println("HustWebAuth service is stopped.")
			case service.StatusRunning:
				log.Println("HustWebAuth service is running.")
			}
		},
	}

	stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop HustWebAuth service",
		Run: func(cmd *cobra.Command, args []string) {

			s, err := newSVC(&program{}, newSVCConfig())
			if err != nil {
				log.Fatal(err)
				return
			}
			err = s.Stop()
			if err != nil {
				log.Fatal(err)
				return
			}
			log.Println("HustWebAuth service stoped.")
		},
	}

	restartCmd = &cobra.Command{
		Use:   "restart",
		Short: "Restart HustWebAuth service",
		Run: func(cmd *cobra.Command, args []string) {
			s, err := newSVC(&program{}, newSVCConfig())
			if err != nil {
				log.Fatal(err)
				return
			}
			err = s.Restart()
			if err != nil {
				log.Fatal(err)
				return
			}
			log.Println("HustWebAuth service has been restarted.")
		},
	}

	uninstallCmd = &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall HustWebAuth service from system",
		Run: func(cmd *cobra.Command, args []string) {
			s, err := newSVC(&program{}, newSVCConfig())
			if err != nil {
				log.Fatal(err)
				return
			}

			err = s.Uninstall()
			if err != nil {
				log.Fatal(err)
				return
			}
			log.Println("HustWebAuth service has been uninstalled")
		},
	}
)
