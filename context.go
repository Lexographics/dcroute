package dcroute

import "github.com/bwmarrin/discordgo"

type Context struct {
	Sender    *User
	ChannelID string

	MessageCreate *discordgo.MessageCreate

	session *discordgo.Session
}

func (c *Context) SendMessage(channel string, m string) error {
	_, err := c.session.ChannelMessageSend(channel, m)
	return err
}
