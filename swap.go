//go:build !windows

package selfupdate

import (
	"os"
)

// swap replaces the running executable.
// On most systems this simply renames the source file.
func (c config) swap(source, target string) error {
	c.Printf("moving to %#v\n", target)
	return os.Rename(source, target)
}
