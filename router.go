package dcroute

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Router struct {
	errorFunc HandlerFunc
	session   *discordgo.Session
	prefix    string

	groups []*Group
}

func New(token string) *Router {
	h := &Router{}
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Panicf("Failed to initialize discord bot: %s", err.Error())
	}

	h.session = session
	session.AddHandler(h.handlerMessageCreate)
	session.Identify.Intents = discordgo.IntentGuildMessages

	return h
}

func (h *Router) Group() *Group {
	g := Group{
		handler:      h,
		messageFuncs: map[string]HandlerFunc{},
		errorFunc: func(ctx *Context) error {
			return nil
		},
		middlewares: []MiddlewareFunc{},
	}

	h.groups = append(h.groups, &g)
	return &g
}

func (h *Router) Start() error {
	err := h.session.Open()

	if err != nil {
		return errors.New("Failed to initialize websocket connection: " + err.Error())
	}

	text := `
  ____     ____    ____    U  ___ u   _   _  _____  U _____ u 
 |  _"\ U /"___|U |  _"\ u  \/"_ \/U |"|u| ||_ " _| \| ___"|/ 
/| | | |\| | u   \| |_) |/  | | | | \| |\| |  | |    |  _|"   
U| |_| |\| |/__   |  _ <.-,_| |_| |  | |_| | /| |\   | |___   
 |____/ u \____|  |_| \_\\_)-\___/  <<\___/ u |_|U   |_____|  
  |||_   _// \\   //   \\_    \\   (__) )(  _// \\_  <<   >>  
 (__)_) (__)(__) (__)  (__)  (__)      (__)(__) (__)(__) (__) 
`
	fmt.Print(cyan)
	fmt.Println(text)
	fmt.Print(green)
	fmt.Printf("Discord bot '%s' started\n", h.session.State.User.Username)
	fmt.Print(reset)

	return nil
}

func (h *Router) SetErrorFunc(f HandlerFunc) {
	h.errorFunc = f
}

func (h *Router) SetPrefix(prefix string) {
	h.prefix = prefix
}

func (h *Router) Wait() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func (h *Router) processGroup(cmd string, group *Group, ctx *Context) {
	handlerfn, ok := group.messageFuncs[cmd]
	if !ok {
		return
	}

	for _, fn := range group.middlewares {
		err := fn(ctx)
		if err != nil {
			return
		}
	}

	err := handlerfn(ctx)
	if err != nil {
		fmt.Println("Handler error: " + err.Error())
	}
}

func (h *Router) handlerMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	ctx := &Context{
		Sender: &User{
			Username: m.Author.Username,
			ID:       m.Author.ID,
		},
		MessageID: m.ID,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,

		MessageCreate: m,
		session:       s,
	}

	if h.prefix != "" {
		if !strings.HasPrefix(m.Content, h.prefix) {
			return
		}
	}

	cmd := strings.TrimPrefix(m.Content, h.prefix)

	for _, group := range h.groups {
		h.processGroup(cmd, group, ctx)
	}
}

const reset = "\033[0m"
const red = "\033[31m"
const green = "\033[32m"
const yellow = "\033[33m"
const blue = "\033[34m"
const purple = "\033[35m"
const cyan = "\033[36m"
const white = "\033[37m"
