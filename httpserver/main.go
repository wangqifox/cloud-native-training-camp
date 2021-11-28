package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

const stopTimeout = time.Second * 10

func main() {
	var port int
	flag.IntVar(&port, "port", 80, "服务端口")
	log.Printf("serving at %s\n", ":"+strconv.Itoa(port))

	sigs := make(chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(sigs, os.Interrupt)

	mux := http.NewServeMux()
	mux.Handle("/", middleware(http.HandlerFunc(rootHandler)))
	mux.Handle("/longTimeRequest", middleware(http.HandlerFunc(longTimeRequest)))
	mux.Handle("/healthz", middleware(http.HandlerFunc(healthz)))
	mux.Handle("/err", middleware(http.HandlerFunc(err)))

	var srv = &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}

	// 优雅终止
	go func() {
		<-sigs
		ctx, cancel := context.WithTimeout(context.Background(), stopTimeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("Server shutdown with error: ", err)
		}
		close(done)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
	<-done
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	version := os.Getenv("VERSION")
	if len(version) > 0 {
		w.Header().Add("version", version)
	}
	for h := range r.Header {
		w.Header().Add(h, r.Header.Get(h))
	}

	_, err := w.Write([]byte("ok"))
	if err != nil {
		log.Fatal(err)
	}
}

// 健康检查
func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

// 异常处理模拟
func err(w http.ResponseWriter, r *http.Request) {
	if true {
		panic("this is an error")
	}
	w.WriteHeader(200)
}

// 模拟长时间处理的请求，测试优雅终止功能
func longTimeRequest(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 5)
	w.WriteHeader(200)
}

type responseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// 打印请求ip和错误信息
func middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{
			ResponseWriter: w,
			StatusCode:     200,
		}

		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				fmt.Println("http code:", 500)
				w.WriteHeader(500)
			} else {
				fmt.Println("http code:", rw.StatusCode)
			}
		}()

		if ip, _, e := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); e == nil {
			fmt.Println("client ip:", ip)
		}

		handler.ServeHTTP(rw, r)
	})
}
