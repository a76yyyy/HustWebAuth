//go:build !windows && !plan9

package cmd

import (
	"io"
	"log"
	"log/syslog"
	"os"
)

func initLog() {
	logWriter := os.Stderr
	if logFile != "" {
		var err error
		logWriter, err = os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("Open log file failed, err:", err)
		}
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
