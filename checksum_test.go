package selfupdate

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

type checksumTransport struct {
	assets    string
	downloads map[string]string
}

func (t checksumTransport) Tags(_ string, _, _ io.Writer) error {
	return nil
}

func (t checksumTransport) Assets(_, _ string, out, _ io.Writer) error {
	_, err := io.WriteString(out, t.assets)
	return err
}

func (t checksumTransport) Download(_, _, asset string, out, _ io.Writer) error {
	_, err := io.WriteString(out, t.downloads[asset])
	return err
}

func TestVerifyChecksum(t *testing.T) {
	asset := "project_1.0.0_linux_amd64.tar.gz"
	content := "archive bytes"
	sum := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
	path := writeTempFile(t, content)

	c := New("owner", "project", WithTransport(checksumTransport{
		assets: fmt.Sprintf(`{"assets":[{"name":"project_1.0.0_checksums.txt"},{"name":%q}]}`, asset),
		downloads: map[string]string{
			"project_1.0.0_checksums.txt": fmt.Sprintf("%s  %s\n", sum, asset),
		},
	}))

	if err := c.verifyChecksum("v1.0.0", asset, path); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyChecksumMismatch(t *testing.T) {
	asset := "project_1.0.0_linux_amd64.tar.gz"
	path := writeTempFile(t, "archive bytes")

	c := New("owner", "project", WithTransport(checksumTransport{
		assets: fmt.Sprintf(`{"assets":[{"name":"project_1.0.0_checksums.txt"},{"name":%q}]}`, asset),
		downloads: map[string]string{
			"project_1.0.0_checksums.txt": strings.Repeat("0", sha256.Size*2) + "  " + asset + "\n",
		},
	}))

	err := c.verifyChecksum("v1.0.0", asset, path)
	if err == nil || !strings.Contains(err.Error(), "checksum mismatch") {
		t.Fatalf("expected checksum mismatch, got %v", err)
	}
}

func TestChecksumMissingAsset(t *testing.T) {
	c := New("owner", "project", WithTransport(checksumTransport{
		assets: `{"assets":[{"name":"project_1.0.0_linux_amd64.tar.gz"}]}`,
	}))

	_, err := c.Checksum("v1.0.0", "project_1.0.0_linux_amd64.tar.gz")
	if err == nil || !strings.Contains(err.Error(), "checksum asset not found") {
		t.Fatalf("expected missing checksum asset error, got %v", err)
	}
}

func TestChecksumMissingEntry(t *testing.T) {
	c := New("owner", "project", WithTransport(checksumTransport{
		assets: `{"assets":[{"name":"project_1.0.0_checksums.txt"}]}`,
		downloads: map[string]string{
			"project_1.0.0_checksums.txt": strings.Repeat("a", sha256.Size*2) + "  other.tar.gz\n",
		},
	}))

	_, err := c.Checksum("v1.0.0", "project_1.0.0_linux_amd64.tar.gz")
	if err == nil || !strings.Contains(err.Error(), "checksum for") {
		t.Fatalf("expected missing checksum entry error, got %v", err)
	}
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()

	f, err := os.CreateTemp(t.TempDir(), "archive-*")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.WriteString(f, content); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	return f.Name()
}
