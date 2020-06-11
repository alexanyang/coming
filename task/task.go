package task

import (
	"fmt"
	"sync"
)

//任务主体，需要抽象属性合实现方法
//任务主题需要存在一个私有字段int 用以表示数据的执行状态,需要保证这个字段的原子性
// 需要一个就收chan int  用以接收数据的状态变化,在数据结束后返回一个int值进行表示
// 可执行的代码不能直接执行需要与开始和结束嵌套以获取执行状态,当执行完成后通过通道表示数据结束
// 需要通过*chan 保存一个通道，必须只想地址，方便在调用的时候能，通道功效，需要用到通道的方法也必须使用*struct{}来去对象本身的引用

//   todo 需要注意通道在缓存不足的情况下会进入阻塞状态  目前还不能用
//定义可执行的方法

type TaskInterface interface {
	Start() error
	Run() error
	Stop() error
	New() interface{}
}

type Task struct {
	stop   chan int
	start  chan int
	works  []TaskItem
	status bool
	finish bool
	data   interface{}
	wg     sync.WaitGroup
}

type TaskItem func(Task, ...interface{}) error

func (Task) New() *Task {
	return &Task{
		stop:   make(chan int),
		start:  make(chan int),
		works:  make([]Task, 0),
		status: false,
		finish: false,
		data:   nil,
	}
}

func (this *Task) Add(t TaskItem) {
	if this.works == nil {
		panic("please use Task.New() create Task")
	}
	this.works = append(this.works, t)
}
func (this *Task) Start() error {
	go func() {
		for {
			select {
			case <-this.stop:
				//关闭其他通道
				return
			default:
				this.work(*this)
				return
			}
		}
	}()
}

func (this *Task) Stop() error {
	var err error

	defer func() {
		if err1 := recover(); err != nil {
			err = fmt.Sprint(err1)
		}
	}()

	this.stop <- 1

	return err
}

func (this *Task) Run() error {
	return this.work(*this)
}

func (this *Task) Pause() {

}

func (this *Task) Resume() {

}

func (this *Task) Call() (interface{}, error) {

}
