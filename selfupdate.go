package selfupdate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/carapace-sh/carapace-selfupdate/filter"
	"github.com/carapace-sh/carapace-selfupdate/transport"
)

type config struct {
	repo   string
	binary string
	filter func(asset string) bool
	t      transport.Transport
}

func New(owner, repository string, opts ...func(c *config)) *config {
	c := &config{
		repo:   fmt.Sprintf("%v/%v", owner, repository),
		binary: repository,
		filter: filter.Goreleaser(repository),
		t:      &transport.Curl{},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithBinary(s string) func(c *config) {
	return func(c *config) {
		c.binary = s
	}
}

func WithAssetFilter(f func(s string) bool) func(c *config) {
	return func(c *config) {
		c.filter = f
	}
}

func WithTransport(t transport.Transport) func(c *config) {
	return func(c *config) {
		c.t = t
	}
}

func (c config) Assets(tag string) ([]string, error) {
	var b bytes.Buffer
	if err := c.t.Assets(c.repo, tag, &b); err != nil {
		return nil, err
	}

	var response struct {
		Name   string
		Assets []struct {
			Name string
		}
	}
	if err := json.Unmarshal(b.Bytes(), &response); err != nil {
		return nil, err
	}

	names := make([]string, 0, len(response.Assets))
	for _, asset := range response.Assets {
		if c.filter == nil || c.filter(asset.Name) {
			names = append(names, asset.Name)
		}
	}
	return names, nil
}

func (c config) Tags() ([]string, error) {
	var b bytes.Buffer
	if err := c.t.Tags(c.repo, &b); err != nil {
		return nil, err
	}

	var tags []struct {
		Name string
	}
	if err := json.Unmarshal(b.Bytes(), &tags); err != nil {
		return nil, err
	}

	names := make([]string, 0, len(tags))
	for _, tag := range tags {
		names = append(names, tag.Name)
	}
	return names, nil
}

func (c config) Download(tag, asset string) error {
	tmpfile, err := os.CreateTemp(os.TempDir(), "carapace-selfupdate_")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	f, err := os.Open(tmpfile.Name())
	if err != nil {
		return err
	}
	defer f.Close()

	sum, err := c.Checksum(tag, asset)
	if err != nil {
		return err
	}
	println("checksum:" + sum)

	return c.t.Download(c.repo, tag, asset, f)
}

func (c config) Checksum(tag, asset string) (string, error) {
	r := regexp.MustCompile(`^(?P<prefix>[^_]+_[^_]+)_.*$`)
	matches := r.FindStringSubmatch(asset)
	if matches == nil {
		return "", errors.New(`asset does not match checksum pattern`)
	}

	b := &bytes.Buffer{}
	if err := c.t.Download(c.repo, tag, fmt.Sprintf("%v_checksums.txt", matches[1]), b); err != nil {
		return "", err
	}

	m := make(map[string]string)
	for _, line := range strings.Split(b.String(), "\n") {
		if sum, file, ok := strings.Cut(line, "  "); ok {
			m[file] = sum
		}
	}
	return m[asset], nil
}
