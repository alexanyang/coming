package coming

import (
	"fmt"
	"testing"
	"time"
)

func TestTask_task(t *testing.T) {
	task := Task{}.New()
	task.Add(func(this *Task, i int) error {
		fmt.Printf("%s 执行第一个任务，序号：%v \n", time.Now(), i)
		this.data[i] = fmt.Sprintf("完成第一个任务，序号：%v \n", i)
		time.Sleep(2 * time.Second)
		return nil
	}).Add(func(this *Task, i int) error {
		fmt.Printf("%s 执行第二个任务，序号：%v \n", time.Now(), i)
		this.data[i] = fmt.Sprintf("完成第二个任务，序号：%v \n", i)
		time.Sleep(2 * time.Second)
		return nil
	}).Add(func(this *Task, i int) error {
		fmt.Printf("%s 执行第三个任务，序号：%v \n", time.Now(), i)
		this.data[i] = fmt.Sprintf("完成第三个任务，序号：%v \n", i)
		time.Sleep(2 * time.Second)
		return nil
	})

	task.Start()
	fmt.Println("任务数量：", len(task.works))
	task.Run()

	time.Sleep(1 * time.Second)
	go task.Pause()
	fmt.Println("暂停中")

	time.Sleep(10 * time.Second)
	task.Resume()
	fmt.Println("重启")
	task.Stop()
	//开始等待结束
	task.WaitFinish()
	fmt.Printf("%s main执行等待完成 \n", time.Now())
	call, _ := task.Call()
	if len(call) > 0 {
		for k, v := range call {
			fmt.Printf("结果 %v:%v \n", k, v)
		}
	}
	time.Sleep(2 * time.Second)
	fmt.Println("main执行完成")

}

func ExampleTask_Add() {
	task := Task{}.New()
	task.Add(func(this *Task, i int) error {
		fmt.Printf("%s 执行第一个任务，Order：%v \n", time.Now(), i)
		this.data[i] = fmt.Sprintf("完成第一个任务，Order：%v ,you can save result to this.data[i] \n", i)
		return nil
	})
}
