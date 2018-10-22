# Go中的请求处理概述

> 原文链接 [A Recap of Request Handling in Go](https://www.alexedwards.net/blog/a-recap-of-request-handling)

<!-- TOC -->

- [Go中的请求处理概述](#go中的请求处理概述)
    - [自定义处理程序](#自定义处理程序)
    - [作为处理程序的功能](#作为处理程序的功能)
    - [DefaultServeMux](#defaultservemux)

<!-- /TOC -->

使用Go处理HTTP请求主要涉及两件事：ServeMuxes和Handlers。

甲[ServeMux](http://golang.org/pkg/net/http/#ServeMux)本质上是一个HTTP请求路由器（或*多路复用器*）。它将传入的请求与预定义的URL路径列表进行比较，并在找到匹配项时调用路径的关联处理程序。

处理程序负责编写响应头和主体。几乎任何对象都可以是处理程序，只要它满足[`http.Handler`接口即可](http://golang.org/pkg/net/http/#Handler)。在非专业术语中，这仅仅意味着它必须具有`ServeHTTP`以下签名的方法：

`ServeHTTP(http.ResponseWriter, *http.Request)`

Go的HTTP包附带了一些函数来生成常用处理程序，例如[`FileServer`](http://golang.org/pkg/net/http/#FileServer)， [`NotFoundHandler`](http://golang.org/pkg/net/http/#NotFoundHandler)和[`RedirectHandler`](http://golang.org/pkg/net/http/#RedirectHandler)。让我们从一个简单但人为的例子开始：

```
$ mkdir handler-example
$ cd handler-example
$ touch main.go
```

文件：main.go

```go
package main

import (
    "log"
    "net/http"
)

func main() {
    mux := http.NewServeMux()

    rh := http.RedirectHandler("http://example.org", 307)
    mux.Handle("/foo", rh)

    log.Println("Listening...")
    http.ListenAndServe(":3000", mux)
}
```

让我们快速介绍一下：

*   在`main`函数中，我们使用[`http.NewServeMux`](http://golang.org/pkg/net/http/#NewServeMux)函数来创建一个空的ServeMux。
*   然后我们使用该[`http.RedirectHandler`](http://golang.org/pkg/net/http/#RedirectHandler)函数创建一个新的处理程序。该处理程序307将其接收的所有请求重定向到`http://example.org`。
*   接下来我们使用该[`mux.Handle`](http://golang.org/pkg/net/http/#ServeMux.Handle)函数向我们的新ServeMux注册它，因此它充当所有带URL路径的传入请求的处理程序`/foo`。
*   最后，我们创建一个新服务器并开始使用该[`http.ListenAndServe`](http://golang.org/pkg/net/http/#ListenAndServe)函数侦听传入请求，并传入我们的ServeMux以匹配请求。

继续运行应用程序：

```
$ go run main.go
```

并访问[`http://localhost:3000/foo`](http://localhost:3000/foo)您的浏览器。您应该会发现您的请求已成功重定向。

您可能已经注意到了一些有趣的东西：ListenAndServe函数的签名是`ListenAndServe(addr string, handler Handler)`，但我们传递了一个ServeMux作为第二个参数。

我们能够做到这一点，因为ServeMux类型也[有一个`ServeHTTP`方法](http://golang.org/pkg/net/http/#ServeMux.ServeHTTP)，这意味着它也满足Handler接口。

对我来说，它简化了将ServeMux视为*一种特殊处理程序的*事情，而不是提供响应本身将请求传递给第二个处理程序。这并不像它最初听起来那么大的飞跃 - 将处理程序链接在一起在Go中相当普遍。

## 自定义处理程序

让我们创建一个自定义处理程序，它以给定格式响应当前本地时间：

```go
type timeHandler struct {
    format string
}

func (th *timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    tm := time.Now().Format(th.format)
    w.Write([]byte("The time is: " + tm))
}
```

这里的确切代码并不太重要。

所有真正重要的是我们有一个对象（在这种情况下它是一个`timeHandler`结构，但它同样可以是一个字符串或函数或其他任何东西），并且我们已经实现了一个带有签名的方法`ServeHTTP(http.ResponseWriter, *http.Request)`。这就是我们制作处理程序所需的全部内容。

让我们把它嵌入一个具体的例子中：

文件：main.go

```go
package main

import (
    "log"
    "net/http"
    "time"
)

type timeHandler struct {
    format string
}

func (th *timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    tm := time.Now().Format(th.format)
    w.Write([]byte("The time is: " + tm))
}

func main() {
    mux := http.NewServeMux()

    th := &timeHandler{format: time.RFC1123}
    mux.Handle("/time", th)

    log.Println("Listening...")
    http.ListenAndServe(":3000", mux)
}
```

在`main`函数中，我们`timeHandler`使用`&`符号产生指针，以与完全正常结构完全相同的方式初始化。然后，与前面的示例一样，我们使用该`mux.Handle`函数将其注册到ServeMux。

现在，当我们运行应用程序时，ServeMux会将任何`/time`直接请求传递给我们的`timeHandler.ServeHTTP`方法。

来吧试一试：[`http://localhost:3000/time`](http://localhost:3000/time)。

另请注意，我们可以轻松地`timeHandler`在多个路径中重复使用：

```go
func main() {
    mux := http.NewServeMux()

    th1123 := &timeHandler{format: time.RFC1123}
    mux.Handle("/time/rfc1123", th1123)

    th3339 := &timeHandler{format: time.RFC3339}
    mux.Handle("/time/rfc3339", th3339)

    log.Println("Listening...")
    http.ListenAndServe(":3000", mux)
}
```

## 作为处理程序的功能

对于简单的情况（如上面的示例），定义新的自定义类型和ServeHTTP方法感觉有点冗长。让我们看一个替代方法，我们利用Go的[`http.HandlerFunc`](http://golang.org/pkg/net/http/#HandlerFunc)类型来强制正常的函数来满足Handler接口。

任何具有签名的函数`func(http.ResponseWriter, *http.Request)`都可以转换为HandlerFunc类型。这很有用，因为HandleFunc对象带有一个内置的`ServeHTTP`方法 - 相当巧妙和方便 - 执行原始函数的内容。

如果这听起来令人困惑，请尝试查看[相关的源代码](https://golang.org/src/net/http/server.go?s=57023:57070#L1904)。你会发现它是一个非常简洁的方法，使函数满足Handler接口。

让我们使用这种技术重现timeHandler应用程序：

文件：main.go

```go
package main

import (
    "log"
    "net/http"
    "time"
)

func timeHandler(w http.ResponseWriter, r *http.Request) {
    tm := time.Now().Format(time.RFC1123)
    w.Write([]byte("The time is: " + tm))
}

func main() {
    mux := http.NewServeMux()

    // Convert the timeHandler function to a HandlerFunc type
    th := http.HandlerFunc(timeHandler)
    // And add it to the ServeMux
    mux.Handle("/time", th)

    log.Println("Listening...")
    http.ListenAndServe(":3000", mux)
}
```

实际上，将函数转换为HandlerFunc类型，然后将其添加到ServeMux中就像这样常见，Go提供了一种快捷方式：[`mux.HandleFunc`](http://golang.org/pkg/net/http/#ServeMux.HandleFunc)方法。

`main()`如果我们使用这个快捷方式，这就是函数的样子：

```go
func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("/time", timeHandler)

    log.Println("Listening...")
    http.ListenAndServe(":3000", mux)
}
```

大多数时候使用函数作为这样的处理程序很有效。但是当事情变得越来越复杂时，会有一些限制。

您可能已经注意到，与以前的方法不同，我们必须对`timeHandler`函数中的时间格式进行硬编码。*当我们想要将信息或变量传递`main()`给处理程序时会发生什么？*

一个简洁的方法是将我们的处理程序逻辑放入一个闭包中，并*关闭*我们想要使用的变量：

文件：main.go<

```go
package main

import (
    "log"
    "net/http"
    "time"
)

func timeHandler(format string) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        tm := time.Now().Format(format)
        w.Write([]byte("The time is: " + tm))
    }
    return http.HandlerFunc(fn)
}

func main() {
    mux := http.NewServeMux()

    th := timeHandler(time.RFC1123)
    mux.Handle("/time", th)

    log.Println("Listening...")
    http.ListenAndServe(":3000", mux)
}
```

该`timeHandler`功能现在具有微妙的不同作用。我们现在使用它来*返回处理程序*，而不是将函数强制转换为处理程序（就像我们之前所做的那样）。这项工作有两个关键因素。

首先它创建`fn`一个匿名函数，它访问 - 或者关闭 - `format`形成一个*闭包*的变量。无论我们如何处理闭包，它总是能够访问创建它的作用域的局部变量 - 在这种情况下意味着它总是可以访问`format`变量。

其次我们的封闭有签名`func(http.ResponseWriter, *http.Request)`。您可能还记得，这意味着我们可以将其转换为HandlerFunc类型（以便它满足Handler接口）。`timeHandler`然后我们的函数返回转换后的闭包。

在这个例子中，我们刚刚将一个简单的字符串传递给处理程序。但在实际应用程序中，您可以使用此方法传递数据库连接，模板映射或任何其他应用程序级上下文。它是使用全局变量的一个很好的替代方案，并且具有为测试制作整洁的自包含处理程序的额外好处。

您可能还会看到相同的模式，如下所示：

```go
func timeHandler(format string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tm := time.Now().Format(format)
        w.Write([]byte("The time is: " + tm))
    })
}
```

或者在返回时使用隐式转换为HandlerFunc类型：

```go
func timeHandler(format string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tm := time.Now().Format(format)
        w.Write([]byte("The time is: " + tm))
    }
}
```

## DefaultServeMux

您可能已经看到很多地方提到过DefaultServeMux，从最简单的Hello World示例到Go源代码。

我花了很长时间才意识到它并不特别。DefaultServeMux只是一个普通的'ServeMux'，就像我们已经使用的那样，在使用HTTP包时默认情况下会实例化。以下是Go源代码中的相关行：

`var DefaultServeMux = NewServeMux()`

通常，您不应使用DefaultServeMux，因为**它会带来安全风险**。

由于DefaultServeMux存储在全局变量中，因此任何程序包都可以访问它并注册路由 - 包括应用程序导入的任何第三方软件包。如果其中一个第三方软件包遭到破坏，他们可以使用DefaultServeMux向Web公开恶意处理程序。

因此，根据经验，避免使用DefaultServeMux是一个好主意，而是使用您自己的本地范围的ServeMux，就像我们到目前为止一样。但如果你确定决定使用它......

HTTP包提供了一些使用DefaultServeMux的快捷方式：[http.Handle](http://golang.org/pkg/net/http/#Handle)和[http.HandleFunc](http://golang.org/pkg/net/http/#HandleFunc)。这些与我们已经查看过的同名函数完全相同，不同之处在于它们将处理程序添加到DefaultServeMux而不是您创建的处理程序。

此外，如果没有提供其他处理程序（即第二个参数设置为`nil`），ListenAndServe将回退到使用DefaultServeMux 。

因此，作为最后一步，让我们更新我们的timeHandler应用程序以使用DefaultServeMux：

```go
package main

import (
    "log"
    "net/http"
    "time"
)

func timeHandler(format string) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        tm := time.Now().Format(format)
        w.Write([]byte("The time is: " + tm))
    }
    return http.HandlerFunc(fn)
}

func main() {
    // Note that we skip creating the ServeMux...

    var format string = time.RFC1123
    th := timeHandler(format)

    // We use http.Handle instead of mux.Handle...
    http.Handle("/time", th)

    log.Println("Listening...")
    // And pass nil as the handler to ListenAndServe.
    http.ListenAndServe(":3000", nil)
}
```
