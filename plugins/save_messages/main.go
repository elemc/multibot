package main

import (
	"fmt"
	"io/ioutil"
	"multibot/context"
	"os"

	"gopkg.in/telegram-bot-api.v4"
)

var (
	f   *os.File
	ctx *context.MultiBotContext
)

// InitPlugin initialize plugin if it needed
func InitPlugin(c *context.MultiBotContext) (err error) {
	ctx = c
	return
}

// GetName function returns plugin name
func GetName() string {
	return "save_messages"
}

// GetDescription function returns plugin description
func GetDescription() string {
	return "Simple log all messages plugin for multibot"
}

// GetCommands return plugin commands for bot
func GetCommands() []string {
	return []string{
		"show_file",
	}
}

// UpdateHandler function call for each update
func UpdateHandler(update tgbotapi.Update) (err error) {
	if f, err = os.OpenFile("qwe.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644); err != nil {
		return
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("%s\t-\t%s\n", update.Message.From.String(), update.Message.Text))
	return nil
}

// RunCommand handler start if bot get one of commands
func RunCommand(command string, update tgbotapi.Update) (err error) {
	switch command {
	case "show_file":
		err = runShowFile(update.Message.Chat.ID)
	}
	return
}

func runShowFile(chatID int64) (err error) {
	if f, err = os.OpenFile("qwe.txt", os.O_RDONLY, 0644); err != nil {
		return
	}
	defer f.Close()

	var data []byte
	if data, err = ioutil.ReadAll(f); err != nil {
		return
	}
	ctx.SendMessageText(chatID, string(data), 0)
	return
}

// StartCommand handler start if bot get one command 'start'
func StartCommand(update tgbotapi.Update) (err error) {
	return
}
