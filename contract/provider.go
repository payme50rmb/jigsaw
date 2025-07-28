package contract

import (
	"context"
)

// Provider 定义了一个提供依赖的接口
type Provider interface {
	Module
	Apply(ctx context.Context) error
}

type Runnable interface {
	Run(ctx context.Context) error
}

type Closer interface {
	Close() error
}

// Rootable 替换掉根的实现
type Rootable interface {
	Root(Core) Runnable
}

// Initable 定义了一个初始化的接口, 会先于所有 module 和 provider 的执行
type Initable interface {
	Init() error
}
