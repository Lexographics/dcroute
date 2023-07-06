package dcroute

import (
	"bufio"
	"os"

	"github.com/bwmarrin/discordgo"
)

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

func (c *Context) SendFile(channel string, name string, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)

	_, err = c.session.ChannelFileSend(channel, name, reader)
	return err
}