//go:build !windows && !plan9

package cmd

import (
	"io"
	"io/fs"
	"log"
	"log/syslog"
	"os"
	"path/filepath"
)

func initLog() {
	logWriter := os.Stderr
	if logFile != "" {
		var err error
		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			os.Mkdir(logDir, fs.ModeDir)
		}
		if logRandom {
			logWriter, err = os.CreateTemp(logDir, logFile)
		} else if logAppend {
			logWriter, err = os.OpenFile(filepath.Join(logDir, logFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		} else {
			logWriter, err = os.OpenFile(filepath.Join(logDir, logFile), os.O_CREATE|os.O_WRONLY, 0644)
		}
		if err != nil {
			log.Fatal("Open log file failed, err:", err)
		}
		log.Println("Log file:", logWriter.Name())
	}

	if sysType != "windows" && sysLog {
		var err error
		sysLogWriter, err := syslog.New(syslog.LOG_INFO, "ruijie_web_login")
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(io.MultiWriter(logWriter, sysLogWriter))
	} else {
		log.SetOutput(logWriter)
	}
}
