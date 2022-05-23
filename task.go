package coming

import (
	"errors"
	"fmt"
	"sync/atomic"
)

// Task struct should be used for task that will be running for long time
// so that you can use other methods like Pause() before finish
type Task struct {
	BaseTask
	status      int32
	finish      bool
	stop, start chan struct{}
	pause, wait chan struct{}
	works       []TaskItem
	data        []interface{}
	index       int
	err         error
}

const (
	TaskStateDefault int32 = iota
	TaskStateInit
	TaskStateReady
	TaskStateRunning
	TaskStatePause
	TaskStateFinish
)

var (
	ErrorRunBeforeStart = errors.New("call  Run() function before Start()")
	ErrorNotInit        = errors.New("task not init")
	ErrorStarted        = errors.New("task had started ,not support add")
	ErrorPushData       = errors.New("push data can only call in TaskItem")
	ErrorCacheUsed      = errors.New("this cache had been used")
)

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

// Init New create a new Task schedule,you can use method to operate you taskItem
func (task *Task) Init() *Task {
	task.stop = make(chan struct{}, 1)
	task.start = make(chan struct{}, 1)
	task.pause = make(chan struct{})
	task.wait = make(chan struct{})
	if len(task.works) == 0 {
		task.works = make([]TaskItem, 0, 4)
	}
	task.status = TaskStateInit
	return task
}

// Add function add a TaskItem into Task.
// if we do not pass check ,just record error and return this task point.
// You can also get error by using method Error of it.
func (task *Task) Add(t TaskItem) *Task {
	if err := task.check(); err != nil {
		task.err = err
	} else {
		task.works = append(task.works, t)
	}
	return task
}

// work wait a start signal and execute function in works of it
// util get a signal of stop or pause .
// the pause function worked by channel receive waiting .
func (task *Task) work() {
	task.data = make([]interface{}, len(task.works))
	// wait for start signal
	<-task.start
	for index := range task.works {
		select {
		case <-task.stop:
			//关闭其他通道
			task.close()
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
	atomic.StoreInt32(&task.status, TaskStateFinish)
	task.wait <- struct{}{}
}

// Start function will start a waiting for run TaskItem of slice
// It means that user had added all task ready.
func (task *Task) Start() error {
	if atomic.CompareAndSwapInt32(&task.status, TaskStateInit, TaskStateReady) {
		go task.work()
	}
	return nil
}

func (task *Task) close() {
	close(task.start)
	close(task.pause)
	close(task.stop)
}

// Exec will execute taskItem function in
func (task *Task) Exec(index int) error {
	// keep task index is position of running item,
	// that make sure push data into right position,
	// and reset index to -1 after function done
	task.index = index
	defer func() {
		task.index = -1
	}()
	return task.works[index](task, index)
}

// WaitFinish wait for all taskItem done
func (task *Task) WaitFinish() {
	_ = <-task.wait
	close(task.wait)
	return
}

// Flush 当需要再次执行Start时可以重新刷新,在此之前需要执行stop,并等待关闭完成
func (task *Task) Flush() bool {
	if atomic.CompareAndSwapInt32(&task.status, TaskStateFinish, TaskStateInit) {
		task.Init()
		return true
	}
	return false
}

func (task *Task) Stop() (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			err = fmt.Errorf("%v", err1)
		}
	}()
	task.stop <- struct{}{}
	return
}

func (task *Task) Run() error {
	if !atomic.CompareAndSwapInt32(&task.status, TaskStateReady, TaskStateRunning) {
		return ErrorRunBeforeStart
	}
	task.start <- struct{}{}
	return nil
}

// Pause send a signal to tell work to wait,
// use atomic to avoid multi signal makes channel block.
// This causes the status update to precede the actual operation
func (task *Task) Pause() {
	if atomic.CompareAndSwapInt32(&task.status, TaskStateRunning, TaskStatePause) {
		task.pause <- struct{}{}
	}
}

// Resume receive pause channel to free it from block to used
// atomic can avoid  multi signal makes function block.
// This causes the status update to precede the actual operation
func (task *Task) Resume() {
	if atomic.CompareAndSwapInt32(&task.status, TaskStatePause, TaskStateRunning) {
		<-task.pause
	}
}

// Call task need user put result of TaskItem into data ,and key is index ,value is result
func (task *Task) Call() ([]interface{}, error) {
	return task.data, nil
}

func (task *Task) check() error {
	if task.status == TaskStateDefault || task.works == nil {
		return ErrorNotInit
	}
	if task.status >= TaskStateReady {
		return ErrorStarted
	}
	return nil
}

func (task *Task) Error() error {
	return task.err
}

func (task *Task) PushData(data interface{}) (err error) {
	if task.index >= 0 {
		return task.pushData(task.index, data)
	}
	return ErrorPushData
}

func (task *Task) pushData(index int, data interface{}) (err error) {
	if task.data[index] != nil {
		return ErrorCacheUsed
	}
	task.data[index] = data
	return nil
}

// PushDataForce can push result of taskItem into cache.
// It can avoid index out of range without user index
func (task *Task) PushDataForce(data interface{}) (err error) {
	if task.index >= 0 {
		return task.pushDataForce(task.index, data)
	}
	return ErrorPushData
}

// PushDataForce can push result of taskItem into cache.
// It can avoid index out of range without user index
func (task *Task) pushDataForce(index int, data interface{}) (err error) {
	task.data[index] = data
	return nil
}

func (task *Task) Len() int {
	return len(task.works)
}
