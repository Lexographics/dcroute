package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	dcroute "github.com/Lexographics/dcroute/v1"
	"github.com/Lexographics/dcroute/v1/utils"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		panic("Invalid token, see '.env.example'")
	}

	r := dcroute.New(token)

	userGroup := r.Group()

	userGroup.Use(func(ctx *dcroute.Context) error {
		if ctx.ChannelID == "11111111111111111111" {
			return errors.New("Wrong channel")
		}

		return nil
	})

	userGroup.Message("create-channel", func(ctx *dcroute.Context) error {
		channel, err := ctx.Router().Session().Channel(ctx.ChannelID)
		if err != nil {
			return err
		}

		err = ctx.Router().CreateChannel(dcroute.CreateChannelArgs{
			GuildID:    ctx.GuildID,
			Name:       "test-channel",
			Topic:      "Topic",
			CategoryID: channel.ParentID,
			Type:       dcroute.ChannelTypeGuildText,
		})

		return nil
	})

	userGroup.Message("ping", func(ctx *dcroute.Context) error {
		return ctx.SendReply(ctx.ChannelID, ctx.MessageID, ctx.GuildID, "pong")
	})

	userGroup.Message("image", func(ctx *dcroute.Context) error {
		fmt.Printf("Got message form %s with id %s\n", ctx.Sender.Username, ctx.Sender.ID)

		err := ctx.SendFile(ctx.ChannelID, "image.jpg", "example/image.png")
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		
		return nil
	})

	userGroup.Message("emoji", func(ctx *dcroute.Context) error {
		return ctx.SendReply(ctx.ChannelID, ctx.MessageID, ctx.GuildID, utils.GetEmoji("new_emoji", "1126605484644909106"))
	})

	userGroup.Message("react", func(ctx *dcroute.Context) error {
		ctx.SendReaction(ctx.ChannelID, ctx.MessageID, utils.GetReaction("new_emoji", "1126605484644909106"))
		time.Sleep(time.Second * 2)
		ctx.RemoveReaction(ctx.ChannelID, ctx.MessageID, utils.GetReaction("new_emoji", "1126605484644909106"))
		return nil
	})
	
	userGroup.Command("test", dcroute.Command{
		Description: "Description",
		GuildID:     "1095712396921806849",
	}, func(ctx *dcroute.Context) error {
		ctx.CommandRespond("test2")
		return nil
	})
	
	r.Start()
	
	// err := r.CreateEmoji("1095712396921806849", "new_emoji", "example/image.png")

	r.Wait()
}
