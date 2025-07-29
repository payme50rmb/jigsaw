package friendly

import (
	"context"

	"github.com/payme50rmb/jigsaw/contract"
	"github.com/payme50rmb/jigsaw/pkg/logger"
	"github.com/robfig/cron/v3"
)

type UseCron interface {
	UseCron(cron *cron.Cron)
}

type CronProvider struct {
	c    contract.Core
	cron *cron.Cron
	log  logger.Logger
}

func NewCronProvider(c contract.Core) *CronProvider {
	return &CronProvider{
		c:    c,
		cron: cron.New(cron.WithSeconds()),
		log:  logger.New(),
	}
}

func (p *CronProvider) Apply(ctx context.Context) error {
	for s, module := range p.c.Modules() {
		if m, ok := module.(UseCron); ok {
			m.UseCron(p.cron)
			p.log.F("module", s).Info("cron provider apply module")
		}
	}
	return nil
}

func (p *CronProvider) Run(ctx context.Context) error {
	go func() {
		p.log.Info("cron provider run")
		p.cron.Start()
	}()
	<-ctx.Done()
	p.log.Info("cron provider stop")
	p.cron.Stop()
	return nil
}
