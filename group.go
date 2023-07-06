package dcroute

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

type Group struct {
	router *Router

	messageFuncs map[string]HandlerFunc
	commandFuncs map[string]HandlerFunc

	middlewares []MiddlewareFunc

	errorFunc HandlerFunc
}

func (g *Group) Message(command string, f HandlerFunc) {
	g.messageFuncs[command] = f
}

func (g *Group) Command(command string, description string, guildID string, f HandlerFunc) error {
	if g.router.session == nil {
		return errors.New("Session is nil")
	}

	g.commandFuncs[command] = f
	_, err := g.router.Session().ApplicationCommandCreate(g.router.Session().State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:                     command,
		Description:              description,
	})

	if err != nil {
		return err
	}

	return nil
}

func (g *Group) Use(f MiddlewareFunc) {
	g.middlewares = append(g.middlewares, f)
}

func (g *Group) SetErrorFunc(f HandlerFunc) {
	g.errorFunc = f
}
