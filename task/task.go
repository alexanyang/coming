package task

import (
	"errors"
	"fmt"
)

//任务主体，需要抽象属性合实现方法
//任务主题需要存在一个私有字段int 用以表示数据的执行状态,需要保证这个字段的原子性
// 需要一个就收chan int  用以接收数据的状态变化,在数据结束后返回一个int值进行表示
// 可执行的代码不能直接执行需要与开始和结束嵌套以获取执行状态,当执行完成后通过通道表示数据结束
// 需要通过*chan 保存一个通道，必须只想地址，方便在调用的时候能，通道功效，需要用到通道的方法也必须使用*struct{}来去对象本身的引用

//   todo 需要注意通道在缓存不足的情况下会进入阻塞状态

type Task struct {
	//当我们的工作状态进行转变的时候通过这个通道来传递方法执行完毕的信息
	status chan bool //通过这个通带进行状态的返回
	//任务的执行程序，最终我们讨论的是这个执行程序的状态
	runner Runable //这个
	//用来装返回的数据参数
	resultData []byte
	//当meta执行完成，标为true
	finished bool
}

func NewTask(runner Runable) *Task {
	return &Task{
		status:     make(chan bool, 1),
		runner:     runner,
		resultData: make([]byte, 0),
		finished:   false,
	}
}

func (this *Task) Finish() {
	this.finished = true
}

func (this *Task) Run() error {

	defer func() {
		if r := recover(); r != nil {
			fmt.Sprintf("Execute task err : ", r)
		}
	}()

	bytes, err := this.runner.Run()
	if err == nil {
		fmt.Println(string(bytes))
		//copy不会改变目标数据长度，所以需要重新定义
		this.resultData = make([]byte, len(bytes))
		copy(this.resultData, bytes)
	}
	this.status <- true
	this.Finish()
	return err
}

//当程序开始中的时候，我们需要在协程中使用一个通道进行信息传递，当程序执行完成，用通道接收，改变任务的状态
//todo 待完成
func (this *Task) Start() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Start task err : ", r))
		}
	}()

	go this.Run()
	fmt.Println("任务执行中")
	return err
}

func (this *Task) GetResult() []byte {
	return this.resultData
}

func (this *Task) IsFinished() bool {
	this.finished = <-this.status
	return this.finished
}

//虽然兼容返回参数,但是并不执行对返回参数的获取,实现通信后,对回执参数进行处理目前准备采取[]byte进行处理返回,由用户自己进行反序列化
//r是返回结果的序列化数据,提供给实现方用以完成数据返回,如果未能获取任何数据返回空的[]byte对象,如果在获取结果的过程中失败了,将错误
//尝试使用。。。interface进行任意参数接受失败，需要进行数据格式转化才可以
type Runner func(ops ...Runable) (r []byte, err error)

//由于 Runner 类型存在具体参数如何将这个参数 输入到这个函数

func (r Runner) Run() (bytes []byte, err error) {
	bytes, err = r()
	return
}

//定义可执行的方法
type Executable interface {
	Start() error
	Stop() error
}

//实现改接口,以调用Stop()进行停止当前携程,本质是在执行过程中,通过通道传递信号,使得协程停止
type Stopable interface {
	Stop() error
}

type Startable interface {
	Start() error
}

type Runable interface {
	Run() ([]byte, error)
}
