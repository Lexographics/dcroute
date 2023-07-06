package dcroute

type Group struct {
	router *Router

	messageFuncs map[string]HandlerFunc
	commandFuncs map[string]HandlerFunc

	middlewares []MiddlewareFunc

	errorFunc HandlerFunc
}

func (g *Group) Message(name string, f HandlerFunc) {
	g.messageFuncs[name] = f
}

func (g *Group) Command(name string, command Command, f HandlerFunc) {
	g.commandFuncs[name] = f
	g.router.commands[name] = command
}

func (g *Group) Use(f MiddlewareFunc) {
	g.middlewares = append(g.middlewares, f)
}

func (g *Group) SetErrorFunc(f HandlerFunc) {
	g.errorFunc = f
}
