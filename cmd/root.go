/*
Copyright © 2022 a76yyyy q981331502@163.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	daemon "github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile       string
	account       string
	password      string
	service       string
	register      bool
	pingIP        string
	pingCount     int
	pingTimeout   time.Duration
	pingPrivilege bool
	redirectURL   string
	logFile       string
	sysLog        bool
	saveCfg       bool
	daemonEnable  bool
	daemonPidFile string
	daemonTimeout time.Duration
	daemonRetry   int
	connectLog    bool
)

var path, _ = os.Executable()
var _, filenameWithSuffix = filepath.Split(path)
var sysType = runtime.GOOS

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   filenameWithSuffix,
	Short: "A program used to implement Ruijie web authentication",
	Long:  `ruijie_web_login is a program used to implement Ruijie web authentication.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if sysType != "windows" && daemonEnable {
			if logFile == "" {
				logFile = "/var/log/" + filenameWithSuffix + ".log"
			}
			if daemonPidFile == "" {
				daemonPidFile = "/var/run/" + filenameWithSuffix + ".pid"
			}
			cntxt := &daemon.Context{
				PidFileName: "/var/run/" + filenameWithSuffix + ".pid",
				PidFilePerm: 0644,
				LogFileName: logFile,
				LogFilePerm: 0644,
			}

			// Reborn()返回 子进程为nil 父进程不为nil
			child, err := cntxt.Reborn()
			if err != nil {
				log.Fatal("Unable to run: ", err)
			}
			if child != nil {
				return
			}
			defer cntxt.Release()

			log.Print("- - - - - - - - - - - - - - -")
			log.Print("daemon started")

			eventsTick := time.NewTicker(daemonTimeout)
			defer eventsTick.Stop()
			retryCount := 0
			for range eventsTick.C {
				res, err := Login()
				if err != nil {
					if daemonRetry < 0 {
						log.Println("Login failed, Err: ", err)
						log.Println("Login retrying...")
					} else if retryCount < daemonRetry {
						retryCount++
						log.Println("Login failed, Err: ", err)
						log.Println("Login retry", strconv.Itoa(retryCount), "times after", daemonTimeout.String())
					} else {
						log.Println("Login failed, Err: ", err)
						log.Println("Exceed the maximum number of retries, daemon stopped")
						os.Exit(1)
					}
				} else {
					if res != "" {
						log.Println(res)
					}
					retryCount = 0
				}
			}
		}

		res, err := Login()
		if err != nil {
			log.Fatal(err)
		}
		if res != "" {
			log.Println(res)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLog)
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(saveConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "", "Config file (default is localDir/ruijie.yaml)")
	rootCmd.PersistentFlags().StringVarP(&account, "account", "n", "", "Account")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Password")
	rootCmd.PersistentFlags().StringVarP(&service, "service", "S", "internet", "Service, options: [internet, local]")
	rootCmd.PersistentFlags().BoolVarP(&register, "register", "r", false, "Register Mac address")
	rootCmd.PersistentFlags().StringVar(&pingIP, "pingIP", "202.114.0.131", "IP address to ping")
	rootCmd.PersistentFlags().IntVar(&pingCount, "pingCount", 3, "ping count")
	rootCmd.PersistentFlags().DurationVar(&pingTimeout, "pingTimeout", 3*time.Second, "Ping timeout")
	rootCmd.PersistentFlags().BoolVar(&pingPrivilege, "pingPrivilege", true, `Sets the type of ping pinger will send. 
false means pinger will send an "unprivileged" UDP ping. 
true means pinger will send a "privileged" raw ICMP ping. 
NOTE: setting to true requires that it be run with super-user privileges.
`)
	rootCmd.PersistentFlags().StringVar(&redirectURL, "redirectURL", "http://123.123.123.123", "Redirect URL")
	rootCmd.PersistentFlags().StringVarP(&logFile, "logFile", "l", "", "Log file address (default means output to os.stdout)")
	rootCmd.PersistentFlags().BoolVar(&sysLog, "syslog", false, "Enable syslog, not support windows")
	rootCmd.PersistentFlags().BoolVarP(&saveCfg, "save", "s", false, "Save config file")
	rootCmd.PersistentFlags().BoolVarP(&daemonEnable, "daemon", "d", false, "Enable daemon mode, not support windows")
	rootCmd.PersistentFlags().StringVar(&daemonPidFile, "daemonPidFile", "", "Daemon pid file")
	rootCmd.PersistentFlags().DurationVar(&daemonTimeout, "daemonTimeout", 1*time.Minute, "Daemon cycle time")
	rootCmd.PersistentFlags().IntVar(&daemonRetry, "daemonRetry", 3, "Daemon retry times, -1 means retry forever")
	rootCmd.PersistentFlags().BoolVar(&connectLog, "connectLog", false, "Enable connect log")

	rootCmd.MarkFlagRequired("account")
	rootCmd.MarkFlagRequired("password")

	viper.BindPFlag("ping.IP", rootCmd.PersistentFlags().Lookup("pingIP"))
	viper.BindPFlag("ping.Count", rootCmd.PersistentFlags().Lookup("pingCount"))
	viper.BindPFlag("ping.Timeout", rootCmd.PersistentFlags().Lookup("pingTimeout"))
	viper.BindPFlag("ping.Privilege", rootCmd.PersistentFlags().Lookup("pingPrivilege"))
	viper.BindPFlag("redirect.URL", rootCmd.PersistentFlags().Lookup("redirectURL"))
	viper.BindPFlag("account", rootCmd.PersistentFlags().Lookup("account"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("service", rootCmd.PersistentFlags().Lookup("service"))
	// viper.BindPFlag("register", rootCmd.PersistentFlags().Lookup("register"))
	viper.BindPFlag("syslog", rootCmd.PersistentFlags().Lookup("syslog"))
	// viper.BindPFlag("save", rootCmd.PersistentFlags().Lookup("save"))
	viper.BindPFlag("daemon.Enable", rootCmd.PersistentFlags().Lookup("daemon"))
	viper.BindPFlag("daemon.PidFile", rootCmd.PersistentFlags().Lookup("daemonPidFile"))
	viper.BindPFlag("daemon.Timeout", rootCmd.PersistentFlags().Lookup("daemonTimeout"))
	viper.BindPFlag("daemon.Retry", rootCmd.PersistentFlags().Lookup("daemonRetry"))
	viper.BindPFlag("connectLog", rootCmd.PersistentFlags().Lookup("connectLog"))

	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		ex, err := os.Executable()
		if err != nil {
			log.Panic(err)
		}
		exDir := filepath.Dir(ex)
		cfgFile = filepath.Join(exDir, "ruijie.yaml")
		// Search config in home directory with name "ruijie" (without extension).
		viper.AddConfigPath(exDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("ruijie")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file: " + viper.ConfigFileUsed())
	}
}

func saveConfig() {
	if saveCfg {
		err := viper.WriteConfigAs(cfgFile)
		if err != nil {
			log.Fatal(err)
		}
	}
}
