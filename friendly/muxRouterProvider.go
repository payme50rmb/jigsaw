package friendly

import (
	"context"
	"github.com/payme50rmb/jigsaw/contract"
	"github.com/payme50rmb/jigsaw/pkg/logger"
	"log"
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
	server   *http.Server
	router   *mux.Router
	isClosed bool
	done     chan struct{}
	log      logger.Logger
}

func NewMuxRouterProvider(c contract.Core) *MuxRouterProvider {
	return &MuxRouterProvider{
		c:   c,
		log: logger.New("init", "MuxRouterProvider"),
	}
}

func (p *MuxRouterProvider) Name() contract.ModuleName {
	return "default.providers.muxRouter"
}

func (p *MuxRouterProvider) Apply(ctx context.Context) error {
	if p.router == nil {
		p.router = mux.NewRouter()
	}
	for name, module := range p.c.Modules() {
		if contract, ok := module.(UseMuxRouter); ok {
			log.Printf("the module %s is register at MuxRouterProvider \n", name)
			contract.UseMuxRouter(p.router)
		}
	}
	return nil
}

func (p *MuxRouterProvider) Run(ctx context.Context) error {
	if p.router == nil {
		return nil
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
		defer func() {
			_ = p.Close()
		}()
		p.log.Info("the MuxRouterProvider is start")
		if err := p.server.ListenAndServe(); err != nil {
			p.log.Info("the MuxRouterProvider is closed")
		}
	}()
	for {
		select {
		case <-ctx.Done():
			p.log.Info("the MuxRouterProvider is close by ctx")
			_ = p.Close()
			return nil
		case <-p.done:
			p.log.Info("the MuxRouterProvider is close by done")
			return nil
		}
	}
}

func (p *MuxRouterProvider) Close() error {
	p.done <- struct{}{}
	defer func() {
		p.isClosed = true
	}()
	if p.server == nil {
		return nil
	}
	_ = p.server.Close()
	if p.isClosed {
		return nil
	}
	return nil
}
