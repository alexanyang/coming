package coming

import (
	"errors"
	"fmt"
	"log"
)

// Task struct should be used for task that will running for long time so
// that you can use other methods like Pause() before finish
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

// TaskItem e.g.
// func(task *Task,index int){
//    //do something
//    //if there is result you need ,you can put it into task.data[index]
//    task.data[index] = result
//}
type TaskItem func(*Task, int) error

// New create a new Task schedule,you can use method to operate you taskItem
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

// Add function add a TaskItem into Task
func (task *Task) Add(t TaskItem) *Task {
	if task.works == nil {
		panic("please use Task.New() create Task")
	}
	if task.status {
		panic("This Task had start ,can not add")
	}
	task.works = append(task.works, t)
	return task
}

// Start function will start a waiting for run TaskItem of slice
func (task *Task) Start() error {
	task.status = true
	go func() {
		<-task.start
		for index, _ := range task.works {
			select {
			case <-task.stop:
				//关闭其他通道
				close(task.start)
				close(task.pause)
				close(task.stop)
				log.Print(&task, " has stop")
				goto END
			case <-task.pause:
				task.pause <- 2
			default:
			}
			task.Exec(index)
		}
	END:
		task.finish = true
		task.wait <- 1
		task.status = false
	}()
	return nil
}

// Exec will execute taskItem function in
func (task *Task) Exec(index int) {
	task.works[index](task, index)
}

func (task *Task) WaitFinish() int {
	result := <-task.wait
	close(task.wait)
	return result
}

//当需要再次执行Start时可以重新刷新,在此之前需要执行stop,并等待关闭完成
func (task *Task) Flush() bool {
	if task.finish {
		task.start = make(chan int, 1)
		task.stop = make(chan int, 1)
		task.pause = make(chan int)
		task.wait = make(chan int)
		task.status = false
		task.finish = false
		return true
	}
	return false
}

func (task *Task) Stop() error {
	var err error
	defer func() {
		if err1 := recover(); err != nil {
			err = fmt.Errorf("%v", err1)
		}
	}()
	task.stop <- 1
	return err
}

func (task *Task) Run() error {
	var err error
	if !task.status {
		err = errors.New("please do Start() function before Run()")
		return err
	}
	task.start <- 1
	return err
}

// 如果多处触发Pause()可能会导致阻塞，建议 go task.Pause(),利用channel阻塞进行暂停
func (task *Task) Pause() {
	task.pause <- 1
}

// 如果多处触发Resume()可能会导致阻塞，建议 go task.Pause()
func (task *Task) Resume() {
	<-task.pause
}

//task need user put result of TaskItem into data ,and key is index ,value is result
func (task *Task) Call() (map[int]interface{}, error) {
	return task.data, nil
}
