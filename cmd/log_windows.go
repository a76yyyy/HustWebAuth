//go:build windows || plan9

package cmd

import (
	"log"
	"os"
)

func initLog() {
	logWriter := os.Stderr
	if logFile != "" {
		var err error
		logWriter, err = os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("Open log file failed, Err:", err)
		}
	}
	log.SetOutput(logWriter)
}
