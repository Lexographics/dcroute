package dcroute

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
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
	commands map[string]Command
}

func New(token string) *Router {
	r := &Router{
		errorFunc: func(ctx *Context) error {
			return nil
		},
		session:  nil,
		prefix:   "",
		groups:   []*Group{},
		commands: map[string]Command{},
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Panicf("Failed to initialize discord bot: %s", err.Error())
	}

	r.session = session
	session.AddHandler(r.handlerMessageCreate)
	session.AddHandler(r.handlerInteractionCreate)
	session.Identify.Intents = discordgo.IntentGuildMessages

	return r
}

func (r *Router) Group() *Group {
	g := Group{
		router:       r,
		messageFuncs: map[string]HandlerFunc{},
		commandFuncs: map[string]HandlerFunc{},
		errorFunc: func(ctx *Context) error {
			return nil
		},
		middlewares: []MiddlewareFunc{},
	}

	r.groups = append(r.groups, &g)
	return &g
}

func (r *Router) Start() error {
	err := r.session.Open()

	if err != nil {
		return errors.New("Failed to initialize websocket connection: " + err.Error())
	}

	for name, command := range r.commands {
		_, err := r.Session().ApplicationCommandCreate(r.Session().State.User.ID, command.GuildID, &discordgo.ApplicationCommand{
			Name:                     name,
			Description:              command.Description,
		})
		
		if err != nil {
			return err
		}
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
	fmt.Printf("Discord bot '%s' started\n", r.session.State.User.Username)
	fmt.Print(reset)

	return nil
}

func (r *Router) SetErrorFunc(f HandlerFunc) {
	r.errorFunc = f
}

func (r *Router) SetPrefix(prefix string) {
	r.prefix = prefix
}

func (r *Router) Wait() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func (r *Router) Session() *discordgo.Session {
	return r.session
}

func (r *Router) CreateChannel(args CreateChannelArgs) (*discordgo.Channel, error) {
	channel, err := r.session.GuildChannelCreateComplex(args.GuildID, discordgo.GuildChannelCreateData{
		Name:                 args.Name,
		Type:                 discordgo.ChannelType(args.Type),
		Topic:                args.Topic,
		Bitrate:              0,
		UserLimit:            0,
		RateLimitPerUser:     0,
		Position:             0,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{},
		ParentID:             args.CategoryID,
		NSFW:                 false,
	})

	return channel, err
}

func (r *Router) CreateEmoji(guildID string, name string, path string) (*discordgo.Emoji, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	b64 := base64.StdEncoding.EncodeToString(bytes)
	if strings.HasSuffix(file.Name(), ".png") {
		b64 = "data:image/png;base64," + b64
	} else if strings.HasSuffix(file.Name(), ".jpg") {
		b64 = "data:image/jpg;base64," + b64
	} else {
		return nil, errors.New("Invalid file extension")
	}

	emoji, err := r.Session().GuildEmojiCreate(guildID, &discordgo.EmojiParams{
		Name:  name,
		Image: b64,
		Roles: []string{},
	})

	return emoji, err
}

func (r *Router) processMessageCreate(cmd string, group *Group, ctx *Context) {
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

func (r *Router) processInteractionCreate(cmd string, group *Group, ctx *Context) {
	if h, ok := group.commandFuncs[cmd]; ok {
		h(ctx)
	}
}

func (r *Router) handlerMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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

		MessageCreate:     m,
		InteractionCreate: nil,
		session:           s,
		router:            r,
	}

	if r.prefix != "" {
		if !strings.HasPrefix(m.Content, r.prefix) {
			return
		}
	}

	cmd := strings.TrimPrefix(m.Content, r.prefix)

	for _, group := range r.groups {
		r.processMessageCreate(cmd, group, ctx)
	}
}

func (r *Router) handlerInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd := i.ApplicationCommandData().Name

	ctx := &Context{
		Sender: &User{
			Username: i.Member.User.Username,
			ID:       i.Member.User.ID,
		},
		MessageID: i.ID,
		ChannelID: i.ChannelID,
		GuildID:   i.GuildID,

		MessageCreate:     nil,
		InteractionCreate: i,
		session:           s,
		router:            r,
	}

	for _, group := range r.groups {
		r.processInteractionCreate(cmd, group, ctx)
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
