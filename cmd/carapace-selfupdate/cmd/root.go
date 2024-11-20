package cmd

import (
	"github.com/carapace-sh/carapace"
	selfupdate "github.com/carapace-sh/carapace-selfupdate"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "carapace-selfupdate",
	Short: "",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute(version string) error {
	rootCmd.Version = version
	return rootCmd.Execute()
}
func init() {
	carapace.Gen(rootCmd).Standalone()

	cmd := selfupdate.Command("carapace-sh", "carapace-bin", selfupdate.WithBinary("carapace"))
	cmd.Use = "carapace"
	rootCmd.AddCommand(cmd)

	cmd = selfupdate.Command("carapace-sh", "carapace-bridge")
	cmd.Use = "bridge"
	rootCmd.AddCommand(cmd)

	cmd = selfupdate.Command("carapace-sh", "carapace-shlex")
	cmd.Use = "shlex"
	rootCmd.AddCommand(cmd)

	cmd = selfupdate.Command("carapace-sh", "carapace-spec")
	cmd.Use = "spec"
	rootCmd.AddCommand(cmd)

}
