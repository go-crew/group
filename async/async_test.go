package async

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"
)

// 同步协程测试
func TestAsync_Run(t *testing.T) {
	ctx := context.Background()
	gp := NewGroup()
	// 添加第一个异步协程
	gp.Add(func(ctx context.Context, params ...interface{}) (err error) {
		log.Printf("start no 1, param1:%d, param2:%d", params[0], params[1])
		select {
		case <-time.After(2 * time.Second):
			return errors.New("no 1 is timeout")
		}
	}, func(cancel context.CancelFunc, err error) {
		cancel()
		log.Printf("no 1 收到错误：%s", err.Error())
	}, 10, 20)

	// 添加第二个异步协程
	gp.Add(func(ctx context.Context, params ...interface{}) (err error) {
		log.Println("start no 2")
		select {
		case <-time.After(5 * time.Second):
			return errors.New("no 2 is timeout")
		case <-ctx.Done():
			return errors.New("no 2 is canceled")
		}
	}, func(cancel context.CancelFunc, err error) {
		log.Printf("no 2 收到错误：%s", err.Error())
		cancel()
	})

	// 运行
	// 可以设置超时，如：2 * time.Second，永不超时为 -1
	if err := gp.Run(ctx, -1); nil != err {
		log.Println(err)
	}

	time.Sleep(20 * time.Second)

}
