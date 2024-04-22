package main

import (
	"context"
	"dexbot/bot"
	"dexbot/service"
	"os"
	"os/signal"
)

func main() {
	service.GetConfig()

	ctx, cancel := context.WithCancel(context.Background())
	bot.StartBot(ctx)

	//gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	cancel()
}
