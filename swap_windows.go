package selfupdate

import (
	"os"
	"syscall"
)

// swap replaces the running executable.
// On windows it cannot be overwritten/deleted.
// It is thus renamed to an intermediate file (`{target}.old`) and simply hidden.
//
// see https://github.com/minio/selfupdate
// see https://github.com/sanbornm/go-selfupdate
func (c config) swap(source, target string) error {
	old := target + ".old"
	if _, err := os.Stat(old); err == nil {
		if err := c.confirm("remove %#v", old); err != nil {
			return err
		}
		if err := os.Remove(old); err != nil {
			return err
		}
	}

	if _, err := os.Stat(target); err == nil {
		c.Printf("moving current executable to %#v\n", old)
		if err := os.Rename(target, old); err != nil {
			return err
		}
	}

	c.Printf("moving new executable to %#v\n", target)
	if err := os.Rename(source, target); err != nil {
		_ = os.Rename(old, target) // try to restore old executable
		return err
	}

	if _, err := os.Stat(old); err == nil {
		c.Printf("hiding %#v\n", old)
		return hide(old)
	}
	return nil
}

func hide(path string) error {
	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	err = syscall.SetFileAttributes(p, syscall.FILE_ATTRIBUTE_HIDDEN)
	if err != nil {
		return err
	}

	return nil
}
