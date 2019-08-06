package group

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

// run测试方法
func TestGroup_Run(t *testing.T) {
	gp := NewGroup()
	// 第一个Go程
	gp.Add(func(ctx context.Context) error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-ctx.Done():
			return ctx.Err()
		}
	}, func(cancel context.CancelFunc, e error) {
		cancel()
	})

	// 第二个Go程
	gp.Add(func(ctx context.Context) error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-time.After(3 * time.Second):
			return errors.New("超时错误")
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-ctx.Done():
			return ctx.Err()
		}

	}, func(cancel context.CancelFunc, e error) {
		cancel()
	})

	// 运行所有使用group对象的goroutine，当有一个超时时，取消所有
	if err := gp.Run(); nil != err {
		fmt.Println(err) // 超时错误
	}
}

func TestGroup_Wait(t *testing.T) {
	gp := NewGroup()

	gp.Add(func(ctx context.Context) error {
		time.Sleep(2 * time.Second)
		return errors.New("第一个Go程出错")
	}, func(cancel context.CancelFunc, e error) {
		log.Printf("1号执行完毕, 错误信息：%s", e.Error())
	})

	gp.Add(func(ctx context.Context) error {
		time.Sleep(3 * time.Second)
		return errors.New("第二个Go程出错")
	}, func(cancel context.CancelFunc, e error) {
		log.Printf("2号执行完毕, 错误信息：%s", e.Error())
	})

	gp.Add(func(ctx context.Context) error {
		time.Sleep(3 * time.Second)
		return nil
	}, func(cancel context.CancelFunc, e error) {
		log.Printf("3号执行完毕, 错误信息：%s", e)
	})

	if err := gp.Wait(); nil != err {
		fmt.Println(err) // 超时错误
	} else {
		fmt.Println("no err")
	}
}
