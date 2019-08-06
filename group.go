package group

import (
	"context"
	"errors"
	"log"
	"sync"
)

// Go程组结构对象
type Group struct {
	ctx     context.Context
	cancel  context.CancelFunc
	actors  []actor
	wg      sync.WaitGroup
	onceErr sync.Once
}

// 创建Group对象
func NewGroup() *Group {
	ctx, cancel := context.WithCancel(context.Background())
	return &Group{
		ctx:    ctx,
		cancel: cancel,
	}
}

// 同步Go程执行，一个Go程出错，其他Go程继续执行，一直到所有Go程都完成
func (g *Group) Wait() (err error) {
	l := len(g.actors)
	if 0 == l {
		return
	}

	g.wg.Add(l)
	errCh := make(chan error, l)
	for _, a := range g.actors {
		go func(a actor) {
			defer func() {
				err := recover()
				ecp := catchGroupException(err)
				errCh <- ecp
			}()
			e := a.execute(g.ctx)
			a.interrupt(g.cancel, e)
			errCh <- e
			g.wg.Done()
		}(a)
	}

	for i := 0; i < cap(errCh); i++ {
		e := <-errCh
		if nil != e {
			g.onceErr.Do(func() {
				err = e
			})
		}
	}

	g.wg.Wait()
	close(errCh)
	return
}

// Go程执行环境
type AddContext func(context.Context) error

// Go程终止环境
type AddInterrupt func(context.CancelFunc, error)

// 添加Go程运行程序
func (g *Group) Add(execute AddContext, interrupt AddInterrupt) {
	g.actors = append(g.actors, actor{execute, interrupt})
}

// 异步执行Go程，一个Go程终止，所有Go程连带一起终止
func (g *Group) Run() (err error) {
	if len(g.actors) == 0 {
		return nil
	}

	errCh := make(chan error, len(g.actors))
	for _, a := range g.actors {
		go func(a actor) {
			defer func() {
				err := recover()
				ecp := catchGroupException(err)
				errCh <- ecp
			}()
			errCh <- a.execute(g.ctx)
		}(a)
	}

	err = <-errCh
	for _, a := range g.actors {
		a.interrupt(g.cancel, err)
	}

	for i := 1; i < cap(errCh); i++ {
		<-errCh
	}

	return
}

// 获取ctx顶级对象
func (g *Group) Ctx() context.Context {
	return g.ctx
}

// 获取取消对象
func (g *Group) Cancel() context.CancelFunc {
	return g.cancel
}

type actor struct {
	execute   AddContext
	interrupt AddInterrupt
}

// 处理group触发的panic，并打印
func catchGroupException(err interface{}) error {
	var errMsg error
	switch err.(type) {
	case string:
		msg := err.(string)
		errMsg = errors.New(msg)
	case error:
		errMsg = err.(error)
	}

	if nil != errMsg {
		log.Println("group panic:" + errMsg.Error())
	}

	return errMsg
}
