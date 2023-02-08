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

// A program used to implement Ruijie web authentication
package cmd

import (
	"io/fs"
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
	serviceType   string
	pingIP        string
	pingCount     int
	pingTimeout   time.Duration
	pingPrivilege bool
	redirectURL   string
	logDir        string
	logFile       string
	logRandom     bool
	logAppend     bool
	sysLog        bool
	saveCfg       bool
	daemonEnable  bool
	daemonPidFile string
	cycleEnable   bool
	cycleDuration time.Duration
	cycleRetry    int
	logConnected  bool
)

var path, _ = os.Executable()
var _, filenameWithSuffix = filepath.Split(path)
var sysType = runtime.GOOS

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   filenameWithSuffix,
	Short: "A program used to implement Ruijie web authentication",
	Long:  `HustWebAuth is a program used to implement Ruijie web authentication.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		runDaemon()
	},
}

func runDaemon() {
	if saveCfg {
		return
	}
	if sysType != "windows" && daemonEnable {
		if logFile == "" {
			tmpDir := filepath.Join(os.TempDir(), "HustWebAuth")
			if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
				os.Mkdir(tmpDir, fs.ModeDir)
			}
			logFile = filepath.Join(tmpDir, filenameWithSuffix+".log")
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

		log.Println("- - - - - - - - - - - - - - - - - - -")
		log.Println("HustWebAuth Daemon started.")
	}

	runCycle()
}

func runCycle() {
	log.Println("- - - - - - - - - - - - - - - - - - -")
	log.Println("HustWebAuth started.")
	retryCount := 0
	res, err := Login()
	if err != nil {
		if cycleEnable {
			if cycleRetry < 0 {
				log.Println("Login failed, Err: ", err)
				log.Println("Login retrying...")
			} else if retryCount < cycleRetry {
				retryCount++
				log.Println("Login failed, Err: ", err)
				log.Println("Login retry ", strconv.Itoa(retryCount), "times after "+cycleDuration.String())
			} else {
				log.Fatal("Login failed, Err: ", err)
			}
		} else {
			log.Fatal("Login failed, Err: ", err)
		}
	}
	if res != "" {
		log.Println(res)
	}

	if cycleEnable {
		eventsTick := time.NewTicker(cycleDuration)
		defer eventsTick.Stop()
		for range eventsTick.C {
			res, err := Login()
			if err != nil {
				if cycleRetry < 0 {
					log.Println("Login failed, Err: ", err)
					log.Println("Login retrying...")
				} else if retryCount < cycleRetry {
					retryCount++
					log.Println("Login failed, Err: ", err)
					log.Println("Login retry", strconv.Itoa(retryCount), "times after", cycleDuration.String())
				} else {
					log.Println("Login failed, Err: ", err)
					log.Fatal("Exceed the maximum number of retries, daemon stopped!")
				}
			} else {
				if res != "" {
					log.Println(res)
				}
				retryCount = 0
			}
		}
	}
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
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initLog)
	cobra.OnFinalize(saveConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "", "Config file (default is $HOME/HustWebAuth.yaml)")
	rootCmd.PersistentFlags().StringVarP(&account, "account", "a", "", "Account for ruijie web authentication")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Password for ruijie web authentication")
	rootCmd.PersistentFlags().StringVarP(&serviceType, "serviceType", "s", "internet", "Service type, options: [internet, local]")

	rootCmd.PersistentFlags().StringVar(&pingIP, "pingIP", "202.114.0.131", "IP address to ping")
	rootCmd.PersistentFlags().IntVar(&pingCount, "pingCount", 3, "ping count")
	rootCmd.PersistentFlags().DurationVar(&pingTimeout, "pingTimeout", 3*time.Second, "Ping timeout")
	rootCmd.PersistentFlags().BoolVar(&pingPrivilege, "pingPrivilege", true, `Sets the type of ping pinger will send. 
