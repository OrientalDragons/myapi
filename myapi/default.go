package myapi

import (
	"net/http"
	"regexp"
)

var serverMux = http.NewServeMux()
var handlers []*handler

type Routes interface {
	Get()
	Post()
	SetW(w http.ResponseWriter)
	SetR(r *http.Request)
	Runfilter() bool
}
type StaticRoute interface {
	Routes
	GetRoute() string
}
type handler struct {
	url         string
	handlerFunc func(http.ResponseWriter, *http.Request)
}

func Route(hd string, c Routes) {

	handlers = append(handlers, &(handler{hd, MuxHandler(c)}))
	// serverMux.HandleFunc(hd, MuxHandler(c))
}

func RouteStatic(staticHandle StaticRoute) {
	r := staticHandle.GetRoute()
	if r == "" {
		r = "/static/"
	}
	if r == "/" {
		panic("Route Error: Can not use '/' to be absoluted route")
	}
	serverMux.HandleFunc(r, MuxHandler(staticHandle))
}

func Listen(port string) error {
	serverMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		found := false
		url := r.URL.Path
		if len(handlers) != 0 {
			for _, v := range handlers {
				if match, _ := regexp.MatchString(v.url, url); match {
					found = true
					v.handlerFunc(w, r)
				}
			}
		}
		if found == false {
			http.NotFoundHandler().ServeHTTP(w, r)
			return
		}

	})
	// serverMux.HandleFunc()
	err := http.ListenAndServe(port, serverMux)
	if err != nil {
		return err
	}
	return nil
}

func MuxHandler(route Routes) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		route.SetW(w)
		route.SetR(r)
		if route.Runfilter() {
			if r.Method == "GET" {
				route.Get()
			}
			if r.Method == "POST" {
				route.Post()
			}
			return
		}
		return
	}
}
