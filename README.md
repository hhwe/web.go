# web.go -- golang restful api framework

web.go is a golang web framework, you can use it build a slightly restful api.

## middlerware

there is three middlware in package: logging, recovery. certainly, you can rewrite it in your project.

by default a middlware must apply like:

```go 
type middleware func(http.Handler) http.Handler
```

## api route

router is the core part of this package, I use a regexp package to mapping url and handler function, so that the dynamic url also can be handled. 

