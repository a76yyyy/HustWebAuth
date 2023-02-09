//go:build darwin || openbsd || freebsd

package cmd

func isOpenWrt() (ok bool) {
	return false
}
