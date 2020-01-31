package sync

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"
)

// 同步协程测试
func TestSync_Run(t *testing.T) {
	ctx := context.Background()
	gp := NewGroup()
	if err := gp.Add("one", func(ctx context.Context, task TaskResult, params ...interface{}) TaskResult {
		task.Data = params[0]
		time.Sleep(1 * time.Second)
		return task
	}, 100); nil != err {
		log.Println(err)
	}

	if err := gp.Add("two", func(ctx context.Context, task TaskResult, params ...interface{}) TaskResult {
		task.Data = params[0]
		task.Err = errors.New("timeout")
		time.Sleep(2 * time.Second)
		return task
	}, 200); nil != err {
		log.Println(err)
	}

	res := gp.Run(ctx)
	log.Println(res["one"], res["two"])
}
