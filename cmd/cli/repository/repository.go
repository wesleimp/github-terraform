package repository

import (
	"os"

	"github.com/apex/log"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/wesleimp/github-terraform/cmd/cli/client"
	eerror "github.com/wesleimp/github-terraform/cmd/cli/error"
	"github.com/wesleimp/github-terraform/cmd/cli/repository/collaborator"
	"github.com/wesleimp/github-terraform/cmd/cli/repository/project"
	"github.com/wesleimp/github-terraform/pkg/config"
	"github.com/wesleimp/github-terraform/pkg/context"
	"github.com/wesleimp/github-terraform/pkg/repository"
)

// Cmd config
type Cmd struct {
	Cmd     *cobra.Command
	options options
}

type options struct {
	name     string
	user     string
	dest     string
	repoType string
	token    string
	perPage  int
	page     int
}

// NewCmd creates a repository cmd
func NewCmd() *Cmd {
	root := &Cmd{}

	var commands = &cobra.Command{
		Use:           "repositories",
		Aliases:       []string{"repos"},
		Short:         "Import user repositories",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info(color.New(color.Bold).Sprint("Importing..."))

			_, err := startImport(root.options)
			if err != nil {
				return eerror.Wrap(err, color.New(color.Bold).Sprintf("Import error"))
			}

			log.Infof(color.New(color.Bold).Sprintf("Import Succeeded"))
			return nil
		},
	}

	commands.Flags().StringVarP(&root.options.name, "name", "n", "", `Repository name. The name must contains owner/repo`)
	commands.Flags().StringVarP(&root.options.user, "user", "u", "", "Repository user")
	commands.Flags().StringVarP(&root.options.dest, "dest", "d", "./output", "Path that will contains the output files")
	commands.Flags().StringVarP(&root.options.repoType, "type", "t", "", "Repository type. Could be public or private")
	commands.Flags().StringVar(&root.options.token, "token", "", "Github token. This property is not necessary if you already exported GITHUB_TOKEN")
	commands.Flags().IntVar(&root.options.perPage, "per-page", 100, "Items per page")
	commands.Flags().IntVar(&root.options.page, "page", 1, "Current page")

	commands.AddCommand(
		collaborator.NewCmd().Cmd,
		project.NewCmd().Cmd,
	)

	root.Cmd = commands
	return root
}

func startImport(o options) (*context.Context, error) {
	ctx := context.New(&config.Config{
		Repository: config.Repository{},
	})
	ctx = setupContext(ctx, o)

	err := repository.Import(ctx)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func setupContext(ctx *context.Context, o options) *context.Context {
	ctx.Config.Repository.Dest = o.dest
	ctx.Config.Repository.Name = o.name
	ctx.Config.Repository.User = o.user
	ctx.Config.Repository.Type = o.repoType
	ctx.Config.Repository.PerPage = o.perPage
	ctx.Config.Repository.Page = o.page

	if o.token == "" {
		ctx.Token = os.Getenv("GITHUB_TOKEN")
	} else {
		ctx.Token = o.token
	}

	ctx.Client = client.New(ctx)

	return ctx
}
