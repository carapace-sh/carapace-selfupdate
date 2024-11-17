package transport

import "io"

type Transport interface {
	Tags(repo string, out, progress io.Writer) error
	Assets(repo, tag string, out, progress io.Writer) error
	Download(repo, tag, asset string, out, progress io.Writer) error
}
