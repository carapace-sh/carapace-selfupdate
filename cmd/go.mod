module github.com/carapace-sh/carapace-selfupdate/cmd

go 1.23.1

require (
	github.com/carapace-sh/carapace v1.8.6
	github.com/carapace-sh/carapace-selfupdate v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.9.1
)

require (
	github.com/carapace-sh/carapace-shlex v1.0.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/carapace-sh/carapace-selfupdate => ../
