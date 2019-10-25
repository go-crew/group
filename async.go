package group

import (
	"errors"
	"log"
)

// 动作
type asyncActor struct {
	execute   AddContext
	interrupt AddInterrupt
}

// 标记
func (a asyncActor) result(interface{}, error) {}

// 添加异步Go程操作
func (g *Group) AddAsync(execute AddContext, interrupt AddInterrupt) {
	g.setStatus(Async)
	g.actors = append(g.actors, asyncActor{execute, interrupt})
}

// 异步执行
func (g *Group) async() (err error) {
	l := len(g.actors)
	if l == 0 {
		return
	}

	errCh := make(chan error, l)
	for _, a := range g.actors {
		go func(a actor) {
			defer func() {
				err := recover()
				ecp := catchGroupException(err)
				errCh <- ecp
			}()

			if act, ok := a.(asyncActor); ok {
				errCh <- act.execute(g.ctx)
			}
		}(a)
	}

	err = <-errCh
	for _, a := range g.actors {
		if act, ok := a.(asyncActor); ok {
			act.interrupt(g.cancel, err)
		}
	}

	for i := 1; i < cap(errCh); i++ {
		<-errCh
	}

	return
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