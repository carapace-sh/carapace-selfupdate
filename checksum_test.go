package selfupdate

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"testing"
)

func TestVerifyChecksum(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "archive_*.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("archive contents"); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	sum := sha256.Sum256([]byte("archive contents"))
	if err := verifyChecksum(f.Name(), hex.EncodeToString(sum[:])); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyChecksumMismatch(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "archive_*.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("archive contents"); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	if err := verifyChecksum(f.Name(), "deadbeef"); err == nil {
		t.Fatal("expected checksum mismatch")
	}
}
