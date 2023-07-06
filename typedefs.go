package dcroute

type HandlerFunc func(ctx *Context) error

type MiddlewareFunc func(ctx *Context) error