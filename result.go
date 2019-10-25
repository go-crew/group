package group

// 结果集对象
type Result struct {
	val interface{}
	err error
}

// 获取值
func (res Result) Val() interface{} {
	return res.val
}

// 获取错误
func (res Result) Err() error {
	return res.err
}