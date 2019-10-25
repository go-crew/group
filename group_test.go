package group

import (
	"context"
	"errors"
	"fmt"
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
	res1 := gp.AddSync(func(ctx context.Context) (val interface{}, err error) {
		time.Sleep(2 * time.Second)
		val = 100
		err = nil
		return
	})

	res2 := gp.AddSync(func(ctx context.Context) (val interface{}, err error) {
		time.Sleep(3 * time.Second)
		val = 200
		err = nil
		return
	})

	res3 := gp.AddSync(func(ctx context.Context) (val interface{}, err error) {
		time.Sleep(5 * time.Second)
		val = nil
		err = errors.New("发生错误")
		return
	})

	// 运行
	if err := gp.Run(); err != nil {
		fmt.Println(err)
	}

	fmt.Println(res1.Val(), res1.Err())
	fmt.Println(res2.Val(), res2.Err())
	fmt.Println(res3.Val(), res3.Err())
}
