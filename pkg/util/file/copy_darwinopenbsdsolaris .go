//go:build darwin || openbsd || solaris
// +build darwin openbsd solaris

package file

import (
	"fmt"
)

func CopyDirectory(scrDir, dest string) error {
	return fmt.Errorf("not implemented")
}
