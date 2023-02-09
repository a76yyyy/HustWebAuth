package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

type program struct {
	// cmd  *cobra.Command
	// args []string
}

func newSVCConfig() *service.Config {
	var logOutput = false
	if logFile != "" {
		logOutput = true
	}

	c := &service.Config{
		Name:        "HustWebAuth",
		DisplayName: "HustWebAuth",
		Description: "A service used to implement Ruijie web authentication.",
		Arguments:   []string{"service"},
		EnvVars:     map[string]string{"HOME": homeDir},
		Option:      service.KeyValue{"LogOutput": logOutput, "LogDirectory": logDir},
	}

	// Start only once network is up on Linux/systemd.
	if sysType == "linux" {
		c.Dependencies = []string{
			"After=syslog.target network.target",
		}
	}

	// Use different scripts on OpenWrt and FreeBSD.
	if IsOpenWrt() {
		c.Option["SysvScript"] = openWrtScript
	}

	return c
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
var (
	serviceCmd = &cobra.Command{
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

	installCmd = &cobra.Command{
		Use:   "install",
		Short: "Install HustWebAuth service",
		Run: func(cmd *cobra.Command, args []string) {
			svcConfig := newSVCConfig()

			s, err := newSVC(&program{}, svcConfig)
			if err != nil {
				log.Fatal(err)
				return
			}

			err = svcAction(s, "install")
			if err != nil {
				log.Fatal(err)
				return
			}
			if IsOpenWrt() {
				// On OpenWrt it is important to run enable after the service
				// installation.  Otherwise, the service won't start on the system
				// startup.
				_, err = runInitdCommand(s.String(), "enable")
				if err != nil {
					log.Fatalf("service: running init enable: %s", err)
				}
			}
			log.Println("HustWebAuth service has been installed")

			err = svcAction(s, "start")
			if err != nil {
				log.Fatal(err)
				return
			}
			log.Println("HustWebAuth service started.")
			saveCfg = true
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

			err = svcAction(s, "start")
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

			status, err := svcStatus(s)
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
			err = svcAction(s, "stop")
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
			err = svcAction(s, "restart")
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

			if IsOpenWrt() {
				// On OpenWrt it is important to run disable command first
				// as it will remove the symlink
				_, err := runInitdCommand(s.String(), "disable")
				if err != nil {
					log.Fatalf("service: running init disable: %s", err)
				}
			}

			status, err := svcStatus(s)
			if err != nil {
				log.Fatal(err)
				return
			}
			if status == service.StatusRunning {
				err = svcAction(s, "stop")
				if err != nil {
					log.Println(err)
				}
			}

			err = svcAction(s, "uninstall")
			if err != nil {
				log.Fatal(err)
				return
			}
			log.Println("HustWebAuth service has been uninstalled")
		},
	}
)

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(installCmd, startCmd, statusCmd, stopCmd, restartCmd, uninstallCmd)
}

// runInitdCommand runs init.d service command
// returns command code or error if any
func runInitdCommand(serviceName, action string) (int, error) {
	confPath := "/etc/init.d/" + serviceName
	// Pass the script and action as a single string argument.
	code, _, err := RunCommand("sh", "-c", confPath+" "+action)

	return code, err
}

// svcAction performs the action on the service.
//
// On OpenWrt, the service utility may not exist.  We use our service script
// directly in this case.
func svcAction(s service.Service, action string) (err error) {
	if sysType == "darwin" && action == "start" {
		var exe string
		if exe, err = os.Executable(); err != nil {
			log.Println("Starting service error: getting executable path: ", err)
		} else if exe, err = filepath.EvalSymlinks(exe); err != nil {
			log.Println("Starting service error: evaluating executable symlinks: ", err)
		} else if !strings.HasPrefix(exe, "/Applications/") {
			log.Println("warning: service must be started from within the /Applications directory")
		}
	}

	err = service.Control(s, action)
	if err != nil && service.Platform() == "unix-systemv" &&
		(action == "start" || action == "stop" || action == "restart") {
		_, err = runInitdCommand(s.String(), action)

		return err
	}

	return err
}

// svcStatus returns the service's status.
//
// On OpenWrt, the service utility may not exist.  We use our service script
// directly in this case.
func svcStatus(s service.Service) (status service.Status, err error) {
	status, err = s.Status()
	if err != nil && service.Platform() == "unix-systemv" {
		var code int
		code, err = runInitdCommand(s.String(), "status")
		if err != nil || code != 0 {
			return service.StatusStopped, nil
		}

		return service.StatusRunning, nil
	}

	return status, err
}

// OpenWrt procd init script
// https://github.com/AdguardTeam/AdGuardHome/issues/1386
const openWrtScript = `#!/bin/sh /etc/rc.common

START=90
STOP=01

cmd="{{.Path}}{{range .Arguments}} {{.|cmd}}{{end}}"
name="{{.Name}}"
pid_file="/var/run/${name}.pid"
stdout_log="{{.LogDirectory}}/$name.log"
stderr_log="{{.LogDirectory}}/$name.err"

{{range $k, $v := .EnvVars -}}
export {{$k}}={{$v}}
{{end -}}

EXTRA_COMMANDS="status"
EXTRA_HELP="$(printf "\t%-16s%s\n" "status" "Print the service status")"

get_pid() {
    cat "${pid_file}"
}

is_running() {
    [ -f "${pid_file}" ] && ps | grep -v grep | grep $(get_pid) >/dev/null 2>&1
}

start() {
    if is_running; then
        echo "Already started"
    else
        echo "Starting $name"
        {{if .WorkingDirectory}}cd '{{.WorkingDirectory}}'{{end}}
        mkdir -p {{.LogDirectory}}
        $cmd >> "$stdout_log" 2>> "$stderr_log" &
        echo $! > "$pid_file"
        if ! is_running; then
            echo "Unable to start, see $stdout_log and $stderr_log"
            exit 1
        fi
    fi
}

stop() {
    if is_running; then
        echo -n "Stopping $name.."
        kill $(get_pid)
        for i in $(seq 1 10)
        do
            if ! is_running; then
                break
            fi
            echo -n "."
            sleep 1
        done
        echo
        if is_running; then
            echo "Not stopped; may still be shutting down or shutdown may have failed"
            exit 1
        else
            echo "Stopped"
            if [ -f "$pid_file" ]; then
                rm "$pid_file"
            fi
        fi
    else
        echo "Not running"
    fi
}

restart() {
    stop
    if is_running; then
        echo "Unable to stop, will not attempt to start"
        exit 1
    fi
    start
}

status() {
    if is_running; then
        echo "Running"
    else
        echo "Stopped"
        exit 1
    fi
}

`
