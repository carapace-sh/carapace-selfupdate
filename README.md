# carapace-selfupdate

Simple self-update mechanism to install nightly/stable releases.

[![asciicast](https://asciinema.org/a/692857.svg)](https://asciinema.org/a/692857)


- Relies on [curl] to avoid additional dependencies.
- Installs to [GOBIN] directory which essentially shadows any system installation.
- [PATH] needs to contain the [GOBIN] directory for this to work.
  ```sh
  export PATH="$HOME/.local/bin:$HOME/go/bin:$PATH"
  ```
  > Executables are installed in the directory named by the GOBIN environment variable, which defaults to $GOPATH/bin or $HOME/go/bin if the GOPATH environment variable is not set. Executables in $GOROOT are installed in $GOROOT/bin or $GOTOOLDIR instead of $GOBIN.

[curl]:https://curl.se
[GOBIN]:https://pkg.go.dev/cmd/go#hdr-Compile_and_install_packages_and_dependencies
[PATH]:https://en.wikipedia.org/wiki/PATH_(variable)
