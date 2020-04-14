package task

//任务主体，需要抽象属性合实现方法
type Task struct {
}

type Runner func()

//定义可执行的方法
type Executable interface {
	Start() (n int, err error)
	Stop() (err error)
}
