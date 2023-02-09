package cmd

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// 获取当前执行文件绝对路径
func getCurrentAbPath() string {
	execPath := getCurrentAbPathByExecutable()
	if strings.Contains(execPath, getTmpDir()) {
		return getCurrentAbPathByCaller()
	}
	return execPath
}

// 获取当前执行文件绝对目录
func getCurrentAbDir() string {
	return filepath.Dir(getCurrentAbPath())
}

// 获取系统临时目录，兼容go run
func getTmpDir() string {
	dir := os.Getenv("TEMP")
	if dir == "" {
		dir = os.Getenv("TMP")
		if dir == "" {
			dir = os.TempDir()
		}
	}
	res, _ := filepath.EvalSymlinks(dir)
	return res
}

// 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(exePath)
	return res
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath, _ = filepath.Abs(filename)
	}
	return abPath
}
