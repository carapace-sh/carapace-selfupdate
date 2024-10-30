package filter

import (
	"fmt"
	"regexp"
	"runtime"
)

func Goreleaser(project string) func(s string) bool {
	goos := runtime.GOOS
	if goos == "android" {
		goos = "termux" // TODO check for termux path (currently only termux builds on android)
	}

	goarch := runtime.GOARCH
	if goos == "termux" {
		switch goarch {
		case "arm":
			goarch = "armv6" // TODO is this really needed?
		}
	}

	r := regexp.MustCompile(fmt.Sprintf(
		`^%v_(?P<version>.+)_%v_%v(\.zip|\.tar\.gz)$`,
		regexp.QuoteMeta(project),
		regexp.QuoteMeta(goos),
		regexp.QuoteMeta(goarch),
	))

	return func(s string) bool {
		return r.MatchString(s)
	}
}
