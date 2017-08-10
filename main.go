// -*- Go -*-

package main

import (
	"flag"
	"multibot/context"
	"sync"

	"gopkg.in/telegram-bot-api.v4"

	log "github.com/sirupsen/logrus"
)

var (
	configName string
	wg         sync.WaitGroup
	botContext *context.MultiBotContext
)

func init() {
	flag.StringVar(&configName, "config", "multibot", "configuration file name")
}

func main() {
	var err error
	flag.Parse()

	if err = LoadConfig(); err != nil {
		log.Fatalf("Unable to load configuration file %s: %s", configName, err)
	}
	if level, err := log.ParseLevel(options.LogLevel); err != nil {
		log.Warnf("Unable to parse log level %s: %s", options.LogLevel, err)
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(level)
	}

	if err = InitDatabase(); err != nil {
		log.Fatalf("Unable to connect to database: %s", err)
	}

	if bot, err = tgbotapi.NewBotAPI(options.APIKey); err != nil {
		log.Fatalf("Unable to initialize telegram bot: %s", err)
	}
	log.Debug("Telegram bot initialized sucessful")
	botContext = context.InitContext(db, bot, options)

	if err = LoadPlugins(); err != nil {
		log.Fatalf("Unable to load plugins: %s", err)
	}

	wg.Add(1)
	if err = BotServe(); err != nil {
		log.Fatalf("Unable to server bot: %s", err)
	}

	log.Warnf("Application started...")
}
