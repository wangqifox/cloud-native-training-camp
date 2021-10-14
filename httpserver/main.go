package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 80, "服务端口")
	log.Printf("serving at %s\n", ":"+strconv.Itoa(port))

	mux := http.NewServeMux()
	mux.Handle("/", middleware(http.HandlerFunc(rootHandler)))
	mux.Handle("/healthz", middleware(http.HandlerFunc(healthz)))
	mux.Handle("/err", middleware(http.HandlerFunc(err)))
	err := http.ListenAndServe(":"+strconv.Itoa(port), mux)
	if err != nil {
		log.Fatal(err)
	}
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

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func err(w http.ResponseWriter, r *http.Request) {
	if true {
		panic("this is an error")
	}
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
