package group

import (
	"context"
	"sync"
)

// 执行的动作接口
type actor interface {
	result(interface{}, error)
}

// Go程组结构对象
type Group struct {
	ctx     context.Context
	cancel  context.CancelFunc
	actors  []actor
	wg      sync.WaitGroup
	onceErr sync.Once
	status  Status
}

// Go程业务逻辑执行环境
type AddContext func(context.Context) error

// Go程带数据返回的执行环境
type AddResultContext func(context.Context, ...interface{}) (interface{}, error)

// Go程停止后的执行环境回调
type AddInterrupt func(context.CancelFunc, error)

// 执行的类型
type Status int

// 同步或异步类型
const (
	None Status = iota
	Async
	Sync
)

// 创建Group对象
func NewGroup() *Group {
	ctx, cancel := context.WithCancel(context.Background())
	return &Group{
		ctx:    ctx,
		cancel: cancel,
	}
}

// 异步执行Go程，一个Go程终止，所有Go程连带一起终止
func (g *Group) Run() (err error) {
	if g.status == Async {
		return g.async()
	}

	return g.sync()
}

// 获取执行Status
func (g *Group) Status() Status {
	return g.status
}

// 获取Context对象
func (g *Group) Context() context.Context {
	return g.Context()
}

// 设置
func (g *Group) setStatus(status Status) {
	if g.status == None {
		g.status = status
	}

	if g.status != status {
		panic("Must use the same type of add method")
	}
}
