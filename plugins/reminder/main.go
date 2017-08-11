package main

import (
	"fmt"
	"multibot/context"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

const (
	taskAddCommand  = "task_add"
	taskDelCommand  = "task_del"
	taskListCommand = "tast_list"

	taskAddKeyboard  = "Добавить задачу"
	taskDelKeyboard  = "Удалить задачу"
	taskListKeyboard = "Список задач"
)

var (
	ctx                     *context.MultiBotContext
	defaultUserStateTimeout time.Duration
	options                 map[string]interface{}
)

// InitPlugin initialize plugin if it needed
func InitPlugin(mbc *context.MultiBotContext) error {
	ctx = mbc
	options = ctx.GetOptions(GetName())
	if st, ok := options["state_timeout"]; ok && st != nil {
		if d, err := time.ParseDuration(st.(string)); err != nil {
			ctx.Log().Errorf("Unable to parse configuration option %s.state_timeout to duration: %s", GetName(), err)
			defaultUserStateTimeout = time.Hour * 24
		} else {
			defaultUserStateTimeout = d
		}
	}
	if err := ctx.DBCreateTable(&ReminderUserState{}); err != nil {
		ctx.Log().Errorf("Unable to create table for reminder user state: %s", err)
		return err
	}
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
	switch update.Message.Text {
	case taskAddKeyboard:
		addTask(update.Message)
	case taskDelKeyboard:
		delTask(update.Message)
	case taskListKeyboard:
		listTask(update.Message)
	}
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

	row := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(taskAddKeyboard),
		tgbotapi.NewKeyboardButton(taskDelKeyboard),
		tgbotapi.NewKeyboardButton(taskListKeyboard),
	)
	rm := tgbotapi.NewReplyKeyboard(row)

	ctx.SendMessageText(update.Message.Chat.ID, msg, 0, rm)
	return
}

func addTask(msg *tgbotapi.Message) {
	var (
		rus *ReminderUserState
		err error
	)
	if rus, err = getReminderUserState(ctx, msg.Chat.ID, taskAddCommand); err != nil {
		ctx.Log().WithField("plugin", GetName()).Errorf("Unable to add task: %s", err)
		return
	}
	if rus == nil {
		rus = initReminderUserState(msg.Chat.ID, taskAddCommand)
	}
	ctx.Log().Debugf("State %d (%d) for chat ID %d and command %s", rus.State, taskAddStateSelectType, rus.ChatID, rus.Command)
	switch rus.State {
	case taskAddStateSelectType:
		sendSelectType(msg)
	}
}

func delTask(msg *tgbotapi.Message) {
	ctx.Log().Infof("send command \"%s\"", taskDelCommand)
}

func listTask(msg *tgbotapi.Message) {
	ctx.Log().Infof("send command \"%s\"", taskListCommand)
}

func sendSelectType(msg *tgbotapi.Message) {
	ctx.Log().Debugf("For chat %d send select task types", msg.Chat.ID)
	var buttons []tgbotapi.KeyboardButton
	var rows [][]tgbotapi.KeyboardButton
	var counter int

	for _, tt := range taskTypes {
		buttons = append(buttons, tgbotapi.NewKeyboardButton(tt))
		counter++
		if counter == 2 {
			rows = append(rows, buttons)
			buttons = []tgbotapi.KeyboardButton{}
			counter = 0
		}
	}

	rm := tgbotapi.NewReplyKeyboard(rows...)
	ctx.SendMessageText(msg.Chat.ID, "Выберите тип задачи", 0, rm)
}
