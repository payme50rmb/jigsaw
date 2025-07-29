package friendly

import (
	"context"
	"github.com/payme50rmb/jigsaw/contract"
	"github.com/payme50rmb/jigsaw/pkg/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jigsaw",
	Short: "jigsaw is a monorepo management tool",
}

type UseCobraCommand interface {
	UseCobraCommand(root *cobra.Command)
}

type CobraCommandProvider struct {
	c        contract.Core
	root     *cobra.Command
	log      logger.Logger
	isClosed bool
}

func NewCobraCommandProvider(c contract.Core) *CobraCommandProvider {
	return &CobraCommandProvider{
		c:   c,
		log: logger.New("init", "CobraCommandProvider"),
	}
}

func NewCobraCommandProviderWithRoot(c contract.Core, root *cobra.Command) *CobraCommandProvider {
	p := NewCobraCommandProvider(c)
	p.root = root
	return p
}

func (p *CobraCommandProvider) Name() contract.ModuleName {
	return "default.providers.cobraCommand"
}

func (p *CobraCommandProvider) Apply(ctx context.Context) error {
	for _, module := range p.c.Modules() {
		if u, ok := module.(UseCobraCommand); ok {
			u.UseCobraCommand(p.root)
		}
	}
	return nil
}

func (p *CobraCommandProvider) Root(m contract.Core) contract.Runnable {
	p.c = m
	if p.root == nil {
		p.root = rootCmd
	}
	var configPath string
	p.root.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "config file path")
	p.root.AddCommand(&cobra.Command{
		Use:   "serve",
		Short: "serve the monica server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return p.c.RunAsRoot(cmd.Context())
		},
	})
	p.c.Register(NewConfig(configPath, "yaml"))
	return p
}

func (p *CobraCommandProvider) Run(ctx context.Context) error {
	return p.root.ExecuteContext(ctx)
}
