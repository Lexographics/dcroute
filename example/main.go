package main

import (
	"errors"
	"fmt"
	"os"

	dcroute "github.com/Lexographics/dcroute/v1"
	"github.com/joho/godotenv"
)


func main() {
	godotenv.Load()
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		panic("Invalid token, see '.env.example'")
	}

	h := dcroute.New(token)

	userGroup := h.Group()

	userGroup.Use(func(ctx *dcroute.Context) error {
		if ctx.ChannelID == "11111111111111111111" {
			return errors.New("Wrong channel")
		}

		return nil
	})

	userGroup.Message("ping", func(ctx *dcroute.Context) error {
		fmt.Printf("Got message form %s with id %s\n", ctx.Sender.Username, ctx.Sender.ID)

		ctx.SendMessage(ctx.ChannelID, "pong")
		return nil
	})

	h.Start()
}