package main

import (
	"fmt"
	"github.com/zmoog/go-bender/foundation/logger"
	"os"

	"github.com/zmoog/go-bender/bot"
)

func main() {
	log, err := logger.New("bender")
	if err != nil {
		fmt.Printf("Error creating logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	token, ok := os.LookupEnv("DISCORD_TOKEN")
	if !ok {
		log.Error("DISCORD_TOKEN environment variable not set")
		return
	}

	log.Info("Starting bot")

	bender := bot.New(log, token)
	err = bender.Run()
	if err != nil {
		log.Errorw("Error running bot: %v", err)
	}

	log.Info("Bot stopped")
}
