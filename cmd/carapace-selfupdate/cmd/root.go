package cmd

import (
	"fmt"
	"strings"

	"github.com/carapace-sh/carapace"
	selfupdate "github.com/carapace-sh/carapace-selfupdate"
	"github.com/carapace-sh/carapace-selfupdate/filter"
	"github.com/carapace-sh/carapace/pkg/traverse"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "carapace-selfupdate",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		su := selfupdate.New("carapace-sh", "carapace")
		releases, err := su.Tags()
		if err != nil {
			return err
		}
		fmt.Println(strings.Join(releases, "\n"))

		path, err := traverse.GoBinDir(carapace.NewContext())
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	}}

func Execute(version string) error {
	rootCmd.Version = version
	return rootCmd.Execute()
}
func init() {
	carapace.Gen(rootCmd).Standalone()

	cmd := selfupdate.Command("carapace-sh", "carapace-shlex")
	cmd.Use = "shlex"
	rootCmd.AddCommand(cmd)

	cmd = selfupdate.Command("carapace-sh", "carapace-bin", selfupdate.WithBinary("carapace"), selfupdate.WithAssetFilter(filter.Goreleaser("carapace-bin")))
	cmd.Use = "carapace"
	rootCmd.AddCommand(cmd)
}
