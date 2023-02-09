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
