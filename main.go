package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Router struct {
	http.Handler
	mux *http.ServeMux
}

func NewRouter() *Router {
	mux := http.NewServeMux()
	return &Router{
		Handler: mux,
		mux:     mux,
	}
}

func (r *Router) Handle(pattern string, handler http.Handler) {
	r.mux.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		// 此处可以增加扩展逻辑

		handler.ServeHTTP(w, req)
	})
}

func (r *Router) HandleFunc(pattern string, handler http.HandlerFunc) {
	r.Handle(pattern, handler)
}

type Middleware func(http.Handler) http.Handler

func Chain(middlewares ...Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		length := len(middlewares)
		for i := range middlewares {
			h = middlewares[length-1-i](h)
		}
		return h
	}
}

func Recover() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Println("recover from panic!")
				}
			}()

			log.Println("middleware recover: before")

			h.ServeHTTP(w, r)

			log.Println("middleware recover: after")
		})
	}
}

func Timeout(d time.Duration) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("middleware timeout: before")

			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()
			h.ServeHTTP(w, r.WithContext(ctx))

			log.Println("middleware timeout: after")
		})
	}
}

func main() {
	router := NewRouter()

	router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.ServeContent(w, r, "", time.Time{}, strings.NewReader(`User-agent: *
Disallow: /search?*
`))
	})

	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		log.Println("path:", r.URL.Path)
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("path:", r.URL.Path)
		fmt.Fprint(w, "Hello pkgsite!")
	})

	mw := Chain(
		Recover(),
		Timeout(30*time.Second),
	)

	addr := ":2020"
	log.Println("Listening on addr ", addr)
	log.Fatal(http.ListenAndServe(addr, mw(router)))
}
