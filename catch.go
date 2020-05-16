package group

import (
	"errors"
	"log"
)

// 处理group触发的panic，并打印
func CatchPanic(errInfo interface{}) (err error) {
	switch errInfo.(type) {
	case string:
		msg := errInfo.(string)
		err = errors.New(msg)
	case error:
		err = err.(error)
	}

	if nil != err {
		log.Println("group panic:" + err.Error())
	}

	return
}
