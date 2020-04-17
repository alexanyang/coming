package task

//任务主体，需要抽象属性合实现方法
//任务主题需要存在一个私有字段int 用以表示数据的执行状态,需要保证这个字段的原子性
// 需要一个就收chan int  用以接收数据的状态变化,在数据结束后返回一个int值进行表示
// 可执行的代码不能直接执行需要与开始和结束嵌套以获取执行状态,当执行完成后通过通道表示数据结束
type Task struct {
}

type Runner func(...interface{}) (...interface{})

//定义可执行的方法
type Executable interface {
	Start() (n int, err error)
	Stop() (err error)
}
