package async

import (
	"context"
	"github.com/go-crew/group"
	"time"
)

type (
	// Execute Go程业务逻辑执行环境
	Execute func(ctx context.Context, params ...interface{}) (err error)
	// Interrupt Go程停止后的执行环境回调
	Interrupt func(cancel context.CancelFunc, err error)
)

// Async 同步业务对象
type Async struct {
	actors []*actor
	length int
}

// NewGroup 创建分组对象1
func NewGroup() *Async {
	return &Async{}
}

// Add 添加方法
func (a *Async) Add(exec Execute, inter Interrupt, params ...interface{}) {
	actor := &actor{
		params: params,
		exec:   exec,
		inter:  inter,
	}

	a.actors = append(a.actors, actor)
	a.length++
}

// Run 执行
// timeout设置 -1 表示不会自动超时
func (a *Async) Run(ctx context.Context, timeout time.Duration) (err error) {
	var cancel context.CancelFunc
	l := len(a.actors)
	if 0 == l {
		return
	}

	if timeout == -1 {
		ctx, cancel = context.WithCancel(ctx)
	} else {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	}

	errCh := make(chan error, l)
	for _, act := range a.actors {
		go func(ctx context.Context, act *actor) {
			defer func() {
				errCh <- group.CatchPanic(recover())
			}()
			errCh <- act.exec(ctx, act.params...)
		}(ctx, act)
	}

	err = <-errCh
	// 全局通知
	for _, act := range a.actors {
		act.inter(cancel, err)
	}

	// 清空
	for i := 1; i < cap(errCh); i++ {
		<-errCh
	}

	defer cancel()
	return
}

// 任务对象
type actor struct {
	params []interface{}
	exec   Execute
	inter  Interrupt
}
