package cmd

import (
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

type program struct {
	// cmd  *cobra.Command
	// args []string
}

func newSVCConfig() *service.Config {
	return &service.Config{
		Name:        "HustWebAuth",
		DisplayName: "HustWebAuth",
		Description: "A service used to implement Ruijie web authentication.",
		Arguments:   []string{"service"},
	}
}

func newSVC(prg *program, conf *service.Config) (service.Service, error) {
	s, err := service.New(prg, conf)
	if err != nil {
		// log.Fatal(err)
		return nil, err
	}
	return s, nil
}

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "System service related commands",
	Long:  `Use HustWebAuth as a system service: install, start, stop, uninstall, etc.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := newSVC(&program{}, newSVCConfig())
		if err != nil {
			return err
		}
		return s.Run()
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(installCmd, startCmd, stopCmd, restartCmd, uninstallCmd)
}
