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
	status      bool
	finish      bool
	stop, start chan struct{}
	pause, wait chan struct{}
	works       []TaskItem
	data        map[int]interface{}
}

// TaskItem e.g.
// func(task *Task,index int){
//    //do something
//    //if there is result you need ,you can put it into task.data[index]
//    task.data[index] = result
//}
type TaskItem func(*Task, int) error

func NewTask() *Task {
	t := Task{}
	return t.Init()
}

// New create a new Task schedule,you can use method to operate you taskItem
func (task *Task) Init() *Task {
	task.stop = make(chan struct{}, 1)
	task.start = make(chan struct{}, 1)
	task.pause = make(chan struct{})
	task.wait = make(chan struct{})
	task.works = make([]TaskItem, 0, 4)
	task.data = make(map[int]interface{})
	return task
}

// Add function add a TaskItem into Task
func (task *Task) Add(t TaskItem) *Task {
	if task.works == nil {
		panic("please Init Task")
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
		// wait for start signal
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
				task.pause <- struct{}{}
			default:
			}
			if err := task.Exec(index); err != nil {
				task.data[index] = err
			}
		}
	END:
		task.finish = true
		task.wait <- struct{}{}
		task.status = false
	}()
	return nil
}

func (task *Task) close() {
	close(task.start)
	close(task.pause)
	close(task.stop)
}

// Exec will execute taskItem function in
func (task *Task) Exec(index int) error {
	return task.works[index](task, index)
}

// WaitFinish wait for all taskItem done
func (task *Task) WaitFinish() {
	_ = <-task.wait
	close(task.wait)
	return
}

//当需要再次执行Start时可以重新刷新,在此之前需要执行stop,并等待关闭完成
func (task *Task) Flush() bool {
	if task.finish {
		task.start = make(chan struct{}, 1)
		task.stop = make(chan struct{}, 1)
		task.pause = make(chan struct{})
		task.wait = make(chan struct{})
		task.status = false
		task.finish = false
		return true
	}
	return false
}

func (task *Task) Stop() error {
	var err error
	defer func() {
		if err1 := recover(); err1 != nil {
			err = fmt.Errorf("%v", err1)
		}
	}()
	task.stop <- struct{}{}
	return err
}

func (task *Task) Run() error {
	var err error
	if !task.status {
		err = errors.New("please do Start() function before Run()")
		return err
	}
	task.start <- struct{}{}
	return err
}

// 如果多处触发Pause()可能会导致阻塞，建议 go task.Pause(),利用channel阻塞进行暂停
func (task *Task) Pause() {
	task.pause <- struct{}{}
}

// 如果多处触发Resume()可能会导致阻塞，建议 go task.Resume()
func (task *Task) Resume() {
	<-task.pause
}

//task need user put result of TaskItem into data ,and key is index ,value is result
func (task *Task) Call() (map[int]interface{}, error) {
	return task.data, nil
}
