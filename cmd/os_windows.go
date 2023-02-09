//go:build windows

package cmd

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

func rootDirFS() (fsys fs.FS) {
	// TODO(a.garipov): Use a better way if golang/go#44279 is ever resolved.
	sysDir, err := windows.GetSystemDirectory()
	if err != nil {
		log.Println("Error: Getting root filesystem: %s; using C:", err)

		// Assume that C: is the safe default.
		return os.DirFS("C:")
	}

	return os.DirFS(filepath.VolumeName(sysDir))
}

func isOpenWrt() (ok bool) {
	return false
}
