package selfupdate

import (
	"fmt"
	"io"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-selfupdate/filter"
	"github.com/carapace-sh/carapace/pkg/style"
	"github.com/spf13/cobra"
)

func Command(owner, repository string, opts ...option) *cobra.Command {
	opts = append([]option{
		WithBinary(repository),
		WithAssetFilter(filter.Goreleaser(repository)),
	}, opts...)
	repo := map[string]string{
		"stable":  repository,
		"nightly": "nightly",
	}

	cmd := &cobra.Command{
		Use:   "selfupdate [source] [tag] [asset]",
		Short: "install nightly/stable releases",
		Args:  cobra.MinimumNArgs(3), // TODO implicit update to latest
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 2 { // TODO test
				opts := append([]option{WithProgress(cmd.ErrOrStderr())}, opts...)
				if cmd.Flag("force").Changed {
					opts = append(opts, WithForce(true))
				}
				c := New(owner, repo[args[0]], opts...)
				return c.Install(args[1], args[2])
			}
			return nil
		},
	}

	cmd.Flags().BoolP("all", "a", false, "show all tags/assets")
	cmd.Flags().BoolP("force", "f", false, "force")
	cmd.Flags().BoolP("help", "h", false, "help for selfupdate") // TODO use cobras help flag
	// cmd.Flags().Bool("no-verify", false, "disable checksum verification") // TODO disable verification

	carapace.Gen(cmd).Standalone()

	carapace.Gen(cmd).PositionalCompletion(
		carapace.ActionStyledValuesDescribed(
			"stable", fmt.Sprintf("https://github.com/%v/%v", owner, repository), style.Green,
			"nightly", fmt.Sprintf("https://github.com/%v/nightly", owner), style.Red,
		),
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			opts = append(opts, WithProgress(io.Discard))
			if cmd.Flag("all").Changed {
				opts = append(opts, WithAssetFilter(nil))
			}
			tags, err := New(owner, repo[c.Args[0]], opts...).Tags()
			if err != nil {
				return carapace.ActionMessage(err.Error())
			}
			return carapace.ActionValues(tags...)
		}),
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			opts = append(opts, WithProgress(io.Discard))
			if cmd.Flag("all").Changed {
				opts = append(opts, WithAssetFilter(nil))
			}
			assets, err := New(owner, repo[c.Args[0]], opts...).Assets(c.Args[1])
			if err != nil {
				return carapace.ActionMessage(err.Error())
			}
			return carapace.ActionValues(assets...).StyleF(style.ForPathExt)
		}),
	)

	return cmd
}
