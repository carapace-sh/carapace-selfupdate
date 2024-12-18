package selfupdate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-selfupdate/filter"
	"github.com/carapace-sh/carapace-selfupdate/transport"
	"github.com/carapace-sh/carapace/pkg/traverse"
)

type config struct {
	repo     string
	binary   string
	filter   func(asset string) bool
	force    bool
	progress io.Writer
	t        transport.Transport
}

type option func(c *config)

func New(owner, repository string, opts ...option) *config {
	c := &config{
		repo:   fmt.Sprintf("%v/%v", owner, repository),
		binary: repository,
		filter: filter.Goreleaser(repository),
		t:      &transport.Curl{},
	}
	for _, opt := range opts {
		opt(c)
	}
	if runtime.GOOS == "windows" {
		c.binary += ".exe"
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

func WithForce(b bool) func(c *config) {
	return func(c *config) {
		c.force = b
	}
}

func WithTransport(t transport.Transport) func(c *config) {
	return func(c *config) {
		c.t = t
	}
}

func WithProgress(w io.Writer) func(c *config) {
	return func(c *config) {
		c.progress = w
	}
}

func (c config) Assets(tag string) ([]string, error) {
	var b bytes.Buffer
	if err := c.t.Assets(c.repo, tag, &b, c.progress); err != nil {
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
	if err := c.t.Tags(c.repo, &b, c.progress); err != nil {
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

func (c config) Println(s string) {
	c.Printf(s + "\n")
}

func (c config) Printf(format string, any ...any) {
	fmt.Fprintf(c.progress, "\x1b[1;2m"+format+"\x1b[0m", any...)
}

func (c config) confirm(format string, any ...any) error {
	if c.force {
		return nil
	}
	fmt.Fprintf(os.Stderr, "\x1b[1;2m"+format+" [y/n]: \x1b[0m", any...) // stderr should be configurable
	var input string
	if _, err := fmt.Scanln(&input); err != nil {
		return err
	}
	if strings.ToLower(input) != "y" {
		return errors.New("aborted")
	}
	return nil
}

func (c config) Install(tag, asset string) error {
	if !strings.HasSuffix(asset, ".tar.gz") && !strings.HasSuffix(asset, ".zip") {
		return errors.New("invalid extension [expected: .tar.gz|.zip]") // fail early
	}

	ext := strings.Replace(filepath.Ext(asset), ".gz", ".tar.gz", 1)
	tmpArchive, err := os.CreateTemp(os.TempDir(), "carapace-selfupdate_*"+ext)
	if err != nil {
		return err
	}
	defer os.Remove(tmpArchive.Name())

	f, err := os.Create(tmpArchive.Name())
	if err != nil {
		return err
	}
	defer f.Close()

	c.Printf("downloading to %#v\n", tmpArchive.Name())
	if err := c.Download(tag, asset, f); err != nil {
		return err
	}

	// sum, err := c.Checksum(tag, asset)
	// if err != nil {
	// 	return err
	// }

	// TODO verify checksum

	binDir, err := traverse.GoBinDir(carapace.NewContext())
	if err != nil {
		return err
	}

	if _, err := os.Stat(binDir); errors.Is(err, os.ErrNotExist) {
		if err := c.confirm("create directory %#v?", binDir); err != nil {
			return err
		}

		c.Printf("creating directory %#v\n", binDir)
		if err := os.MkdirAll(binDir, os.ModePerm); err != nil {
			return err
		}
	}

	tmpBinary := filepath.Join(binDir, c.binary+".selfupdate")
	if _, err := os.Stat(tmpBinary); err == nil {
		if err := c.confirm("overwrite %#v?", tmpBinary); err != nil {
			return err
		}
	}

	fExecutable, err := os.Create(tmpBinary)
	if err != nil {
		return err
	}
	defer os.Remove(fExecutable.Name())
	defer f.Close()

	c.Printf("extracting to %#v\n", fExecutable.Name())
	if err := c.extract(tmpArchive.Name(), fExecutable); err != nil {
		return err
	}

	if err := os.Chmod(fExecutable.Name(), 0755); err != nil {
		return err
	}

	if err := fExecutable.Close(); err != nil {
		return err
	}

	c.Printf("executing %#v\n", fExecutable.Name()+" --version")
	if err := c.verify(fExecutable.Name()); err != nil {
		return err
	}

	target := filepath.Join(binDir, c.binary)

	if _, err := os.Stat(target); err == nil {
		if err := c.confirm("overwrite %#v?", target); err != nil {
			return err
		}
	}

	return c.swap(fExecutable.Name(), target)
}

func (c config) verify(executable string) error {
	command := exec.Command(executable, "--version")
	command.Stdout = c.progress
	command.Stderr = c.progress
	return command.Run()
}

func (c config) extract(source string, out io.Writer) error {
	switch {
	case strings.HasSuffix(source, ".tar.gz"):
		command := exec.Command("tar", "--to-stdout", "-xzvf", source, c.binary)
		command.Stdout = out
		command.Stderr = c.progress
		return command.Run()
	case strings.HasSuffix(source, ".zip"):
		command := exec.Command("tar", "--to-stdout", "-xvf", source, c.binary)
		command.Stdout = out
		command.Stderr = c.progress
		return command.Run()
	default:
		return errors.New("invalid extension [expected: .tar.gz|.zip]")
	}
}

func (c config) Download(tag, asset string, out io.Writer) error {
	return c.t.Download(c.repo, tag, asset, out, c.progress)
}

func (c config) Checksum(tag, asset string) (string, error) {
	r := regexp.MustCompile(`^(?P<prefix>[^_]+_[^_]+)_.*$`)
	matches := r.FindStringSubmatch(asset)
	if matches == nil {
		return "", errors.New(`asset does not match checksum pattern`)
	}

	b := &bytes.Buffer{}
	if err := c.t.Download(c.repo, tag, fmt.Sprintf("%v_checksums.txt", matches[1]), b, c.progress); err != nil {
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
