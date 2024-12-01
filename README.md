# carapace-selfupdate

Simple self-update mechanism to install nightly/stable releases.

[![asciicast](https://asciinema.org/a/692857.svg)](https://asciinema.org/a/692857)

Executable is installed to the [GOBIN] directory, essentially shadowing any system installation.


```sh
export PATH="$HOME/.local/bin:$HOME/go/bin:$PATH"
#            │                │            └system installation (e.g. /usr/bin/carapace)
#            │                └selfupdate/go based installation ($GOBIN)
#            └user binaries
```

## Requirements

- [curl] for downloads
- [PATH] containing the [GOBIN] directory

[curl]:https://curl.se
[GOBIN]:https://pkg.go.dev/cmd/go#hdr-Compile_and_install_packages_and_dependencies
[PATH]:https://en.wikipedia.org/wiki/PATH_(variable)
