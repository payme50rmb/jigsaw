package friendly

import (
	"context"
	"errors"
	"github.com/payme50rmb/jigsaw/contract"
	"github.com/payme50rmb/jigsaw/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type UseMuxRouter interface {
	UseMuxRouter(r *mux.Router)
}

type MuxRouterConfig struct {
	Addr string
}

type MuxRouterProvider struct {
	c        contract.Core
	log      logger.Logger
	server   *http.Server
	router   *mux.Router
	isClosed bool
	done     chan struct{}
}

func NewMuxRouterProvider(c contract.Core) *MuxRouterProvider {
	return &MuxRouterProvider{
		c:    c,
		log:  logger.New(),
		done: make(chan struct{}),
	}
}

func (p *MuxRouterProvider) Name() contract.ModuleName {
	return "default.providers.muxRouter"
}

func (p *MuxRouterProvider) Apply(ctx context.Context) error {
	if p.router == nil {
		p.router = mux.NewRouter()
	}
	for _, module := range p.c.Modules() {
		if m, ok := module.(UseMuxRouter); ok {
			m.UseMuxRouter(p.router)
		}
	}
	return nil
}

func (p *MuxRouterProvider) Run(ctx context.Context) error {
	if p.router == nil {
		return errors.New("the MuxRouterProvider is not register any router")
	}
	var cfg MuxRouterConfig
	err := viper.UnmarshalKey("http", &cfg)
	if err != nil {
		return err
	}
	p.server = &http.Server{
		Addr:    cfg.Addr,
		Handler: p.router,
	}
	go func() {
		p.log.F("addr", cfg.Addr).Info("the MuxRouterProvider is start")
		if err := p.server.ListenAndServe(); err != nil {
			p.log.F("addr", cfg.Addr).Error("ListenAndServe failed", err)
			p.done <- struct{}{}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			p.log.F("addr", cfg.Addr).Info("the MuxRouterProvider is close by ctx")
			return p.Close()
		case <-p.done:
			p.log.F("addr", cfg.Addr).Info("the MuxRouterProvider is close by done")
			return p.Close()
		}
	}
}

func (p *MuxRouterProvider) Close() error {
	if p.isClosed {
		return nil
	}
	defer func() {
		p.isClosed = true
	}()
	if p.server != nil {
		_ = p.server.Close()
	}
	return nil
}
