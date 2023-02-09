//go:build darwin || freebsd || linux || openbsd

package cmd

import (
	"io/fs"
	"os"
)

func rootDirFS() (fsys fs.FS) {
	return os.DirFS("/")
}
