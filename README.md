# coming
<!--<img align="right" src="https://raw.githubusercontent.com/alexanyang/coming/master/logo.jpg">-->
<!--<a href="https://circleci.com/gh/alexanyang/coming/tree/dev"><img src="https://img.shields.io/circleci/project/alexanyang/coming/dev.svg" alt="Build Status"></a>-->

[//]: # ([![CircleCI Status]&#40;https://circleci.com/gh/alexanyang/coming.svg?style=shield&#41;]&#40;https://circleci.com/gh/alexanyang/coming&#41;)
[//]: # (![Appveyor]&#40;https://ci.appveyor.com/api/projects/status/github/alexanyang/coming?branch=master&svg=true&#41;)
[//]: # ([![codecov]&#40;https://codecov.io/gh/alexanyang/coming/branch/master/graph/badge.svg&#41;]&#40;https://codecov.io/gh/alexanyang/coming&#41;)
[//]: # ([![Build Status]&#40;https://travis-ci.org/alexanyang/coming.svg&#41;]&#40;https://travis-ci.org/alexanyang/coming&#41;)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexanyang/coming)](https://goreportcard.com/report/github.com/alexanyang/coming)
[![GoDoc](https://godoc.org/github.com/alexanyang/coming?status.svg)](https://godoc.org/github.com/alexanyang/coming)
[![GitHub release](https://img.shields.io/github/release/alexanyang/coming.svg)](https://github.com/alexanyang/coming/releases/latest)
[![Join the chat at https://gitter.im/alexanyang/coming](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/alexanyang/coming?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
<!-- [![Release](https://github-release-version.herokuapp.com/github/alexanyang/coming/release.svg?style=flat)](https://github.com/alexanyang/coming/releases/latest) -->
<!--<a href="https://github.com/alexanyang/coming/releases"><img src="https://img.shields.io/badge/%20version%20-%206.0.0%20-blue.svg?style=flat-square" alt="Releases"></a>-->


"# coming" 
//coming 是一个原生的任务调度系统，目前给出了定时任务和计划任务的解决方案

*已实现功能*:
1.添加定时任务 

2.构造顺序工作链

3.实现暂停/重启功能

4.可循环执行功能，已实现，需要手动

```go
package main

import (
	"fmt"
	"time"
	
	"github.com/alexanyang/coming"
)

func main() {
	task := coming.NewTask()
	task.Add(func(this *coming.Task, i int) error {
		fmt.Printf("%s 执行第一个任务，序号：%v \n", time.Now(), i)
		this.PushData(fmt.Sprintf("第一个任务结果，序号：%v", i))
		time.Sleep(2 * time.Second)
		return nil
	}).Add(func(this *coming.Task, i int) error {
		fmt.Printf("%s 执行第二个任务，序号：%v \n", time.Now(), i)
		this.PushData(fmt.Sprintf("第二个任务结果，序号：%v", i))
		time.Sleep(2 * time.Second)
		return nil
	}).Add(func(this *coming.Task, i int) error {
		fmt.Printf("%s 执行第三个任务，序号：%v \n", time.Now(), i)
		this.PushData(fmt.Sprintf("第三个任务结果，序号：%v", i))
		time.Sleep(2 * time.Second)
		return nil
	})
	// 检查操作是否成功
	if err := task.Error(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s: 任务Start\n", time.Now())
	if err := task.Start(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s: 任务数量：%v, and run\n", time.Now(), task.Len())
	if err := task.Run(); err != nil {
		fmt.Println(err)
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

```

运行定时任务
```go 
func TestNewSchedule(t *testing.T) {
    //ANY 分别对应时分秒,hour,minute,second
	NewDailySchedule(ANY, ANY, ANY,"每秒任务", func(t time.Time) {
    		fmt.Printf("现在时间是%s,定时执行中\n", t)
    	})  
    //延时关闭才能看到协程输出
    time.Sleep(5 * time.Second)
}

```

your star is best support
