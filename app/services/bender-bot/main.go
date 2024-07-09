package main

import (
	"fmt"
	"os"

	"github.com/zmoog/go-bender/foundation/logger"

	"github.com/zmoog/go-bender/bot"
	"github.com/zmoog/go-bender/bot/commands"
	"github.com/zmoog/go-bender/scraper/jsonscraper"
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

	jsonscraper := jsonscraper.New()

	bender := bot.New(log, token)
	bender.AddCommand(commands.ListAppleProducts(jsonscraper))

	err = bender.Run()
	if err != nil {
		log.Errorw("Error running bot: %v", err)
	}

	log.Info("Bot stopped")
}
