package transport

import (
	"fmt"
	"io"
	"os/exec"
)

type Gh struct{}

func (t *Gh) retrieve(args []string, out, progress io.Writer) error {
	command := exec.Command("gh", args...)
	command.Stdout = out
	command.Stderr = progress
	return command.Run()
}

func (t *Gh) Tags(repo string, out, progress io.Writer) error {
	args := []string{"api", fmt.Sprintf("https://api.github.com/repos/%v/tags", repo)}
	return t.retrieve(args, out, progress)
}

func (t *Gh) Assets(repo, tag string, out, progress io.Writer) error {
	args := []string{"api", fmt.Sprintf("repos/%v/releases/tags/%v", repo, tag)}
	return t.retrieve(args, out, progress)
}

func (t *Gh) Download(repo, tag, asset string, out, progress io.Writer) error {
	args := []string{"release", "download", "--repo", repo, tag, "--pattern", asset, "--output", "-"}
	return t.retrieve(args, out, progress)
}
