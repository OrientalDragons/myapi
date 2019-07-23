package controller

import (
	"net/http"
)

//Controller 是 http 的 controller 容器
type Controller struct {
	W http.ResponseWriter
	R *http.Request
}

// Get adds a request function to handle GET request.
func (c *Controller) Get() {
	http.Error(c.W, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// Post adds a request function to handle POST request.
func (c *Controller) Post() {
	http.Error(c.W, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func (c *Controller) SetW(w http.ResponseWriter) {
	c.W = w
}
func (c *Controller) SetR(r *http.Request) {
	c.R = r
}

func (c *Controller) Runfilter() bool {
	return true
}
