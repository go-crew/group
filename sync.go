package group

// 动作
type syncActor struct {
	*Result
	execute AddResultContext
	Params  []interface{}
}

func (a syncActor) result(val interface{}, err error) {
	a.val = val
	a.err = err
}

// 添加同步有返回值的操作
func (g *Group) AddSync(execute AddResultContext, params ...interface{}) (res *Result) {
	g.setStatus(Sync)
	res = new(Result)
	g.actors = append(g.actors, &syncActor{res, execute, params})

	return res
}

// 同步执行
func (g *Group) sync() (err error) {
	l := len(g.actors)
	if 0 == l {
		return
	}

	g.wg.Add(l)
	errCh := make(chan error, l)
	for _, a := range g.actors {
		go func(a actor) {
			defer g.wg.Done()
			if act, ok := a.(*syncActor); ok {
				val, err := act.execute(g.ctx, act.Params)
				act.result(val, err)
				errCh <- err
			}
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
