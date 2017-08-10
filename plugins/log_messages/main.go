package main

import (
	"multibot/context"

	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

var ctx *context.MultiBotContext

// InitPlugin initialize plugin if it needed
func InitPlugin(mbc *context.MultiBotContext) error {
	ctx = mbc
	return nil
}

// GetName function returns plugin name
func GetName() string {
	return "log_messages"
}

// GetDescription function returns plugin description
func GetDescription() string {
	return "Simple plugin save all messages to file for multibot"
}

// GetCommands return plugin commands for bot
func GetCommands() []string {
	return []string{}
}

// UpdateHandler function call for each update
func UpdateHandler(update tgbotapi.Update) (err error) {
	log.Debugf("%s", update.Message.Text)
	ctx.SendMessage(update.Message.Chat.ID, "Thanks", update.Message.MessageID)
	return nil
}

// RunCommand handler start if bot get one of commands
func RunCommand(command string, update tgbotapi.Update) (err error) {
	return
}
