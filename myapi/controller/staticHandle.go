package controller

import (
	"net/http"
	"os"
	"path"
)

type StaticHandler struct {
	rout string
	Controller
}

func (c *StaticHandler) Runfilter() bool {

	wd, osErr := os.Getwd()
	if osErr != nil {
		panic(osErr)
	}
	dirPath := path.Join(wd, c.rout)
	Handler := http.StripPrefix(c.rout, http.FileServer(http.Dir(dirPath)))
	Handler.ServeHTTP(c.W, c.R)
	return false
}
func (c *StaticHandler) GetRoute() string {
	return c.rout
}
