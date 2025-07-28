package contract

import "context"

type Core interface {
	Runnable
	Closer
	Modules() map[string]Module
	Register(module Module)
	RunAsRoot(ctx context.Context) error
}
