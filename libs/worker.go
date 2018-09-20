package libs

import (
	"fmt"
	"time"
)

type WorkerA struct{
	RunningFlag bool  // 是否运行 true:运行 false:停止
	ExitedFlag bool  //  已经退出的标识
	Name string
}

func NewWorkerA(name string) *WorkerA{
	return &WorkerA{RunningFlag:true, ExitedFlag:false, Name:name}
}

func (worker *WorkerA) Start(){
	fmt.Println("WorkerA start")
	for worker.RunningFlag{
		time.Sleep(time.Second)
		fmt.Println("WorkerA do something ...")
	}
	worker.ExitedFlag = true
}
func (worker *WorkerA) Stop(){
	worker.RunningFlag = false
	for !worker.ExitedFlag{
		time.Sleep(time.Millisecond * 50)
	}
	fmt.Println("WorkerA stop")
}

type WorkerB struct{
	RunningFlag bool  // 是否运行 true:运行 false:停止
	ExitedFlag bool  //  已经退出的标识
	Name string
}

func NewWorkerB(name string) *WorkerB{
	return &WorkerB{RunningFlag:true, ExitedFlag:false, Name:name}
}

func (worker *WorkerB) Start(){
	fmt.Println("WorkerB start")
	for worker.RunningFlag{
		time.Sleep(time.Second * 3)
		fmt.Println("WorkerB do something ...")
	}
	worker.ExitedFlag = true
}
func (worker *WorkerB) Stop(){
	worker.RunningFlag = false
	for !worker.ExitedFlag{
		time.Sleep(time.Millisecond * 50)
	}
	fmt.Println("WorkerB stop")
}

