package main

import (
	"multibot/context"

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
	return "template"
}

// GetDescription function returns plugin description
func GetDescription() string {
	return "This is template plugin for multibot"
}

// GetCommands return plugin commands for bot
func GetCommands() []string {
	return []string{}
}

// UpdateHandler function call for each update
func UpdateHandler(update tgbotapi.Update) (err error) {
	return nil
}

// RunCommand handler start if bot get one of commands
func RunCommand(command string, update tgbotapi.Update) (err error) {
	return
}

// StartCommand handler start if bot get one command 'start'
func StartCommand(update tgbotapi.Update) (err error) {
	return
}
