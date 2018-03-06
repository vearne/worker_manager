package libs

import (
	"sync"
	"fmt"
)

type Worker interface{
	Start()
	Stop()
}

type  WorkerManager struct {
	sync.WaitGroup
	// 保存所有worker
	WorkerSlice []Worker
}

func NewWorkerManager() *WorkerManager{
	workerManager := WorkerManager{}
	workerManager.WorkerSlice = make([]Worker, 0, 10)
	return &workerManager
}

func (wm *WorkerManager) AddWorker(w Worker){
	wm.WorkerSlice = append(wm.WorkerSlice, w)
	fmt.Println("size", len(wm.WorkerSlice), wm.WorkerSlice)
}

func (wm *WorkerManager) Start(){
	wm.Add(len(wm.WorkerSlice))
	for _, worker := range wm.WorkerSlice{
		go func(w Worker){
			fmt.Printf("start worker:%v\n", w)
			w.Start()
		}(worker)
	}
}

func (wm *WorkerManager) Stop(){

	for _, worker := range wm.WorkerSlice{
		go func(w Worker){
			fmt.Printf("stop worker:%v\n", w)
			w.Stop()
			wm.Done()
		}(worker)
	}
}



