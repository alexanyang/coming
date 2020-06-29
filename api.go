package coming

//访问Stop()接口失败,没有成功关闭
//
type StopRefuseError struct{}

func (StopRefuseError) Error() string {
	return "Process refused to stop."
}

type TaskInterface interface {
	//正常启动目标返回nil,否则返回错误信息,程序执行期间阻塞
	Run() error
	//实现改接口,以调用Stop()进行停止当前携程,本质是在执行过程中,通过通道传递信号,使得协程停止
	//执行stop方法,当程序正常停止时返回nil,当程序停止被拒绝,返回StopRefuseError,当程序停止
	//成功,但是报错时,将返回对应的错误信息
	Stop() error
	//如果完成了Start()操作则返回nil,否则返回err信息
	Start() error
	//Pause to stop next task , resume to allow run
	Pause()
	//Pause to stop next task , resume to allow run
	Resume()
}

type BaseTask struct{}

//正常启动目标返回nil,否则返回错误信息,程序执行期间阻塞
func (this BaseTask) Run() error {}

//实现改接口,以调用Stop()进行停止当前携程,本质是在执行过程中,通过通道传递信号,使得协程停止
//执行stop方法,当程序正常停止时返回nil,当程序停止被拒绝,返回StopRefuseError,当程序停止
//成功,但是报错时,将返回对应的错误信息
func (this BaseTask) Stop() error {}

//如果完成了Start()操作则返回nil,否则返回err信息
func (this BaseTask) Start() error {}

//Pause to stop next task , resume to allow run
func (this BaseTask) Pause() {}

//Pause to stop next task , resume to allow run
func (this BaseTask) Resume() {}
