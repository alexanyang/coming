"# coming" 
//coming 是一个原生的任务调度系统，目前给出了定时任务和顺序任务的解决方案

*第一步*;
    接口定义，已完成
*第二步*： 
    实现功能，已完成
*第三步*：  
    优化结构，持续中

*已实现功能*:
1.添加定时任务 

2.构造顺序工作链

3.实现暂停/重启功能

4.可循环执行功能，已实现，需要手动

```go
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
```

```go 
func TestNewSchedule(t *testing.T) {

	NewDailySchedule(ANY, ANY, ANY, func(t time.Time) {
    		fmt.Printf("现在时间是%s,定时执行中\n", t)
    	})
    time.Sleep(5 * time.Second)
}

```

如果你觉得这里的代码又帮助到你,可以star表示支持,谢谢!!
