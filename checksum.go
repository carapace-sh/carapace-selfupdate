package selfupdate

import (
	"crypto/sha256"
	"io"
	"os"
)

//lint:ignore U1000 TODO
func verify(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	return nil
}