false means pinger will send an "unprivileged" UDP ping. 
true means pinger will send a "privileged" raw ICMP ping. 
NOTE: setting to true requires that it be run with super-user privileges.
`)
	rootCmd.PersistentFlags().StringVar(&redirectURL, "redirectURL", "http://123.123.123.123", "Redirect URL")
	rootCmd.PersistentFlags().StringVar(&logDir, "logDir", filepath.Join(os.TempDir(), "HustWebAuth"), "Log Directory")
	rootCmd.PersistentFlags().StringVarP(&logFile, "logFile", "l", "", "Log file name (default means output to os.stdout)")
	rootCmd.PersistentFlags().BoolVar(&logRandom, "logRandom", true, "Log file name with random string.\nNOTE: If logFile includes a \"*\", the random string replaces the last \"*\".\n")
	rootCmd.PersistentFlags().BoolVar(&logAppend, "logAppend", true, "Log file append mode. \nNOTE: if logRandom is true, it will be ignored")
	rootCmd.PersistentFlags().BoolVar(&logConnected, "logConnected", true, "Enable logging of \"The network is connected\"")
	rootCmd.PersistentFlags().BoolVar(&sysLog, "syslog", false, "Enable syslog, not support windows")
	rootCmd.PersistentFlags().BoolVarP(&saveCfg, "save", "o", false, "Save config file")
	rootCmd.Flags().BoolVarP(&daemonEnable, "daemon", "d", false, "Enable daemon mode, not support windows")
	rootCmd.Flags().StringVar(&daemonPidFile, "daemonPidFile", "", "Daemon pid file")
	rootCmd.Flags().BoolVarP(&cycleEnable, "cycle", "c", false, "Enable cycle mode")
	rootCmd.Flags().DurationVar(&cycleDuration, "cycleDuration", 5*time.Minute, "Cycle duration")
	rootCmd.Flags().IntVar(&cycleRetry, "cycleRetry", 3, "Cycle retry times, -1 means retry forever")

	rootCmd.MarkFlagRequired("account")
	rootCmd.MarkFlagRequired("password")

	viper.BindPFlag("auth.account", rootCmd.PersistentFlags().Lookup("account"))
	viper.BindPFlag("auth.password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("auth.serviceType", rootCmd.PersistentFlags().Lookup("serviceType"))
	viper.BindPFlag("ping.ip", rootCmd.PersistentFlags().Lookup("pingIP"))
	viper.BindPFlag("ping.count", rootCmd.PersistentFlags().Lookup("pingCount"))
	viper.BindPFlag("ping.timeout", rootCmd.PersistentFlags().Lookup("pingTimeout"))
	viper.BindPFlag("ping.privilege", rootCmd.PersistentFlags().Lookup("pingPrivilege"))
	viper.BindPFlag("redirect.url", rootCmd.PersistentFlags().Lookup("redirectURL"))
	viper.BindPFlag("log.dir", rootCmd.PersistentFlags().Lookup("logDir"))
	viper.BindPFlag("log.file", rootCmd.PersistentFlags().Lookup("logFile"))
	viper.BindPFlag("log.random", rootCmd.PersistentFlags().Lookup("logRandom"))
	viper.BindPFlag("log.append", rootCmd.PersistentFlags().Lookup("logAppend"))
	viper.BindPFlag("log.connected", rootCmd.PersistentFlags().Lookup("logConnected"))
	viper.BindPFlag("log.syslog", rootCmd.PersistentFlags().Lookup("syslog"))
	viper.BindPFlag("daemon.enable", rootCmd.Flags().Lookup("daemon"))
	viper.BindPFlag("daemon.pidFile", rootCmd.Flags().Lookup("daemonPidFile"))
	viper.BindPFlag("cycle.enable", rootCmd.Flags().Lookup("cycle"))
	viper.BindPFlag("cycle.duration", rootCmd.Flags().Lookup("cycleDuration"))
	viper.BindPFlag("cycle.retry", rootCmd.Flags().Lookup("cycleRetry"))

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
		// ex, err := os.Executable()
		// if err != nil {
		// 	log.Panic(err)
		// }
		// exDir := filepath.Dir(ex)
		// viper.AddConfigPath(exDir)

		home, err := os.UserHomeDir()
		if err != nil {
			log.Panic(err)
		}
		// Search config in home directory with name "HustWebAuth" (without extension).
		viper.AddConfigPath(home)
		cfgFile = filepath.Join(home, "HustWebAuth.yaml")

		viper.SetConfigType("yaml")
		viper.SetConfigName("HustWebAuth")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file: " + viper.ConfigFileUsed())
		account = viper.GetString("auth.account")
		password = viper.GetString("auth.password")
		serviceType = viper.GetString("auth.serviceType")
		pingIP = viper.GetString("ping.ip")
		pingCount = viper.GetInt("ping.count")
		pingTimeout = viper.GetDuration("ping.timeout")
		pingPrivilege = viper.GetBool("ping.privilege")
		redirectURL = viper.GetString("redirect.url")
		logDir = viper.GetString("log.dir")
		logFile = viper.GetString("log.file")
		logRandom = viper.GetBool("log.random")
		logAppend = viper.GetBool("log.append")
		logConnected = viper.GetBool("log.connected")
		sysLog = viper.GetBool("log.syslog")
		daemonEnable = viper.GetBool("daemon.enable")
		daemonPidFile = viper.GetString("daemon.pidFile")
		cycleEnable = viper.GetBool("cycle.enable")
		cycleDuration = viper.GetDuration("cycle.duration")
		cycleRetry = viper.GetInt("cycle.retry")
	}
}

func saveConfig() {
	if saveCfg {
		err := viper.WriteConfigAs(cfgFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Save config file: " + cfgFile)
	}
}
