package libs

import (
	"fmt"
	"time"
)

type WorkerA struct{
	Name string
}

func (worker *WorkerA) Start(){
	fmt.Println("WorkerA start")
}
func (worker *WorkerA) Stop(){
	time.Sleep(time.Second)
	fmt.Println("WorkerA stop")
}

type WorkerB struct{
	Name string
}

func (worker *WorkerB) Start(){
	fmt.Println("WorkerB start")
}
func (worker *WorkerB) Stop(){
	time.Sleep(time.Second)
	fmt.Println("WorkerB stop")
}

type WorkerC struct{
	Name string
}

func (worker *WorkerC) Start(){
	fmt.Println("WorkerC start")
}
func (worker *WorkerC) Stop(){
	time.Sleep(time.Second)
	fmt.Println("WorkerC stop")
}