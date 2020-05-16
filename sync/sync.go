package sync

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-crew/group"
)

type (
	// 协程运行时给予的标签
	Tag string
	// Go程业务逻辑执行环境
	Execute func(ctx context.Context, task TaskResult, params ...interface{}) TaskResult
)

// 同步业务对象
type Sync struct {
	tags   []Tag
	actors []*actor
	length int
}

const tagErrMsg = "tag %s exists"

// 创建分组对象
func NewGroup() *Sync {
	return &Sync{}
}

// 添加方法
func (s *Sync) Add(tag Tag, exec Execute, params ...interface{}) error {
	if exists(s.tags, tag) {
		errMsg := fmt.Sprintf(tagErrMsg, tag)
		return errors.New(errMsg)
	}

	actor := &actor{
		TaskResult: TaskResult{
			tag: tag,
		},
		params: params,
		exec:   exec,
	}

	s.tags = append(s.tags, tag)
	s.actors = append(s.actors, actor)
	s.length++

	return nil
}

// 运行
func (s *Sync) Run(ctx context.Context) (mTask map[Tag]*TaskResult) {
	l := len(s.actors)
	if 0 == l {
		return
	}

	taskCh := make(chan TaskResult)
	mTask = make(map[Tag]*TaskResult, l)
	for _, act := range s.actors {
		go func(ctx context.Context, act *actor) {
			task := act.TaskResult
			defer func() {
				if err := group.CatchPanic(recover()); nil != err {
					task.Err = err
					taskCh <- task
				}
			}()
			taskCh <- act.exec(ctx, task, act.params...)
		}(ctx, act)
	}

	for i := 0; i < l; i++ {
		if task, ok := <-taskCh; ok {
			mTask[task.tag] = &task
		}
	}

	close(taskCh)
	return
}

// 任务处理结果对象
type TaskResult struct {
	tag  Tag
	Data interface{}
	Err  error
}

// 获取tag标签
func (tr TaskResult) Tag() string {
	return string(tr.tag)
}

// 任务对象
type actor struct {
	TaskResult

	params []interface{}
	exec   Execute
}

// // 校验标签是否已经存在
func exists(tags []Tag, tag Tag) (ex bool) {
	if 0 == len(tags) {
		return
	}

	for _, t := range tags {
		if t == tag {
			ex = true
			return
		}
	}
	return
}
