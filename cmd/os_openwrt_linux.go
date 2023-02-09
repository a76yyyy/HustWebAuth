//go:build linux

package cmd

import (
	"io"

	"github.com/AdguardTeam/golibs/stringutil"
)

func isOpenWrt() (ok bool) {
	const etcReleasePattern = "etc/*release*"

	var err error
	ok, err = FileWalker(func(r io.Reader) (_ []string, cont bool, err error) {
		const osNameData = "openwrt"

		// This use of ReadAll is now safe, because FileWalker's Walk()
		// have limited r.
		var data []byte
		data, err = io.ReadAll(r)
		if err != nil {
			return nil, false, err
		}

		return nil, !stringutil.ContainsFold(string(data), osNameData), nil
	}).Walk(RootDirFS(), etcReleasePattern)

	return err == nil && ok
}
