package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os/exec"

	"github.com/AdguardTeam/golibs/mathutil"
)

// MaxCmdOutputSize is the maximum length of performed shell command output in
// bytes.
const MaxCmdOutputSize = 64 * 1024

// RunCommand runs shell command.
func RunCommand(command string, arguments ...string) (code int, output []byte, err error) {
	cmd := exec.Command(command, arguments...)
	out, err := cmd.Output()

	out = out[:mathutil.Min(len(out), MaxCmdOutputSize)]

	if err != nil {
		if eerr := new(exec.ExitError); errors.As(err, &eerr) {
			return eerr.ExitCode(), eerr.Stderr, nil
		}

		return 1, nil, fmt.Errorf("command %q failed: %w: %s", command, err, out)
	}

	return cmd.ProcessState.ExitCode(), out, nil
}

// IsOpenWrt returns true if host OS is OpenWrt.
func IsOpenWrt() (ok bool) {
	return isOpenWrt()
}

// RootDirFS returns the [fs.FS] rooted at the operating system's root.  On
// Windows it returns the fs.FS rooted at the volume of the system directory
// (usually, C:).
func RootDirFS() (fsys fs.FS) {
	return rootDirFS()
}
