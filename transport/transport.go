package transport

import "io"

type Transport interface {
	Tags(repo string, out io.Writer) error
	Assets(repo, tag string, out io.Writer) error
	Download(repo, tag, asset string, out io.Writer) error
}
