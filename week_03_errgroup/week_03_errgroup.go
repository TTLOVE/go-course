package week_03_errgroup

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type HelloHandler struct{}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("hello go"))
}

func HttpServer() {
	// 生成errgroup
	g, cxt := errgroup.WithContext(context.Background())

	// 生成处理请求的handler
	mux := http.NewServeMux()

	mux.Handle("/", &HelloHandler{})
	mux.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	// 模拟页面申请退出
	reqOut := make(chan struct{})
	mux.HandleFunc("/out", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "shutting down")
		reqOut <- struct{}{}
	})

	// 设置服务信息
	server := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}

	g.Go(func() error {
		select {
		case <-reqOut:
			log.Println("shutdown from request")
		case <-cxt.Done():
			log.Println("shutdown from errgroup")
		}

		return server.Shutdown(cxt)
	})

	// 发起http服务
	g.Go(func() error {
		log.Println("server start")

		err := server.ListenAndServe()
		if err != nil {
			return errors.Wrap(err, "serve fail")
		}

		return nil
	})

	done := make(chan os.Signal)
	// 创建系统信号接收器
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// 监听linux signal退出通道
	g.Go(func() error {
		select {
		case <-done:
			log.Println("shutdown from signal")
		case <-cxt.Done():
			log.Println("shutdown from errgroup")
		}

		return server.Shutdown(cxt)
	})

	if err := g.Wait(); err != nil {
		log.Println(err)
	}
}
