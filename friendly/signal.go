package friendly

import (
	"context"
	"github.com/payme50rmb/jigsaw/contract"
	"github.com/payme50rmb/jigsaw/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

type Signal struct {
	c   contract.Core
	log logger.Logger
}

func NewSignal(monica contract.Core) *Signal {
	return &Signal{
		c:   monica,
		log: logger.New("init", "Signal"),
	}
}

func (m *Signal) Name() string {
	return "default.modules.signal"
}

func (m *Signal) Run(ctx context.Context) error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-ctx.Done():
		m.log.Info("receive signal from context, shutting down")
		return m.c.Close()
	case sign := <-ch:
		m.log.F("signal", sign).Info("receive signal from os, shutting down")
		return m.c.Close()
	}
}
