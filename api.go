package coming

import "errors"

// TaskInterface contains methods of a task need implement
type TaskInterface interface {
	// Run 正常启动目标返回nil,否则返回错误信息,程序执行期间阻塞
	Run() error
	// Stop can stop task next running,and the running task will run before finished
	// 实现改接口,以调用Stop()进行停止当前携程,本质是在执行过程中,通过通道传递信号,使得协程停止
	// 执行stop方法,当程序正常停止时返回nil,当程序停止被拒绝,返回StopRefuseError,当程序停止
	// 成功,但是报错时,将返回对应的错误信息
	Stop() error
	// Start to run all tasks from header
	//如果完成了Start()操作则返回nil,否则返回err信息
	Start() error
	// Pause to stop next task , resume to allow run
	Pause()
	// Resume to allow run ,Pause to stop next task
	Resume()
}

// BaseTask is an default task with method always return message not nil.
// it will make you only focus on the method you will use
type BaseTask struct{}

// Run 正常启动目标返回nil,否则返回错误信息,程序执行期间阻塞
func (bt BaseTask) Run() error {
	return errors.New("you need override Run() method")
}

// Stop 实现改接口,以调用Stop()进行停止当前携程,本质是在执行过程中,通过通道传递信号,使得协程停止
// 执行stop方法,当程序正常停止时返回nil,当程序停止被拒绝,返回StopRefuseError,当程序停止
// 成功,但是报错时,将返回对应的错误信息
func (bt BaseTask) Stop() error {
	return errors.New("you need override Stop() method")
}

// Start 如果完成了Start()操作则返回nil,否则返回err信息
func (bt BaseTask) Start() error {
	return errors.New("you need override Start() method")
}

// Pause to stop next task
func (bt BaseTask) Pause() {
	panic("you need override Pause() and Resume() method to use this function")
}

// Resume to allow run task
func (bt BaseTask) Resume() {
	panic("you need override Pause() and Resume() method to use this function")
}
