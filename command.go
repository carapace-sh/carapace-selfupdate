package selfupdate

import (
	"fmt"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/style"
	"github.com/spf13/cobra"
)

func Command(owner, repository string, opts ...func(c *config)) *cobra.Command {
	repo := map[string]string{
		"stable":  repository,
		"nightly": "nightly",
	}

	cmd := &cobra.Command{
		Use:  "selfupdate [source] [tag]",
		Args: cobra.MinimumNArgs(1), // TODO
		RunE: func(cmd *cobra.Command, args []string) error {
			println("arr")
			if len(args) > 2 { // TODO test
				c := New(owner, repo[args[0]], opts...)
				return c.Download(args[1], args[2])
			}
			return nil
		},
	}

	cmd.Flags().BoolP("all", "a", false, "show all tags/assets")          // TODO disable filters
	cmd.Flags().BoolP("help", "h", false, "help for selfupdate")          // TODO use cobras help flag
	cmd.Flags().Bool("no-verify", false, "disable checksum verification") // TODO disable verification

	carapace.Gen(cmd).Standalone()

	carapace.Gen(cmd).PositionalCompletion(
		carapace.ActionStyledValuesDescribed(
			"stable", fmt.Sprintf("https://github.com/%v/%v", owner, repository), style.Green,
			"nightly", fmt.Sprintf("https://github.com/%v/nightly", owner), style.Red,
		),
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			tags, err := New(owner, repo[c.Args[0]], opts...).Tags()
			if err != nil {
				return carapace.ActionMessage(err.Error())
			}
			return carapace.ActionValues(tags...)
		}),
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			assets, err := New(owner, repo[c.Args[0]], opts...).Assets(c.Args[1])
			if err != nil {
				return carapace.ActionMessage(err.Error())
			}
			return carapace.ActionValues(assets...).StyleF(style.ForPathExt)
		}),
	)

	return cmd
}
