package dcroute

import (
	"bufio"
	"os"

	"github.com/bwmarrin/discordgo"
)

type Context struct {
	Sender    *User
	MessageID string
	ChannelID string
	GuildID   string

	MessageCreate *discordgo.MessageCreate

	router  *Router
	session *discordgo.Session
}

func (c *Context) SendMessage(channelID string, m string) error {
	_, err := c.session.ChannelMessageSend(channelID, m)
	return err
}

func (c *Context) SendReply(channelID string, messageID string, guildID string, m string) error {
	_, err := c.session.ChannelMessageSendReply(channelID, m, &discordgo.MessageReference{
		MessageID: messageID,
		ChannelID: channelID,
		GuildID:   guildID,
	})

	return err
}

func (c *Context) SendReaction(channelID string, messageID string, reaction string) error {
	return c.session.MessageReactionAdd(channelID, messageID, reaction)
}

func (c *Context) RemoveReaction(channelID string, messageID string, reaction string) error {
	return c.session.MessageReactionRemove(channelID, messageID, reaction, c.Router().Session().State.User.ID)
}

func (c *Context) SendFile(channelID string, name string, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)

	_, err = c.session.ChannelFileSend(channelID, name, reader)
	return err
}

func (c *Context) Router() *Router {
	return c.router
}