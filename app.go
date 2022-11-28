package worker_manager

import (
	slog "github.com/vearne/simplelog"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	wm            *WorkerManager
	exitSigList   []os.Signal
	ignoreSigList []os.Signal
}

func NewApp() *App {
	var app App
	app.wm = NewWorkerManager()
	// default signals
	app.exitSigList = []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
	app.ignoreSigList = make([]os.Signal, 0)
	return &app
}

func (a *App) AddWorker(w Worker) {
	a.wm.AddWorker(w)
}

func (a *App) SetExitSigs(sig ...os.Signal) {
	a.exitSigList = sig
}

func (a *App) SetIgnoreSigs(sig ...os.Signal) {
	a.ignoreSigList = sig
}

func (a *App) Run() {
	if len(a.wm.WorkerSlice) <= 0 {
		panic("The number of workers must be greater than 0!")
	}

	exitSigMap := make(map[os.Signal]struct{})
	for _, sig := range a.exitSigList {
		exitSigMap[sig] = struct{}{}
	}

	ch := make(chan os.Signal, 1)
	// all signal
	sigList := make([]os.Signal, 0)
	sigList = append(sigList, a.exitSigList...)
	sigList = append(sigList, a.ignoreSigList...)
	slog.Debug("sigList:%v", sigList)

	signal.Notify(ch, sigList...)
	go func() {
		for sig := range ch {
			slog.Debug("get sig:%v", sig)
			if _, ok := exitSigMap[sig]; ok {
				close(ch)
				a.wm.Stop()
				break
			}
		}
	}()
	a.wm.Start()
	a.wm.Wait() //nolint: typecheck
}
