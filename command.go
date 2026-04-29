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
		Args:  cobra.RangeArgs(1, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := append([]option{WithProgress(cmd.ErrOrStderr())}, opts...)
			if cmd.Flag("force").Changed {
				opts = append(opts, WithForce(true))
			}
			c := New(owner, repo[args[0]], opts...)

			if len(args) < 2 {
				tags, err := c.Tags()
				if err != nil {
					return err
				}
				if len(tags) == 0 {
					return fmt.Errorf("no tags found for %v", c.repo)
				}
				args = append(args, tags[0])
			}

			if len(args) < 3 {
				assets, err := c.Assets(args[1])
				if err != nil {
					return err
				}
				if len(assets) == 0 {
					return fmt.Errorf("no assets found for %v@%v", c.repo, args[1])
				}
				args = append(args, assets[0])
			}

			return c.Install(args[1], args[2])
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
