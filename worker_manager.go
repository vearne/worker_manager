package worker_manager

import (
	"expvar"
	"fmt"
	"runtime/debug"
	"sync"
)

type Worker interface {
	Start()
	Stop()
}

type WorkerManager struct {
	sync.WaitGroup
	aliveWorkerNum *expvar.Int
	WorkerSlice    []Worker
}

func NewWorkerManager() *WorkerManager {
	workerManager := WorkerManager{}
	workerManager.WorkerSlice = make([]Worker, 0, 10)
	workerManager.aliveWorkerNum = expvar.NewInt("aliveWorkerNum")
	return &workerManager
}

func (wm *WorkerManager) AddWorker(w Worker) {
	wm.WorkerSlice = append(wm.WorkerSlice, w)
}

func (wm *WorkerManager) Start() {
	wm.Add(len(wm.WorkerSlice)) //nolint: typecheck
	wm.aliveWorkerNum.Set(int64(len(wm.WorkerSlice)))
	for _, worker := range wm.WorkerSlice {
		go func(w Worker) {
			defer func() {
				r := recover()
				if r != nil {
					fmt.Printf("WorkerManager error, recover:%v, stack:%v\n",
						r, string(debug.Stack()))
					wm.Done() //nolint: typecheck
					wm.aliveWorkerNum.Add(-1)
				}
			}()
			w.Start()
		}(worker)
	}
}

func (wm *WorkerManager) Stop() {
	for _, worker := range wm.WorkerSlice {
		go func(w Worker) {
			defer func() {
				r := recover()
				if r != nil {
					fmt.Printf("WorkerManager error, recover:%v, stack:%v\n",
						r, string(debug.Stack()))
				}
			}()

			w.Stop()
			wm.Done() //nolint: typecheck
			wm.aliveWorkerNum.Add(-1)
		}(worker)
	}
}
