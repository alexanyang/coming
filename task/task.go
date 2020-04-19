package task

import (
	"fmt"
)

//任务主体，需要抽象属性合实现方法
//任务主题需要存在一个私有字段int 用以表示数据的执行状态,需要保证这个字段的原子性
// 需要一个就收chan int  用以接收数据的状态变化,在数据结束后返回一个int值进行表示
// 可执行的代码不能直接执行需要与开始和结束嵌套以获取执行状态,当执行完成后通过通道表示数据结束

//   todo 需要注意通道在缓存不足的情况下会进入阻塞状态

type Task struct {
	//当我们的工作状态进行转变的时候通过这个通道来传递方法执行完毕的信息
	status chan bool //通过这个通带进行状态的返回
	//任务的执行程序，最终我们讨论的是这个执行程序的状态
	runner Runner //这个
	//用来保存runner的输入参数
	Parameter []interface{}
	//用来装返回的数据参数
	resultData []byte
	//当meta执行完成，标为true
	finished bool
}

func NewTask(runner Runner, ops ...interface{}) *Task {
	return &Task{
		status:     make(chan bool, 1),
		runner:     runner,
		resultData: make([]byte, 1),
		finished:   false,
		Parameter:  ops,
	}
}
func (this Task) Finish() {
	this.finished = true
}

func (this Task) DoSome(t *Task) error {

	defer func() {
		if r := recover(); r != nil {
			fmt.Sprintf("Execute task err : ", r)
		}
	}()
	//todo 缺少参数列表
	var (
		bytes []byte
		err   error
	)
	if this.Parameter == nil || len(this.Parameter) == 0 {
		bytes, err = this.runner()
	} else {
		bytes, err = this.runner(this.Parameter)
	}

	if err == nil {
		t.resultData = bytes
	}

	t.status <- true
	return err
}

//当程序开始中的时候，我们需要在协程中使用一个通道进行信息传递，当程序执行完成，用通道接收，改变任务的状态
//todo 待完成
func (this Task) Start() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			fmt.Sprintf("Start task err : ", r)
		}
	}()
	go this.DoSome(&this)
	select {
	case ok := <-this.status:
		if ok {
			this.Finish()
		}
	default:
	}

	return err
}
func (this Task) GetResult() []byte {
	return this.resultData
}

func (this Task) IsFinished() bool {
	return this.finished
}

//虽然兼容返回参数,但是并不执行对返回参数的获取,实现通信后,对回执参数进行处理目前准备采取[]byte进行处理返回,由用户自己进行反序列化
//r是返回结果的序列化数据,提供给实现方用以完成数据返回,如果未能获取任何数据返回空的[]byte对象,如果在获取结果的过程中失败了,将错误
//尝试使用。。。interface进行任意参数接受失败，需要进行数据格式转化才可以
type Runner func(ops ...interface{}) (r []byte, err error)

//由于 Runner 类型存在具体参数如何将这个参数 输入到这个函数

func NewRunner() {

}

//定义可执行的方法
type Executable interface {
	Start() (n int, err error)
	Stop() (err error)
}

//实现改接口,以调用Stop()进行停止当前携程,本质是在执行过程中,通过通道传递信号,使得协程停止
type Stopable interface {
	Stop() (err error)
}

type Startable interface {
	Start() error
}
