// go 实现动态路由匹配

package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type App struct {
	w http.ResponseWriter
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/public/") {
			//匹配静态文件服务
		} else {
			app := &App{w}
			fmt.Println(app)
			rValue := reflect.ValueOf(app)
			fmt.Println(rValue)
			rType := reflect.TypeOf(app)
			fmt.Println(rType)
			path := strings.Split(r.URL.Path, "/")
			controlName := path[1]
			fmt.Println(controlName)
			method, exist := rType.MethodByName(controlName)
			fmt.Println(method, exist)
			if exist {
				args := []reflect.Value{rValue}
				method.Func.Call(args)
			} else {
				fmt.Fprintf(w, "method %s not found", r.URL.Path)
			}
		}
	})

	log.Fatal(http.ListenAndServe("localhost:8001", nil))
}

//控制器
func (this *App) Say() {
	fmt.Fprintf(this.w, "Say called")
}
