package main

import (
	"context"
	log "github.com/cihub/seelog"
	"net/http"
	"os"
	"os/signal"
	"relaper.com/kubemanage/httpserver"
	"relaper.com/kubemanage/inital"
	"relaper.com/kubemanage/task"
	"relaper.com/kubemanage/utils/tools"
	"syscall"
	"time"
)

// @title K8sManage API
// @version 1.0
// @description This is a K8sManage server

// @contact.name API K8sManage
// @contact.url http://127.0.0.1:7475/docs/index.html
// @contact.email lgy10271416@gmail.com

// @license.name K8sManage API 1.0
// @license.url http://127.0.0.1:7475/docs/index.html
// @license.swagger 2.0

// @host 127.0.0.1:7474
func main() {
	opts := inital.InitOptions(os.Args[1:])
	inital.InitRoot(opts)

	logger, err := log.LoggerFromConfigAsFile(opts.LogPath)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer log.Flush()
	err = log.ReplaceLogger(logger)
	if err != nil {
		log.Error(err.Error())
		return
	}

	err = inital.InitCfg(opts)
	if err != nil {
		log.Error(err.Error())
		return
	}
	go task.ClusterTask()
	var wg tools.WaitGroupWrapper
	var server *http.Server
	done := make(chan struct{})
	wg.Wrap(func() {
		server = &http.Server{Addr: opts.Address, Handler: httpserver.Mux(opts.Dev, opts.Address)}
		if err := server.ListenAndServe(); err != nil {
			log.Error("action", "init http server fatal", "err", err)
		}
		done <- struct{}{}
	})

	signalChan := make(chan os.Signal, 1)
	go func() {
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan
		log.Info("action", "shutdown http server")
		if server != nil {
			server.Shutdown(context.Background())
		}
		select {
		case <-done:
		case <-time.After(3 * time.Second):
			log.Error("action", "server shutdown timeout")
			os.Exit(1)
		}
	}()

	log.Info("action", "success init http server", "addr", opts.Address)
	wg.Wait()
	log.Info("action", "graceful exit")

	inital.CloseCache()
}
