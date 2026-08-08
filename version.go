package selfupdate

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

type version struct {
	major int
	minor int
	patch int
}

func (c config) newerTags(tags []string) []string {
	current, err := c.currentVersion()
	if err != nil {
		return tags
	}
	return filterNewerTags(tags, current)
}

func (c config) currentVersion() (string, error) {
	b, err := exec.Command(c.binary, "--version").Output()
	if err != nil {
		return "", err
	}

	for _, field := range strings.Fields(string(b)) {
		if _, ok := parseVersion(field); ok {
			return field, nil
		}
	}
	return "", errors.New("current version not found")
}

func filterNewerTags(tags []string, current string) []string {
	currentVersion, ok := parseVersion(current)
	if !ok {
		return tags
	}

	filtered := make([]string, 0, len(tags))
	for _, tag := range tags {
		tagVersion, ok := parseVersion(tag)
		if !ok {
			continue
		}
		if compareVersion(tagVersion, currentVersion) > 0 {
			filtered = append(filtered, tag)
		}
	}
	return filtered
}

func parseVersion(s string) (version, bool) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "v")
	if i := strings.IndexAny(s, "-+"); i != -1 {
		s = s[:i]
	}

	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return version{}, false
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return version{}, false
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return version{}, false
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return version{}, false
	}

	return version{major: major, minor: minor, patch: patch}, true
}

func compareVersion(a, b version) int {
	switch {
	case a.major != b.major:
		return a.major - b.major
	case a.minor != b.minor:
		return a.minor - b.minor
	default:
		return a.patch - b.patch
	}
}
