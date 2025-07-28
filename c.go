package jigsaw

import (
	"context"
	"github.com/oklog/run"
	"github.com/payme50rmb/jigsaw/contract"
	"github.com/payme50rmb/jigsaw/pkg/logger"
	"reflect"
)

type C struct {
	modules map[string]contract.Module
	g       *run.Group
	log     logger.Logger
	root    contract.Runnable
	applied bool
}

func (c *C) Modules() map[string]contract.Module {
	return c.modules
}

func (c *C) Register(module contract.Module) {
	name := ""
	if mt, ok := module.(contract.Nameable); ok {
		name = string(mt.Name())
	} else {
		typ := reflect.TypeOf(c).Elem()
		name = typ.String()
	}
	if _, ok := c.modules[name]; ok {
		panic("module name duplicated")
	}
	c.modules[name] = module
	if r, ok := module.(contract.Rootable); ok {
		if c.root == nil {
			c.root = r.Root(c)
		} else {
			panic("root module has been set")
		}
	}
}

func (c *C) Run(ctx context.Context) error {
	if c.root == nil {
		return c.RunAsRoot(ctx)
	}
	c.apply(ctx) // apply all providers
	return c.root.Run(ctx)
}

func (c *C) RunAsRoot(ctx context.Context) error {
	c.apply(ctx) // apply all providers
	return c.g.Run()
}

func (c *C) Close() error {
	for s, module := range c.modules {
		if m, ok := module.(contract.Closer); ok {
			c.log.F("module", s).Error("the module is closed by close", nil)
			_ = m.Close()
		}
	}
	return nil
}

func (c *C) apply(ctx context.Context) {
	defer func() {
		c.applied = true
	}()
	if c.applied {
		return
	}

	// Init all modules
	for s, module := range c.modules {
		if m, ok := module.(contract.Initable); ok {
			if err := m.Init(); err != nil {
				c.log.F("module", s).Error("init module failed", err)
			}
		}
	}

	// Register all providers
	ps := make([]contract.Provider, 0)
	for _, module := range c.modules {
		if p, ok := module.(contract.Provider); ok {
			ps = append(ps, p)
		}
	}
	for _, p := range ps {
		if err := p.Apply(ctx); err != nil {
			c.log.Error("apply provider failed", err)
		}
	}

	// Collect all runnable modules
	for s, module := range c.modules {
		if r, ok := module.(contract.Runnable); ok {
			_ctx, cancel := context.WithCancel(ctx)
			c.g.Add(func() error {
				return r.Run(_ctx)
			}, func(err error) {
				cancel()
				c.log.F("module", s).Error("the module is closed", err)
				_ = r.Close()
			})
		}
	}
}
