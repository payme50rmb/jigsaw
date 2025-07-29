package jigsaw

import (
	"github.com/oklog/run"
	"github.com/payme50rmb/jigsaw/contract"
	"github.com/payme50rmb/jigsaw/friendly"
	"github.com/payme50rmb/jigsaw/pkg/logger"
	"github.com/spf13/cobra"
)

func New() contract.Core {
	return &C{
		modules: make(map[string]contract.Module),
		g:       &run.Group{},
		log:     logger.New("init", "monica"),
		root:    nil,
	}
}

func Default() contract.Core {
	c := New()
	c.Register(friendly.NewSignal(c))
	return c
}

func Commandable() contract.Core {
	c := New()
	c.Register(friendly.NewSignal(c))
	c.Register(friendly.NewCobraCommandProvider(c))
	c.Register(friendly.NewLogger())
	return c
}

func CommandableWithRoot(root *cobra.Command) contract.Core {
	c := New()
	c.Register(friendly.NewSignal(c))
	c.Register(friendly.NewCobraCommandProviderWithRoot(c, root))
	return c
}
