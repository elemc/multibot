// -*- Go -*-

package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

var bot *tgbotapi.BotAPI

// BotServe function run telegram bot listener
func BotServe() (err error) {
	var updates <-chan tgbotapi.Update
	defer wg.Done()

	updOptions := tgbotapi.NewUpdate(0)
	updOptions.Timeout = 60

	if updates, err = bot.GetUpdatesChan(updOptions); err != nil {
		return
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		go botEachUpdateHandler(update)

		if update.Message.Command() != "" {
			if update.Message.Command() == "start" {
				go botStartHandler(update)
			} else {
				go botCommandsHandler(update)
			}
		}
	}

	return
}

func botEachUpdateHandler(update tgbotapi.Update) {
	for name, botPlugin := range botPlugins {
		if err := botPlugin.EachUpdateHandler(update); err != nil {
			log.Errorf("Error in plugin %s on each update handler: %s", name, err)
		}
	}
}

func botCommandsHandler(update tgbotapi.Update) {
	cmd := update.Message.Command()
	if botPlugin, ok := botPluginsByCommand[cmd]; ok {
		if err := botPlugin.RunCommandHandler(cmd, update); err != nil {
			log.Errorf("Unable to run command '%s' for plugin '%s': %s", cmd, botPlugin.Name, err)
		}
	}
}

func botStartHandler(update tgbotapi.Update) {
	for name, botPlugin := range botPlugins {
		if err := botPlugin.StartCommandHandler(update); err != nil {
			log.Errorf("Unable to run command start for plugin '%s': %s", name, err)
		}
	}
}
