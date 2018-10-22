# HTTP Request Contexts & Go

> 这篇文章是 Golang 开源库 [Negroni](https://github.com/codegangsta/negroni) 的 README.md 中推荐一篇的文章，讲的是 Golang 中如何处理请求的上下文信息。
> 
> 原文链接 [HTTP Request Contexts & Go](https://blog.questionable.services/article/map-string-interface/)

<!-- TOC -->

- [HTTP Request Contexts & Go](#http-request-contexts--go)
    - [简介](#简介)
    - [全局上下文 Map](#全局上下文-map)
    - [每个请求一个 map[string]interface](#每个请求一个-mapstringinterface)
    - [Context Structs](#context-structs)
    - [其他的解决方案？](#其他的解决方案)
    - [总结](#总结)

<!-- /TOC -->

*注：本文中 将 Handler 翻译为处理器，middleware 翻译成中间件，Request Context 翻译为请求上下文。*

## 简介

这篇文章其实也可以用 map[string]interface 作为标题。

如果我们想让HTTP请求时携带一些需要一系列处理程序或中间件处理的数据，一种典型的方式是使用 Request contexts. 这些数据可以是用户ID，CSRF token，或者一些数据表示用户是否已经登录或一些你不想在每个处理器都重复一遍的逻辑。如果你用过Django，request context 就是 `request.META` 字典。

比如：

```go
func CSRFMiddleware(http.Handler) http.Handler {
   return func(w http.ResponseWriter, r *http.Request) {
       maskedToken, err := csrf.GenerateNewToken(r)
       if err != nil {
           http.Error(w, "No good!", http.StatusInternalServerError)
           return
       }
       // 我们如何把 maskedToken 从这里传出去 。。。
   }
}

func MyHandler(w http.ResponseWriter, r *http.Request) {  
   // 。。。传到这里，不依赖 session 的存储，
   // 而且不依赖绑定另一个处理器的处理器？
   // CSRF token 怎么办？或者 request header 里面 auth-key？
   // 我们肯定不想在每个处理器里面把这些逻辑都重写一遍！
}
```

在Go语言的 web库/框架 中，我们有三种方法解决这个问题：

1. 创建一个全局 map, *http.Request 作为 key, 同步互斥写入，中间件负责清理旧的请求 ([gorilla/context](https://link.jianshu.com?t=http://www.gorillatoolkit.org/pkg/context) ).
2. 通过创建自定义的处理器类型(handler types), 为每个 request 指定一个 map ([goji](https://link.jianshu.com?t=https://goji.io/))。
3. Structs, 自定义 structs, 为它创建以指针为接收器的中间件方法，或者将这些结构体对象作为参数传递给处理器 ([gocraft/web](https://link.jianshu.com?t=https://github.com/gocraft/web))。

那么，这些方法有什么区别，各自又有什么优劣呢？

## 全局上下文 Map

[gorilla/context](http://www.gorillatoolkit.org/pkg/context) 的实现方式最简单、最容易集成到已有的架构。Gorilla 其实用的是 `map[interface{}]interface{}` 结构，这意味着你需要(且应该)为你的 keys 创建类型。好处是你能用任何类型最为一个 key; 缺点是如果不想引发 run-time 问题，你得提前实现你的 keys.

为了避免相同的类型断言污染你的处理器，你还应该为这些 types 创建一些 setters 用来存储你的上下文 map。

```go
import (
    "net/http"
    "github.com/gorilla/context"
)

type contextKey int

/ 定义一些竞争平等的 key
const csrfKey contextKey = 0
const userKey contextKey = 1

var ErrCSRFTokenNotPresent = errors.New("CSRF token not present in the request context.")

// 对每一个存储在 context map 里面的 key:value 组合，
// 我们都需要一个像这样的辅助方法，
// 不然我们就会在每个需要这个值的中间件中重复这些代码。
func GetCSRFToken(r *http.Request) (string, error) {
    val, ok := context.GetOk(r, csrfKey)
    if !ok {
        return "", ErrCSRFTokenNotPresent
    }

    token, ok := val.(string)
    if !ok {
        return "", ErrCSRFTokenNotPresent
    }

    return token, nil
}

// 一个简单的示例
func CSRFMiddleware(h http.Handler) http.Handler {
    return func(w http.ResponseWriter, r *http.Request) {
        token, err := GetCSRFToken(r)
        if err != nil {
            http.Error(w, "No good!", http.StatusInternalServerError)
            return
        }

        // 这个 map 是全局的，所以我们只调用 Set 方法
        context.Set(r, csrfKey, token)

        h.ServeHTTP(w, r)
    }
}

func ShowSignupForm(w http.ResponseWriter, r *http.Request) {

    // 我们用自己的辅助方法，
    // 所以我们不需要维护每个触发或处理 POST 请求的处理器中数据类型
    csrfToken, err := GetCSRFToken(r)
    if err != nil {
        http.Error(w, "No good!", http.StatusInternalServerError)
        return
    }

    // 我们在每一个用我们的中间件包装的处理器中都能访问到这个 token. 
    // 而不需要在每次请求中从 session 里多次读取(这种方法很慢)。
    fmt.Fprintf(w, "Our token is %v", csrfToken)
}

func main() {
    r := http.NewServeMux()
    r.Handle("/signup", CSRFMiddleware(http.HandlerFunc(ShowSignupForm)))
    // 这里必须调用 context.ClearHandler,
    // 否则将会把旧的请求留在 map 中。
    http.ListenAndServe("localhost:8000", context.ClearHandler(r))
}

```

[完整示例](t=https://gist.github.com/elithrar/015e2a561eee0ca71a77#file-gorilla-go)

优点？灵活，松耦合，第三方包容易使用。你可以把它用在任何 net/http 应用中，因为你要做的只是访问 http.Request——剩下的依赖就只有 global map 了。

缺点？global map 和其互斥器(mutexes) *有可能* 导致高链接负载，所以你需要在每个请求(即每个处理器) 里最后调用 [context.Clear()](t=http://www.gorillatoolkit.org/pkg/context#Clear). 如果忘了清理(或者对你的上层服务器处理器包了一层), 你会让自己陷入内存泄露的危险，因为老的请求依然存在 map 中。如果你的中间件用了 gorilla/context, 那么得确保你引入了 contex 而且在处理器/路由中调用了 context.ClearHandler。

## 每个请求一个 map[string]interface

作为另一种方式，Goji 提供了一个嵌入到 http.Handler 的请求上下文。 因为它绑定在了 Goji 的路由具体实现代码中， 所以它不需要是一个全局 map 而且也避免了使用互斥锁。Goji 提供了一个 [web.HandlerFunc](https://github.com/zenazn/goji/blob/master/web/web.go#L109-L128) 类型，它通过 `func(c web.C, w http.ResponseWriter, r *http.Request)` 扩展了默认的 `http.HandlerFunc`.

```go
var ErrTypeNotPresent = errors.New("Expected type not present in the request context.")

// 相对简单点：我们对每一个 *类型*, 我们只需要这一个方法
func GetContextString(c web.C, key string) (string, error) {
   val, ok := c.Env[key].(string)
   if !ok {
       return "", ErrTypeNotPresent
   }

   return val, nil
}

// 一个简单示例
func CSRFMiddleware(c *web.C, h http.Handler) http.Handler {
   fn := func(w http.ResponseWriter, r *http.Request) {
       maskedToken, err := GenerateToken(r)
       if err != nil {
           http.Error(w, "No good!", http.StatusInternalServerError)
           return
       }

       // Goji 只在你需要的时候分配一个 map
       if c.Env == nil {
           c.Env = make(map[string]interface{})
       }

       // 不是全局变量 —— 只是一个 context map 的引用，
       // 它显式传递给我们我的处理器。
       c.Env["csrf_token"] = maskedToken

       h.ServeHTTP(w, r)
   }

   return http.HandlerFunc(fn)
}

// Goji 的 web.HandlerFunc 类型是 net/http.HandlerFunc 的扩展，
// 它会被参数传递给 request context(又称 web.C.Env)
func ShowSignupForm(c web.C, w http.ResponseWriter, r *http.Request) {
   csrfToken, err := GetContextString(c, "csrf_token")
   if err != nil {
       http.Error(w, "No good!", http.StatusInternalServerError)
       return
   }

   fmt.Fprintf(w, "Our token is %v", csrfToken)
}
```

[完整示例](https://gist.github.com/elithrar/015e2a561eee0ca71a77#file-goji-go)

最大最直接的好处就是性能提升，因为 Goji 只在你需要的时候它才会分配一个 map: 不存在全局 map锁 的问题。注意，对于很多应用，数据库和模板渲染时将会是性能的瓶颈所在，所以“真正”的影响可能很小，但是这正是二者之间一个合理的联系。

最有用的是，你依然可以写不需要应用信息的模块化中间件，而且如果你想用 request context，马上就可以用，但是除了这个，它就是个 http.Handler. 缺点是，你依然需要维护从 context 拿到的数据的类型，尽管像 gorilla/context, 我们可以用 helper方法 简化这些操作，但 `map[string]interface{}` 还会限制我们只能用 `string key`: 对大多数人(包括我)来说很简单，但可能对有些人来说缺少灵活性。

## Context Structs

第三种方式是为每一个请求初始化一个 struct, 然后把我们的 中间件/处理器 定义为它的方法。最大的好处是类型安全：我们明确地定义了请求的上下文，所以我们知道它的类型(除非我们定义了一个 `interface{}` 的字段).

当然，在你可以保证类型安全的同时你也失去了灵活性。你不能创建一个可以使用 `func(http.Handler) http.Handler` 模式的“模块化”中间件，因为这个中间件不知道你的请求上线文 struct 是什么。它可以提供他自己的 struct 并嵌入到你的 struct 中，但这依然不能重用。然而，这是一种不错的方法：除了 `interface{}` 以外，你不需要维护类型。

```go
import (
   "fmt"
   "log"
   "net/http"

   "github.com/gocraft/web"
)

type Context struct {
   CSRFToken string
   User     string
}

// 我们的中间件*和*处理器必须都定义为 context struct 的方法，或者接受 context struct 作为第一个参数。
// 这就会把 处理器/中间件 和具体的应用 架构/设计 绑在一起。
func (c *Context) CSRFMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
   token, err := GenerateToken(r)
   if err != nil {
       http.Error(w, "No good!", http.StatusInternalServerError)
       return
   }

   c.CSRFToken = token
   next(w, r)
}

func (c *Context) ShowSignupForm(w web.ResponseWriter, r *web.Request) {
   // 不需要维护它：我们知道它的类型。
   // 我们能直接用这个值。
   fmt.Fprintf(w, "Our token is %v", c.CSRFToken)
}

func main() {
   router := web.New(Context{}).Middleware((*Context).CSRFMiddleware)
   router.Get("/signup", (*Context).ShowSignupForm)

   err := http.ListenAndServe(":8000", router)
   if err != nil {
       log.Fatal(err)
   }
}

```

[完整示例](https://gist.github.com/elithrar/015e2a561eee0ca71a77#file-gocraftweb-go)

优点很显然：没有类型维护！我们的 struct 有明确的类型而且对每个请求进行初始化，然后传递给我们的 中间件/处理器。但，缺点是社区的中间件再也不能“即插即用”了，因为它们的设计肯定不会考虑我们的上下文。

我们可以将它们的类型匿名嵌入到我们的 struct 中，但是这会让我们的代码变得很乱，而且如果它们命名和我们的一样，还是行不通。现实的解决方案是 fork 一份代码，然后修改以让我们的 struct 能够使用，代价就是时间和精力。gocraft/web 也是用他们自己的类型包装了 ResponseWriter interface/Request struct，这和框架本身关系更密切。

## 其他的解决方案？

一个解决方案是，Go 的 http.Request struct 提供一个 Context 属性，但是事实上，以合乎情理的方式实现一个可以适合一般场景的请求上下文并不是那么容易。

这个属性应该是一个 `map[string]interface{}` (或者以`interface{}` 作为 key). 这意味着我们要么需要为使用者初始化这个 map ——这对于那些不需要请求上下文的使用者来说没什么用。或者需要包使用者在使用之前检查，这会让一些新手在一开始感到困惑，为什么他们的应用对于有些请求报错而其他的却没有。

我觉得着对于 net/http 本身来说都不是什么大的障碍，但是 Go 的设计理念是清晰、容易理解 —— 这么做代价是有时有点啰嗦 —— 这有可能有些违背设计理念。我也认为有些 第三方包/框架 给我们选择没什么不好的：选择最适合你的方式。

## 总结

所以你会为你的项目选择什么方式呢？这取决于你的使用场景。想写一个便于使用者调用请求上下文，而且相对独立的包？gorilla/context 可能是个不错的选择（提醒：别忘了调用 ClearHandler!）。从头写或者打算扩展一个 net/http 应用？Goji 可以很容易上手。从零做起？gocraft/web 整套解决方案也许会比较适合你。

我个人比较喜欢 Goji 的解决方案：我不介意写一些辅助方法去维护我常用的请求上下文类型(CSRF token，用户名，等), 而且我尽量避免全局 map. 对我来说，可以很容易写一些可以被别人集成到他们项目中的中间件。但这只是我自己的使用场景，所以你得先针对自己的情况做些调查。
