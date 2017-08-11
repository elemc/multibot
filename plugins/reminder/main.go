package main

import (
	"fmt"
	"multibot/context"

	log "github.com/sirupsen/logrus"

	"gopkg.in/telegram-bot-api.v4"
)

const (
	taskAddCommand  = "task_add"
	taskDelCommand  = "task_del"
	taskListCommand = "tast_list"
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
		taskAddCommand,
		taskDelCommand,
		taskListCommand,
	}
}

// UpdateHandler function call for each update
func UpdateHandler(update tgbotapi.Update) (err error) {
	return nil
}

// RunCommand handler start if bot get one of commands
func RunCommand(command string, update tgbotapi.Update) (err error) {
	switch command {
	case taskAddCommand:
		addTask(update.Message)
	case taskDelCommand:
		delTask(update.Message)
	case taskListCommand:
		listTask(update.Message)
	}
	return
}

// StartCommand handler start if bot get one command 'start'
func StartCommand(update tgbotapi.Update) (err error) {
	pluginName := GetName()
	msg := fmt.Sprintf(`Тебя приветствует плагин "Напоминатель"
Для добавления задачи отправь команду /%s_%s.
Для удаления уже созданной задачи отправь команду /%s_%s.
Для вывода списка задач отправь команду /%s_%s.
Приятного пользования "Напоминателем"!`, pluginName, taskAddCommand, pluginName, taskDelCommand, pluginName, taskListCommand)
	ctx.SendMessageText(update.Message.Chat.ID, msg, 0)
	return
}

func addTask(msg *tgbotapi.Message) {
	log.Infof("send command \"%s\"", taskAddCommand)
}

func delTask(msg *tgbotapi.Message) {
	log.Infof("send command \"%s\"", taskDelCommand)
}

func listTask(msg *tgbotapi.Message) {
	log.Infof("send command \"%s\"", taskListCommand)
}
