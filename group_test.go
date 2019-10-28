package group

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"
)

// 异步协程测试
func TestGroup_Run(t *testing.T) {
	gp := NewGroup()
	// 第一个Go程
	gp.AddAsync(func(ctx context.Context) error {
		select {
		case <-time.After(2 * time.Second):
			return errors.New("断开1")
		case <-ctx.Done():
			return ctx.Err()
		}
	}, func(cancel context.CancelFunc, e error) {
		fmt.Println("one:", e)
		cancel()
	})

	// 第二个Go程
	gp.AddAsync(func(ctx context.Context) error {
		select {
		case <-time.After(2 * time.Second):
			return errors.New("断开2")
		case <-ctx.Done():
			return ctx.Err()
		}
	}, func(cancel context.CancelFunc, e error) {
		fmt.Println("two:", e)
		cancel()
	})

	// 第三个Go程
	gp.AddAsync(func(ctx context.Context) error {
		select {
		case <-time.After(2 * time.Second):
			return errors.New("超时错误")
		case <-ctx.Done():
			return ctx.Err()
		}
	}, func(cancel context.CancelFunc, e error) {
		fmt.Println("three:", e)
		cancel()
	})

	// 运行所有使用group对象的goroutine，模拟有一个发生错误，取消所有执行
	if err := gp.Run(); nil != err {
		fmt.Println("Run:", err)
	}
}

// 同步协程版
func TestGroup_Wait(t *testing.T) {
	gp := NewGroup()
	res1 := gp.AddSync(func(ctx context.Context, params ...interface{}) (val interface{}, err error) {
		time.Sleep(1 * time.Second)
		val = nil
		err = errors.New("发生错误1")
		return
	})

	res2 := gp.AddSync(func(ctx context.Context, params ...interface{}) (val interface{}, err error) {
		time.Sleep(2 * time.Second)
		val = 200
		err = nil
		return
	})

	res3 := gp.AddSync(func(ctx context.Context, params ...interface{}) (val interface{}, err error) {
		time.Sleep(3 * time.Second)
		val = 300
		err = nil
		return
	})

	res5 := gp.AddSync(func(ctx context.Context, params ...interface{}) (val interface{}, err error) {
		time.Sleep(5 * time.Second)
		val = nil
		err = errors.New("发生错误5")
		return
	})

	// 运行
	if err := gp.Run(); err != nil {
		fmt.Println(err)
	}

	fmt.Println(res1.Val(), res1.Err())
	fmt.Println(res2.Val(), res2.Err())
	fmt.Println(res3.Val(), res3.Err())
	fmt.Println(res5.Val(), res5.Err())
}

// 同步传递参数
func TestGroup_Wait2(t *testing.T) {
	gp := NewGroup()
	// 设置context的值，让所有同步协程共享
	gp.SetContext("key", "value")
	s := make([]interface{}, 10)
	m := make(map[int]string, 10)
	for i:= 0; i<10; i++ {
		m[i] = "str" + strconv.Itoa(i)
	}

	for key, item := range m {
		// 可以传入参数
		s[key] = gp.AddSync(func(ctx context.Context, params ...interface{}) (i interface{}, e error) {
			// 所有协程中都可以获取ctx设置的值
			log.Println(ctx.Value("key"))
			i = params[0]
			e = nil
			return
		}, item)
	}

	// 启动同步协程任务
	if err := gp.Run(); err != nil {
		log.Println(err)
	}

	// 输出
	for key, item := range s {
		if val, ok := item.(*Result); ok {
			log.Println(key, val.Val())
		}

	}
}
