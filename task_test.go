package coming_test

import (
	"fmt"
	. "github.com/alexanyang/coming"
	"testing"
	"time"
)

func TestTask_task(t *testing.T) {
	task := NewTask()
	task.Add(func(this *Task, i int) error {
		fmt.Printf("%s 执行第一个任务，序号：%v \n", time.Now(), i)
		this.PushData(fmt.Sprintf("第一个任务结果，序号：%v", i))
		time.Sleep(2 * time.Second)
		return nil
	}).Add(func(this *Task, i int) error {
		fmt.Printf("%s 执行第二个任务，序号：%v \n", time.Now(), i)
		this.PushData(fmt.Sprintf("第二个任务结果，序号：%v", i))
		time.Sleep(2 * time.Second)
		return nil
	}).Add(func(this *Task, i int) error {
		fmt.Printf("%s 执行第三个任务，序号：%v \n", time.Now(), i)
		this.PushData(fmt.Sprintf("第三个任务结果，序号：%v", i))
		time.Sleep(2 * time.Second)
		return nil
	})
	// 检查操作是否成功
	if err := task.Error(); err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%s: 任务Start\n", time.Now())
	if err := task.Start(); err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%s: 任务数量：%v, and run\n", time.Now(), task.Len())
	if err := task.Run(); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(1 * time.Second)
	fmt.Printf("%s: 暂停中\n", time.Now())
	task.Pause()

	go func() {
		time.Sleep(5 * time.Second)
		fmt.Printf("%s: resume\n", time.Now())
		task.Resume()
	}()
	//fmt.Println("stop")
	//task.Stop()
	//开始等待结束
	fmt.Printf("%s: main等待执行完成 \n", time.Now())
	task.WaitFinish()
	fmt.Printf("%s: main执行完成\n", time.Now())
	fmt.Println("获取结果回调:")
	call, _ := task.Call()
	if len(call) > 0 {
		for k, v := range call {
			fmt.Printf("结果 %v:%v \n", k, v)
		}
	}
	time.Sleep(2 * time.Second)
	fmt.Println("结果输出完成")
}

func ExampleTask_Add() {
	task := NewTask()
	task.Add(func(this *Task, i int) error {
		fmt.Printf("%s 执行第一个任务，Order：%v \n", time.Now(), i)
		this.PushData(fmt.Sprintf("完成第一个任务，Order：%v ,you can save result to this.data[i] \n", i))
		return nil
	})
}
