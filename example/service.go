package main

import (
	"context"
	"github.com/gin-gonic/gin"
	wm "github.com/vearne/worker_manager"
	"log"
	"net/http"
	"time"
)

func main() {
	app := wm.NewApp()
	// add 2 load worker
	app.AddWorker(NewLoadWorker())
	app.AddWorker(NewLoadWorker())
	// add 1 web worker
	app.AddWorker(NewWebServer())
	// If not set, the default value will be used
	//app.SetSigs(syscall.SIGTERM, syscall.SIGQUIT)
	app.Run()
}

// some worker

type LoadWorker struct {
	RunningFlag *wm.BoolFlag
	ExitedFlag  *wm.BoolFlag
	ExitChan    chan struct{}
}

func NewLoadWorker() *LoadWorker {
	worker := &LoadWorker{}
	worker.RunningFlag = wm.NewBoolFlag()
	worker.ExitedFlag = wm.NewBoolFlag()
	wm.SetTrue(worker.RunningFlag)
	wm.SetTrue(worker.ExitedFlag)
	worker.ExitChan = make(chan struct{})
	return worker
}

func (worker *LoadWorker) Start() {
	log.Println("[start]LoadWorker")
	for wm.IsTrue(worker.RunningFlag) {
		select {
		case <-time.After(1 * time.Minute):
			//do some thing
			log.Println("LoadWorker do something")
			time.Sleep(time.Second * 3)
		case <-worker.ExitChan:
			log.Println("LoadWorker execute exit logic")
		}
	}
	wm.SetTrue(worker.ExitedFlag)
}

func (worker *LoadWorker) Stop() {
	log.Println("LoadWorker exit...")
	wm.SetFalse(worker.RunningFlag)
	close(worker.ExitChan)
	for !wm.IsTrue(worker.ExitedFlag) {
		time.Sleep(50 * time.Millisecond)
	}
	log.Println("[end]LoadWorker")
}

type WebServer struct {
	Server *http.Server
}

func NewWebServer() *WebServer {
	return &WebServer{}
}

func (worker *WebServer) Start() {
	log.Println("[start]WebServer")

	ginHandler := gin.Default()
	ginHandler.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/plain", []byte("hello world!"))
	})
	worker.Server = &http.Server{
		Addr:           ":8080",
		Handler:        ginHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	worker.Server.ListenAndServe()
}

func (worker *WebServer) Stop() {
	log.Println("WebServer exit...")
	cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// gracefull exit web server
	err := worker.Server.Shutdown(cxt)
	if err != nil {
		log.Printf("shutdown error, %v", err)
	}
	log.Println("[end]WebServer exit")
}
