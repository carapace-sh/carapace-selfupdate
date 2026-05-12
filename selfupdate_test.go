package selfupdate

import (
	"io"
	"strings"
	"testing"
)

type fakeTransport struct {
	tags   string
	assets string
}

func (f fakeTransport) Tags(_ string, out, _ io.Writer) error {
	_, err := io.WriteString(out, f.tags)
	return err
}

func (f fakeTransport) Assets(_, _ string, out, _ io.Writer) error {
	_, err := io.WriteString(out, f.assets)
	return err
}

func (f fakeTransport) Download(_, _, _ string, _, _ io.Writer) error {
	return nil
}

func TestResolveDefaultsTagAndAsset(t *testing.T) {
	c := New(
		"owner",
		"repository",
		WithTransport(fakeTransport{
			tags:   `[{"name":"v2.0.0"},{"name":"v1.0.0"}]`,
			assets: `{"assets":[{"name":"checksums.txt"},{"name":"matching-asset.tar.gz"}]}`,
		}),
		WithAssetFilter(func(s string) bool { return strings.HasSuffix(s, ".tar.gz") }),
		WithProgress(io.Discard),
	)

	tag, asset, err := c.resolve("", "")
	if err != nil {
		t.Fatal(err)
	}
	if tag != "v2.0.0" {
		t.Fatalf("expected latest tag, got %q", tag)
	}
	if asset != "matching-asset.tar.gz" {
		t.Fatalf("expected first matching asset, got %q", asset)
	}
}

func TestResolveKeepsProvidedTagAndAsset(t *testing.T) {
	c := New(
		"owner",
		"repository",
		WithTransport(fakeTransport{}),
		WithProgress(io.Discard),
	)

	tag, asset, err := c.resolve("v1.0.0", "repository_v1.0.0_linux_amd64.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	if tag != "v1.0.0" {
		t.Fatalf("expected provided tag, got %q", tag)
	}
	if asset != "repository_v1.0.0_linux_amd64.tar.gz" {
		t.Fatalf("expected provided asset, got %q", asset)
	}
}

func TestResolveErrorsWhenTagsAreEmpty(t *testing.T) {
	c := New(
		"owner",
		"repository",
		WithTransport(fakeTransport{tags: `[]`}),
		WithProgress(io.Discard),
	)

	if _, _, err := c.resolve("", ""); err == nil || err.Error() != "no tags found" {
		t.Fatalf("expected no tags error, got %v", err)
	}
}

func TestResolveErrorsWhenNoAssetsMatch(t *testing.T) {
	c := New(
		"owner",
		"repository",
		WithTransport(fakeTransport{
			tags:   `[{"name":"v1.0.0"}]`,
			assets: `{"assets":[{"name":"checksums.txt"}]}`,
		}),
		WithAssetFilter(func(string) bool { return false }),
		WithProgress(io.Discard),
	)

	_, _, err := c.resolve("", "")
	if err == nil || !strings.Contains(err.Error(), "no matching assets found") {
		t.Fatalf("expected no matching assets error, got %v", err)
	}
}
