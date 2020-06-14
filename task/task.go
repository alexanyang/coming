package task

import (
	"errors"
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

//适合执行顺序长时间任务
type Task struct {
	stop   chan int
	start  chan int
	pause  chan int
	works  []TaskItem
	status bool
	finish bool
	data   map[int]interface{}
	wg     sync.WaitGroup
}

//func(this *Task,index int){
//    //do something
//    //if there is result you need ,you can put it into this.data[index]
//    this.data[index] = result
//}
type TaskItem func(*Task, int) error

func (Task) New() *Task {
	return &Task{
		stop:   make(chan int),
		start:  make(chan int),
		pause:  make(chan int, 1),
		works:  make([]TaskItem, 0),
		status: false,
		finish: false,
		data:   make(map[int]interface{}),
	}
}

//
func (this *Task) Add(t TaskItem) {
	if this.works == nil {
		panic("please use Task.New() create Task")
	}
	this.works = append(this.works, t)
}

//this function will start a waiting for run TaskItem of slice
func (this *Task) Start() error {
	this.status = true
	go func() {
		<-this.start
		for index, work := range this.works {
			select {
			case <-this.stop:
				//关闭其他通道
				close(this.start)
				close(this.pause)
				close(this.stop)
				this.finish = true
				return
			default:
				this.pause <- index
				work(this, index)
				<-this.pause
				return
			}
		}
		this.finish = true
	}()
	return nil
}

//当需要再次执行Start时可以重新刷新
func (this *Task) flush() bool {
	if this.finish {
		close(this.start)
		close(this.stop)
		close(this.pause)
		this.start = make(chan int)
		this.stop = make(chan int)
		this.pause = make(chan int, 1)
		this.status = false
		this.finish = false
		return true
	}
	return false
}

func (this *Task) Stop() error {
	var err error
	defer func() {
		if err1 := recover(); err != nil {
			err = fmt.Errorf("%v", err1)
		}
	}()
	this.stop <- 1
	return err
}

func (this *Task) Run() error {
	var err error
	if !this.status {
		err = errors.New("please do Start() function before Run()")
		return err
	}
	this.start <- 1
	return err
}

// 如果多处触发Pause()可能会导致阻塞，建议 go this.Pause()
func (this *Task) Pause() {
	this.pause <- 1
}

// 如果多处触发Resume()可能会导致阻塞，建议 go this.Pause()
func (this *Task) Resume() {
	<-this.pause
}

//this need user put result of TaskItem into data ,and key is index ,value is result
func (this *Task) Call() (interface{}, error) {
	return this.data, nil
}
