package dcroute

type HandlerFunc func(ctx *Context) error

type MiddlewareFunc func(ctx *Context) error

type CreateChannelArgs struct {
	GuildID    string
	Name       string
	Topic      string
	CategoryID string
	Type       ChannelType
}

type ChannelType int

const (
	ChannelTypeGuildText          ChannelType = 0
	ChannelTypeDM                 ChannelType = 1
	ChannelTypeGuildVoice         ChannelType = 2
	ChannelTypeGroupDM            ChannelType = 3
	ChannelTypeGuildCategory      ChannelType = 4
	ChannelTypeGuildNews          ChannelType = 5
	ChannelTypeGuildStore         ChannelType = 6
	ChannelTypeGuildNewsThread    ChannelType = 10
	ChannelTypeGuildPublicThread  ChannelType = 11
	ChannelTypeGuildPrivateThread ChannelType = 12
	ChannelTypeGuildStageVoice    ChannelType = 13
	ChannelTypeGuildForum         ChannelType = 15
)

type Command struct {
	Description string
	GuildID     string
}

const MessageAny = "_*"
const MessageNotFound = "_?"