package task

/**
 * @author anyang
 * Email: 1300378587@qq.com
 * Created Date:2020-06-11 09:20
 */
type Executable interface {
	Start() error
	Stop() error
}

//实现改接口,以调用Stop()进行停止当前携程,本质是在执行过程中,通过通道传递信号,使得协程停止
//执行stop方法,当程序正常停止时返回nil,当程序停止被拒绝,返回StopRefuseError,当程序停止
//成功,但是报错时,将返回对应的错误信息
type Stopable interface {
	Stop() error
}

//如果完成了Start()操作则返回nil,否则返回err信息
type Startable interface {
	Start() error
}

type Pauable interface {
	Pause()
	Resume()
}

//正常启动目标返回nil,否则返回错误信息,程序执行期间阻塞
type Runable interface {
	Run() error
}

//正常启动目标,程序执行期间阻塞,并在目标执行完了后进行结果返回,如果提前返回,则需要包含错误信息
type Callable interface {
	Call() (interface{}, error)
}

//访问Stop()接口失败,没有成功关闭
type StopRefuseError struct{}

func (StopRefuseError) Error() string {
	return "Process refused to stop."
}
