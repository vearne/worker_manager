package main 

import (
	"os"
	"syscall"
	"os/signal"
	"worker_manager/libs"
)

func main() {
	// 1. 初始化各种worker
	wm := prepareAllWorker()

	// 2. start
	wm.Start()

	// 3. register grace exit
	GracefulExit(wm)
	
	// 4. block and wait
	wm.Wait()
}



func GracefulExit(wm * libs.WorkerManager){
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM)
	switch <-ch {
	case syscall.SIGTERM:
		wm.Stop()
		break
	}
}


func prepareAllWorker() * libs.WorkerManager{
	wm := libs.NewWorkerManager()


	// workerA
	WorkerACount := 2
	for i:=0; i< WorkerACount; i++{
		wm.AddWorker(&libs.WorkerA{Name:"WorkerA"})
	}
	// workerB
	WorkerBCount := 3
	for i:=0; i< WorkerBCount; i++{
		wm.AddWorker(&libs.WorkerB{Name:"WorkerB"})
	}
	// workerC
	WorkerCCount := 4
	for i:=0; i< WorkerCCount; i++{
		wm.AddWorker(&libs.WorkerC{Name:"WorkerC"})
	}
	return wm
}

