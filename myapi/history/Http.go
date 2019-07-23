package myapi

import (
	"fmt"
	"net/http"
	"regexp"
)

//HTTPSv 创建http服务器
type HTTPSv struct {
	port          string //:80
	mux           *http.ServeMux
	RegularRoute  map[string]func(http.ResponseWriter, *http.Request)
	AbsoluteRoute map[string]func(http.ResponseWriter, *http.Request)
}

//Test 测试对象是否创建
func (sv *HTTPSv) Test() {
	fmt.Println("Http Sever Create...")
}

//注册路由
func (sv HTTPSv) setRoute() *http.ServeMux {
	sv.mux = http.NewServeMux()
	//优先绝对路由注册
	for rt, fn := range sv.AbsoluteRoute {
		sv.mux.HandleFunc(rt, fn)
	}
	//正则路由注册
	sv.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		found := false
		url := r.URL.Path
		for rt, fn := range sv.RegularRoute {
			if match, _ := regexp.MatchString(rt, url); match {

				found = true
				fn(w, r)
				return

			}
		}
		if found == false {
			http.NotFoundHandler().ServeHTTP(w, r)
			return
		}
	})
	return sv.mux
}

//SetRegularRoute 设置正则匹配url参数map
func (sv *HTTPSv) SetRegularRoute(m map[string]func(http.ResponseWriter, *http.Request)) {

	for k, v := range m {
		if sv.RegularRoute == nil {
			sv.RegularRoute = make(map[string]func(http.ResponseWriter, *http.Request))
		}
		sv.RegularRoute[k] = v
	}
	// sv.regularRoute = m
}

//SetAbsoluteRoute 设置绝对匹配url参数map
func (sv *HTTPSv) SetAbsoluteRoute(m map[string]func(http.ResponseWriter, *http.Request)) {
	for k, v := range m {
		if sv.AbsoluteRoute == nil {
			sv.AbsoluteRoute = make(map[string]func(http.ResponseWriter, *http.Request))
		}
		sv.AbsoluteRoute[k] = v
	}
	// sv.absoluteRoute = m
}

//Listen 开启server
func (sv HTTPSv) Listen(port string) error {

	mux := sv.setRoute()                  //设置路由mux--regular、absolute
	err := http.ListenAndServe(port, mux) //监听，启动路由mux
	return err

}
