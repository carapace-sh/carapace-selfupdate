package selfupdate

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func (c config) Checksum(tag, asset string) (string, error) {
	checksumAsset, err := c.checksumAsset(tag)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	if err := c.Download(tag, checksumAsset, &b); err != nil {
		return "", err
	}

	return parseChecksum(b.String(), asset)
}

func (c config) checksumAsset(tag string) (string, error) {
	unfiltered := c
	unfiltered.filter = nil

	assets, err := unfiltered.Assets(tag)
	if err != nil {
		return "", err
	}

	for _, asset := range assets {
		if strings.HasSuffix(asset, "_checksums.txt") || asset == "checksums.txt" {
			return asset, nil
		}
	}
	for _, asset := range assets {
		if strings.HasSuffix(asset, "checksums.txt") {
			return asset, nil
		}
	}

	return "", fmt.Errorf("checksum asset not found for tag %q", tag)
}

func parseChecksum(content, asset string) (string, error) {
	for _, line := range strings.Split(content, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		sum := strings.ToLower(fields[0])
		if fields[len(fields)-1] == asset && isSHA256(sum) {
			return sum, nil
		}
	}

	return "", fmt.Errorf("checksum for %q not found", asset)
}

func isSHA256(s string) bool {
	if len(s) != sha256.Size*2 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}

func fileChecksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (c config) verifyChecksum(tag, asset, path string) error {
	expected, err := c.Checksum(tag, asset)
	if err != nil {
		return err
	}

	actual, err := fileChecksum(path)
	if err != nil {
		return err
	}

	if actual != expected {
		return fmt.Errorf("checksum mismatch for %q: expected %s, got %s", asset, expected, actual)
	}
	return nil
}
