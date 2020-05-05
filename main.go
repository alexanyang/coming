package main

import (
	"coming/task"
	"fmt"
)

//func main() {
//
//	start := time.Now()
//	//ss
//
//	ch := make(chan string, 1)
//	for _, url := range os.Args[1:] {
//		go fetch(url, ch) // start a goroutine
//	}
//	for range os.Args[1:] {
//		fmt.Println(<-ch) // receive from channel ch
//	}
//	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
//}
//func fetch(url string, ch chan<- string) {
//	start := time.Now()
//	resp, err := http.Get(url)
//	if err != nil {
//		ch <- fmt.Sprint(err) // send to channel ch
//		return
//	}
//	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
//	resp.Body.Close() // don't leak resources
//	if err != nil {
//		ch <- fmt.Sprintf("while reading %s: %v", url, err)
//		return
//	}
//	secs := time.Since(start).Seconds()
//	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
//}

func main() {
	newTask := task.NewTask(TaskTest{})
	err := newTask.Start()
	fmt.Println("开启任务")
	if err != nil {
		fmt.Println("开启任务失败")
	}
	for newTask.IsFinished() {
		fmt.Printf("task result %s \n", newTask.GetResult())
		fmt.Println("task finished")
		//time.Sleep(2 * time.Second)
		break
	}
}

type TaskTest struct {
	s string
}

func (tt TaskTest) Run() ([]byte, error) {
	return []byte("hello task"), nil
}

//
//func Sum(opts ...interface{}) (r []byte,err error){
//	param := opts[0].([]interface{})
//	sum := param[0].(int) + param[1].(int)
//	bytes ,err:= json.Marshal(sum)
//
//	return bytes,err
//}

//go get github.com/go-task/task
//第一步;
//// 实现这个接口，或者传入一个func（） ，我们能狗通过这个接口的实现，来进行协程调度
//// 定义这个任务调用器的功能，
////第二步： 我们主题结构体，比如;任务，定时器，执行器
//// 抽象化任务调度结构体的属性
////第三步：  改造这个结构，实现协程通信，要求能狗知道协程是进行中还是完成了，
////实现接口
////第四步：  我们就添加中断机制，实现协程中断，，，（最好暂停）
////优化通信
