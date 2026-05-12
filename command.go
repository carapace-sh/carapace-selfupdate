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
			repository, ok := repo[args[0]]
			if !ok {
				return fmt.Errorf("invalid source %q", args[0])
			}

			opts := append([]option{WithProgress(cmd.ErrOrStderr())}, opts...)
			if cmd.Flag("force").Changed {
				opts = append(opts, WithForce(true))
			}

			var tag, asset string
			if len(args) > 1 {
				tag = args[1]
			}
			if len(args) > 2 {
				asset = args[2]
			}

			c := New(owner, repository, opts...)
			resolvedTag, resolvedAsset, err := c.resolve(tag, asset)
			if err != nil {
				return err
			}
			if tag == "" {
				c.Printf("selected tag %#v\n", resolvedTag)
			}
			if asset == "" {
				c.Printf("selected asset %#v\n", resolvedAsset)
			}

			return c.Install(resolvedTag, resolvedAsset)
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
			if len(c.Args) == 0 {
				return carapace.ActionMessage("missing source")
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
			if len(c.Args) < 2 {
				return carapace.ActionMessage("missing tag")
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
