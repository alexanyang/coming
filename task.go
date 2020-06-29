package coming

import (
	"errors"
	"fmt"
	"log"
)

//定义可执行的方法

//适合执行顺序长时间任务
type Task struct {
	BaseTask
	status bool
	finish bool
	stop   chan int
	start  chan int
	pause  chan int
	wait   chan int
	works  []TaskItem
	data   map[int]interface{}
}

//func(this *Task,index int){
//    //do something
//    //if there is result you need ,you can put it into this.data[index]
//    this.data[index] = result
//}
type TaskItem func(*Task, int) error

func (Task) New() *Task {
	return &Task{
		status: false,
		finish: false,
		stop:   make(chan int, 1),
		start:  make(chan int, 1),
		pause:  make(chan int),
		wait:   make(chan int),
		works:  make([]TaskItem, 0),
		data:   make(map[int]interface{}),
	}
}

//
func (this *Task) Add(t TaskItem) *Task {
	if this.works == nil {
		panic("please use Task.New() create Task")
	}
	if this.status {
		panic("This Task had start ,can not add")
	}
	this.works = append(this.works, t)
	return this
}

//this function will start a waiting for run TaskItem of slice
func (this *Task) Start() error {
	this.status = true
	go func() {
		<-this.start
		for index, _ := range this.works {
			select {
			case <-this.stop:
				//关闭其他通道
				close(this.start)
				close(this.pause)
				close(this.stop)
				log.Print(&this, " has stop")
				goto END
			case <-this.pause:
				this.pause <- 2
			default:
			}
			this.Exec(index)
		}
	END:
		this.finish = true
		this.wait <- 1
		this.status = false
	}()
	return nil
}

func (this *Task) Exec(index int) {
	this.works[index](this, index)
}

func (this *Task) WaitFinish() int {
	result := <-this.wait
	close(this.wait)
	return result
}

//当需要再次执行Start时可以重新刷新,在此之前需要执行stop,并等待关闭完成
func (this *Task) Flush() bool {
	if this.finish {
		this.start = make(chan int, 1)
		this.stop = make(chan int, 1)
		this.pause = make(chan int)
		this.wait = make(chan int)
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

// 如果多处触发Pause()可能会导致阻塞，建议 go this.Pause(),利用channel阻塞进行暂停
func (this *Task) Pause() {
	this.pause <- 1
}

// 如果多处触发Resume()可能会导致阻塞，建议 go this.Pause()
func (this *Task) Resume() {
	<-this.pause
}

//this need user put result of TaskItem into data ,and key is index ,value is result
func (this *Task) Call() (map[int]interface{}, error) {
	return this.data, nil
}
