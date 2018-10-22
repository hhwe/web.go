## 中间件

go中处理请求最常见的就是中间件，实际上是一种切片编程理念AOC，在python中就是一个装饰器

如果我们在每个请求中有普遍的处理需求, 我们需要为每个请求设立一个堆栈来隔离彼此并且保持运行顺序, 中间件可以很好的解决这类问题

### 1. 以类型的形式实现

上篇博客中，我们探讨过Go语言实现Web最核心的部分：

```go
http.ListenAndServe(":8000", handler)
```

复制代码http包里面的ListenAndServe函数接受两个参数，即监听地址和处理接口handler，handler是一个接口，我们需要实现这个接口中的唯一方法ServeHTTP便可以实现上述的函数，因此我们处理的整个逻辑和流程都会在这个handler里面，下面我们先来看一个最简单的handler实现。

```go
package main

import (
	"net/http"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func main() {
	http.ListenAndServe(":8000", http.HandlerFunc(myHandler))
}
```

复制代码上面的代码中我们定义一个myHandler，它接受http.ResponseWriter,*http.Request两个参数然后再向ResponseWriter中写入Hello World,在main函数中，我们直接使用了ListenAndServe方法监听本地的8000端口，注意由于Go语言的强类型性，ListenAndServe的第二个参数类型是Handler，因此我们想要将myHandler传递给ListenAndServe就必须实现ServeHTTP这个方法。但其实Go源码里面已经帮我们实现了这个方法。

```go
// Handler that calls f.
type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}
```

复制代码可以看到，Go语言将func(ResponseWriter, *Request)这种类型的函数直接定义了类型HandlerFunc，而且还实现了ServeHTTP这个方法，但是这个方法本身并没有实现任何逻辑，需要我们自己来实现。因此我们实现了myHandler这个方法，它将输出一个最简单的HelloWorld响应。随后我们可以用curl来测试一下：

```
$ curl localhost:8000
Hello World
```

复制代码可以看到，我们通过curl请求本地的8000端口，返回我们一个HelloWorld。这便是一个最简单的Handler实现了。但是我们的目标是实现中间件，有了上述的所采用的的方法我们就可以大致明白，myHandler应该作为最后的调用，在它之前才是中间件应该作用的地方，那么我们有了一个大致的方向，我们可以实现一个逻辑用来包含这个myHandler，但它本身也必须实现Handler这个接口，因为我们要把它传递给ListenAndServe这个方法。好，我们先大致阐述一下这个中间件的作用，它会拦截一切请求除了这个请求的host是我们想要的host，当然这个host有我们定义。

```go
type SingleHost struct {
	handler     http.Handler
	allowedHost string
}
```

复制代码于是我们定义了一个SingleHost的结构体，它里面有两个成员一个是Handler，它将是我们上述的myHandler，另一个是我们允许来请求Server的用户，这个用户他有唯一的Host，只有当他的Host满足我们的要求是才让他请求成功，否则一律返回403。
因为我们需要将这个SingleHost实例化并传递给ListenAndServe这个方法，因此它必须实现ServeHTTP这个方法，所以在ServeHTTP里面可以直接定义我们用来实现中间件的逻辑。即除非来请求的用户的Host是allowedHost否则一律返回403。

```go
func (this *SingleHost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if (r.Host == this.allowedHost) {
		this.handler.ServeHTTP(w, r)
	} else {
		w.WriteHeader(403)
	}
}
```

复制代码好，可以清楚的看到只有Request的Host==allowedHost的时候，我们才调用handler的ServeHTTP方法，否则返回403.下面是完整代码：

```go
package main

import (
	"net/http"
)

type SingleHost struct {
	handler     http.Handler
	allowedHost string
}

func (this *SingleHost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if (r.Host == this.allowedHost) {
		this.handler.ServeHTTP(w, r)
	} else {
		w.WriteHeader(403)
	}
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func main() {
	single := &SingleHost{
		handler:http.HandlerFunc(myHandler),
		allowedHost:"refuse.com",
	}
	http.ListenAndServe(":8000", single)
}
```

复制代码然后我们用curl来请求本地的8000端口，

```
$ curl --head localhost:8000
HTTP/1.1 403 Forbidden
Date: Sun, 21 Jan 2018 08:32:47 GMT
Content-Type: text/plain; charset=utf-8
```

复制代码可以看到我们在中间件中实现了只允许host为refuse.com来访问的逻辑实现了，由于curl的Host是localhost所以我们的服务器直接返回了它一个403。接下来我们改变一下allowedHost
allowedHost:"localhost:8000",
复制代码我们将allowedHost变成为localhost:8000，然后用curl测试

```
$ curl localhost:8000
Hello World
```

复制代码可以看到curl通过了中间件的并直接获得了myHandler返回的HelloWorld。

### 2. 以函数的形式实现

好，在上面我们实现了以类型为基础的中间件，可能对Node.js较熟悉的人都习惯以函数的形式实现中间件。首先，因为我们是以函数来实现中间件的因此这个函数返回的便是Handler,它会接受两个参数，一个是我们定义的myHandler，一个是allowedHost。

```go
func SingleHost(handler http.Handler, allowedHost string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Host == allowedHost {
			handler.ServeHTTP(w, r)
		} else {
			w.WriteHeader(403)
		}
	}
	return http.HandlerFunc(fn)
}
```

复制代码可以看到，我们在函数内部定义可一个匿名函数fn，这个匿名函数便是我们要返回的Handler，如果请求用户的Host满足allowedHost,便可以将调用myHandler的函数返回，否则直接返回一个操作403的函数。整个代码如下：

```go
package main

import "net/http"

func SingleHost(handler http.Handler, allowedHost string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Host == allowedHost {
			handler.ServeHTTP(w, r)
		} else {
			w.WriteHeader(403)
		}
	}
	return http.HandlerFunc(fn)
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func main() {
	single := SingleHost(http.HandlerFunc(myHandler), "refuse.com")
	http.ListenAndServe(":8000", single)
}
```

复制代码我们还是通过curl来测试一下

```
$ curl --head localhost:8000
HTTP/1.1 403 Forbidden
Date: Sun, 21 Jan 2018 08:45:39 GMT
Content-Type: text/plain; charset=utf-8
```

复制代码可以看到由于不满足refuse.com的条件，我们会得到一个403，让我们将refuse.com改为localhost:8000测试一下。

```
$ curl localhost:8000
Hello World
```

复制代码与刚才一样我们得到了HelloWorld这个正确结果。


## reference

[Simple HTTP middleware with Go](https://hackernoon.com/simple-http-middleware-with-go-79a4ad62889b)
