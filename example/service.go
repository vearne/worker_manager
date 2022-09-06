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
	RunningFlag *wm.AtomicBool
	ExitedFlag  *wm.AtomicBool
	ExitChan    chan struct{}
}

func NewLoadWorker() *LoadWorker {
	worker := &LoadWorker{}
	worker.RunningFlag = wm.NewAtomicBool(true)
	worker.ExitedFlag = wm.NewAtomicBool(true)
	worker.ExitChan = make(chan struct{})
	return worker
}

func (worker *LoadWorker) Start() {
	log.Println("[start]LoadWorker")
	for worker.RunningFlag.IsTrue() {
		select {
		case <-time.After(1 * time.Minute):
			//do some thing
			log.Println("LoadWorker do something")
			time.Sleep(time.Second * 3)
		case <-worker.ExitChan:
			log.Println("LoadWorker execute exit logic")
		}
	}
	worker.ExitedFlag.Set(true)
}

func (worker *LoadWorker) Stop() {
	log.Println("LoadWorker exit...")
	worker.RunningFlag.Set(false)
	close(worker.ExitChan)
	for !worker.ExitedFlag.IsTrue() {
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
