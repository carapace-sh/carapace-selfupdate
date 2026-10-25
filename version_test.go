package selfupdate

import (
	"reflect"
	"testing"
)

func TestFilterNewerTags(t *testing.T) {
	tags := []string{"v1.7.0", "v1.6.6", "v1.6.5", "v1.6.4", "nightly"}
	got := filterNewerTags(tags, "1.6.5")
	want := []string{"v1.7.0", "v1.6.6"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestFilterNewerTagsKeepsAllWhenCurrentVersionIsUnknown(t *testing.T) {
	tags := []string{"v1.7.0", "v1.6.5"}
	got := filterNewerTags(tags, "develop")

	if !reflect.DeepEqual(got, tags) {
		t.Fatalf("expected %v, got %v", tags, got)
	}
}

func TestParseVersion(t *testing.T) {
	got, ok := parseVersion("v1.2.3")
	if !ok {
		t.Fatal("expected version to parse")
	}
	want := version{major: 1, minor: 2, patch: 3}
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestParseVersionWithPrerelease(t *testing.T) {
	got, ok := parseVersion("1.2.3-beta.1")
	if !ok {
		t.Fatal("expected version to parse")
	}
	want := version{major: 1, minor: 2, patch: 3}
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
