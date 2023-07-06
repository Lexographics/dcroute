package dcroute

type Group struct {
	router *Router

	messageFuncs map[string]HandlerFunc

	middlewares []MiddlewareFunc

	errorFunc HandlerFunc
}

func (g *Group) Message(command string, f HandlerFunc) {
	g.messageFuncs[command] = f
}

func (g *Group) Use(f MiddlewareFunc) {
	g.middlewares = append(g.middlewares, f)
}

func (g *Group) SetErrorFunc(f HandlerFunc) {
	g.errorFunc = f
}
