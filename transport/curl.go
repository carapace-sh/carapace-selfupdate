package transport

import (
	"fmt"
	"io"
	"os/exec"
)

type Curl struct{}

func (c Curl) retrieve(url string, out, progress io.Writer) error {
	command := exec.Command("curl", "-L", url)
	command.Stdout = out
	command.Stderr = progress
	return command.Run()
}

func (c Curl) Tags(repo string, out, progress io.Writer) error {
	url := fmt.Sprintf("https://api.github.com/repos/%v/tags", repo)
	return c.retrieve(url, out, progress)
}

func (c Curl) Assets(repo, tag string, out, progress io.Writer) error {
	url := fmt.Sprintf("https://api.github.com/repos/%v/releases/tags/%v", repo, tag)
	return c.retrieve(url, out, progress)
}

func (c Curl) Download(repo, tag, asset string, out, progress io.Writer) error {
	url := fmt.Sprintf("https://github.com/%v/releases/download/%v/%v", repo, tag, asset)
	return c.retrieve(url, out, progress)
}
