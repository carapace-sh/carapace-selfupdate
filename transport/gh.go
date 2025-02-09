package transport

import (
	"fmt"
	"io"
	"os/exec"
)

type Gh struct{}

func (c Gh) api(path string, out, progress io.Writer) error {
	command := exec.Command("gh", "api", path)
	command.Stdout = out
	command.Stderr = progress
	return command.Run()
}

func (c Gh) Tags(repo string, out, progress io.Writer) error {
	path := fmt.Sprintf("repos/%v/tags", repo)
	return c.api(path, out, progress)
}

func (c Gh) Assets(repo, tag string, out, progress io.Writer) error {
	path := fmt.Sprintf("repos/%v/releases/tags/%v", repo, tag)
	return c.api(path, out, progress)
}

func (c Gh) Download(repo, tag, asset string, out, progress io.Writer) error {
	command := exec.Command("gh", "release", "download", "--repo", repo, tag, "--pattern", asset, "--output", "-")
	command.Stdout = out
	command.Stderr = progress
	return command.Run()
}
