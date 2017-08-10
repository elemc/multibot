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
	return "reminder"
}

// GetDescription function returns plugin description
func GetDescription() string {
	return "Plugin for create edit delete remind tast and notification about it"
}

// GetCommands return plugin commands for bot
func GetCommands() []string {
	return []string{
		"add_task",
		"del_task",
	}
}

// UpdateHandler function call for each update
func UpdateHandler(update tgbotapi.Update) (err error) {
	return nil
}

// RunCommand handler start if bot get one of commands
func RunCommand(command string, update tgbotapi.Update) (err error) {
	switch command {
	case "add_task":
		addTask(update.Message)
	case "del_task":
		delTask(update.Message)
	}
	return
}

// StartCommand handler start if bot get one command 'start'
func StartCommand(update tgbotapi.Update) (err error) {
	msg := `Привет!
Тебя приветствует плагин "Напоминатель"
Для добавления задачи отправь команду /add_task.
Для удаления уже созданной задачи отправь команду /del_task.
Приятного пользования "Напоминателем"!`
	ctx.SendMessageText(update.Message.Chat.ID, msg, 0)
	return
}

func addTask(msg *tgbotapi.Message) {
	log.Infof("send command \"add_task\"")
}

func delTask(msg *tgbotapi.Message) {
	log.Infof("send command \"del_task\"")
}
