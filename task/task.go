package task

//任务主体，需要抽象属性合实现方法
//任务主题需要存在一个私有字段int 用以表示数据的执行状态,需要保证这个字段的原子性
// 需要一个就收chan int  用以接收数据的状态变化,在数据结束后返回一个int值进行表示
// 可执行的代码不能直接执行需要与开始和结束嵌套以获取执行状态,当执行完成后通过通道表示数据结束
type Task struct {
}

//虽然兼容返回参数,但是并不执行对返回参数的获取,实现通信后,对回执参数进行处理目前准备采取[]byte进行处理返回,由用户自己进行反序列化
//r是返回结果的序列化数据,提供给实现方用以完成数据返回,如果未能获取任何数据返回空的[]byte对象,如果在获取结果的过程中失败了,将错误
type Runner func(...interface{}) (r []byte,err error)

//定义可执行的方法
type Executable interface {
	Start() (n int, err error)
	Stop() (err error)
}
//实现改接口,以调用Stop()进行停止当前携程,本质是在执行过程中,通过通道传递信号,使得协程停止
type Stopable interface {
	Stop() (err error)
}
