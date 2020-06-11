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
type Stopable interface {
	Stop() error
}

type Startable interface {
	Start() error
}

type Runable interface {
	Run() ([]byte, error)
}
