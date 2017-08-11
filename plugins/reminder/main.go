package main

import (
	"fmt"
	"multibot/context"
	"strconv"
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
	if err := ctx.DBCreateTable(&UserTask{}); err != nil {
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
	default:
		getReminderValues(update.Message)
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
	sendWelcome(update.Message.Chat.ID, msg)
	return
}

func sendWelcome(chatID int64, msg string) {

	row := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(taskAddKeyboard),
		tgbotapi.NewKeyboardButton(taskDelKeyboard),
		tgbotapi.NewKeyboardButton(taskListKeyboard),
	)
	rm := tgbotapi.NewReplyKeyboard(row)
	ctx.SendMessageText(chatID, msg, 0, rm)
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
		if err = rus.Save(ctx); err != nil {
			ctx.Log().WithField("plugin", GetName()).Errorf("Unable to save user state: %s", err)
			return
		}
	}
	ctx.Log().Debugf("State %d (%d) for chat ID %d and command %s", rus.State, taskAddStateSelectType, rus.ChatID, rus.Command)
	switch rus.State {
	case taskAddStateSelectType:
		sendSelectType(msg)
	case taskAddStateSelectYear:
		sendSelectNum(msg, "Введите или выберите год:", time.Now().Year(), time.Now().Year()+3, 1)
	case taskAddStateSelectMonth:
		sendSelectStringSlice(msg, "Выберите месяц или введите номер месяца:", months)
	case taskAddStateSelectDay:
		sendSelectNum(msg, "Введите или выберите день:", 1, rus.GetLastDay(ctx), 1)
	case taskAddStateSelectHour:
		sendSelectNum(msg, "Введите или выберите час:", 0, 23, 1)
	case taskAddStateSelectMinute:
		sendSelectNum(msg, "Введите или выберите минуту:", 0, 59, 5)
	case taskAddStateSelectSecond:
		sendSelectNum(msg, "Введите или выберите секунду:", 0, 59, 30)
	case taskAddStateSelectWeekDay:
		sendSelectStringSlice(msg, "Выберите месяц или введите номер месяца:", weekDays)
	case taskAddStateSelectName:
		ctx.SendMessageText(msg.Chat.ID, "Введите название задачи:", 0, nil)
	case taskAddStateSelectFinish:
		ctx.SendMessageText(msg.Chat.ID, "Введите текст задачи:", 0, nil)
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
	var (
		buttons []tgbotapi.KeyboardButton
		rows    [][]tgbotapi.KeyboardButton
		counter int
	)

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

func sendSelectNum(msg *tgbotapi.Message, query string, begin, end, step int) {
	var (
		buttons []tgbotapi.KeyboardButton
		rows    [][]tgbotapi.KeyboardButton
		counter int
	)

	for i := begin; i <= end; i += step {
		buttons = append(buttons, tgbotapi.NewKeyboardButton(fmt.Sprintf("%d", i)))
		counter++
		if counter == 4 {
			rows = append(rows, buttons)
			buttons = []tgbotapi.KeyboardButton{}
			counter = 0
		}
	}
	if len(buttons) != 0 {
		rows = append(rows, buttons)
	}

	rm := tgbotapi.NewReplyKeyboard(rows...)
	ctx.SendMessageText(msg.Chat.ID, query, 0, rm)

}

func sendSelectStringSlice(msg *tgbotapi.Message, query string, ss []string) {
	var (
		buttons []tgbotapi.KeyboardButton
		rows    [][]tgbotapi.KeyboardButton
		counter int
	)

	for _, value := range ss {
		buttons = append(buttons, tgbotapi.NewKeyboardButton(value))
		counter++
		if counter == 3 {
			rows = append(rows, buttons)
			buttons = []tgbotapi.KeyboardButton{}
			counter = 0
		}
	}
	if len(buttons) != 0 {
		rows = append(rows, buttons)
	}

	rm := tgbotapi.NewReplyKeyboard(rows...)
	ctx.SendMessageText(msg.Chat.ID, query, 0, rm)
}

func getReminderValues(msg *tgbotapi.Message) {
	var (
		//rusAdd, rusDel, rusList             *ReminderUserState
		rusAdd *ReminderUserState
		err    error
		//changedAdd, changedDel, changedList bool
		changedAdd bool
	)

	if rusAdd, _, _, err = getReminderUserStates(ctx, msg.Chat.ID); err != nil {
		ctx.Log().WithField("plugin", GetName()).Errorf("Unable to get reminder states: %s", err)
		return
	}

	if rusAdd != nil {
		switch msg.Text {
		case taskTypeOnce:
			rusAdd.State = taskAddStateSelectYear
			rusAdd.Type = taskTypeOnce
			changedAdd = true
		case taskTypeYearly:
			rusAdd.State = taskAddStateSelectMonth
			rusAdd.Type = taskTypeYearly
			changedAdd = true
		case taskTypeMonthly:
			rusAdd.State = taskAddStateSelectDay
			rusAdd.Type = taskTypeMonthly
			changedAdd = true
		case taskTypeWeekly:
			rusAdd.State = taskAddStateSelectWeekDay
			rusAdd.Type = taskTypeWeekly
			changedAdd = true
		case taskTypeDaily:
			rusAdd.State = taskAddStateSelectHour
			rusAdd.Type = taskTypeDaily
			changedAdd = true
		case taskTypeHourly:
			rusAdd.State = taskAddStateSelectMinute
			rusAdd.Type = taskTypeHourly
			changedAdd = true
		default:
			var num int
			if num, err = strconv.Atoi(msg.Text); err != nil {
				getReminderTextNotNum(msg.Text, rusAdd, msg)
				return
			}
			getReminderTextNum(num, rusAdd, msg)
			return
		}
	}

	if changedAdd {
		if err = rusAdd.Save(ctx); err != nil {
			ctx.Log().WithField("plugin", GetName()).Errorf("Unable to save user state: %s", err)
			return
		}
		addTask(msg)
		return
	}
}

func getReminderTextNotNum(text string, rus *ReminderUserState, msg *tgbotapi.Message) {
	var valid bool
	switch rus.State {
	case taskAddStateSelectWeekDay:
		for i, d := range weekDays {
			if d == text {
				rus.WeekDay = i + 1
				rus.State = taskAddStateSelectHour
				valid = true
			}
		}
	case taskAddStateSelectMonth:
		for i, d := range months {
			if d == text {
				rus.Month = i + 1
				rus.State = taskAddStateSelectDay
				valid = true
			}
		}
	case taskAddStateSelectName:
		rus.Name = text
		rus.State = taskAddStateSelectFinish
		valid = true
	case taskAddStateSelectFinish:
		rus.Text = text
		task := rus.ToUserTask()
		if err := rus.Del(ctx); err != nil {
			ctx.Log().Errorf("Unable to delete user state: %s", err)
			return
		}
		if err := task.Save(ctx); err != nil {
			ctx.Log().Errorf("Unable to save user task: %s", err)
			return
		}
		sendWelcome(msg.Chat.ID, "Задача успешно создана.")
		return
	}
	if !valid {
		ctx.SendMessageMarkdown(msg.Chat.ID, "Введено неправильное значение", 0, nil)
		addTask(msg)
		return
	}

	if err := rus.Save(ctx); err != nil {
		ctx.Log().Errorf("Unable to save user status: %s", err)
		return
	}
	addTask(msg)
}

func getReminderTextNum(num int, rus *ReminderUserState, msg *tgbotapi.Message) {
	valid := true

	switch rus.State {
	case taskAddStateSelectYear:
		rus.Year = num
		rus.State = taskAddStateSelectMonth
	case taskAddStateSelectMonth:
		if num < 0 || num > 12 {
			valid = false
		} else {
			rus.Month = num
			rus.State = taskAddStateSelectDay
		}
	case taskAddStateSelectDay:
		lastDay := rus.GetLastDay(ctx)
		if num < 0 || num > lastDay {
			valid = false
		} else {
			rus.Day = num
			rus.State = taskAddStateSelectHour
		}
	case taskAddStateSelectHour:
		if num < 0 || num > 23 {
			valid = false
		} else {
			rus.Hour = num
			rus.State = taskAddStateSelectMinute
		}
	case taskAddStateSelectMinute:
		if num < 0 || num > 59 {
			valid = false
		} else {
			rus.Minute = num
			rus.State = taskAddStateSelectSecond
		}
	case taskAddStateSelectSecond:
		if num < 0 || num > 59 {
			valid = false
		} else {
			rus.Second = num
			rus.State = taskAddStateSelectName
		}
	}
	if !valid {
		ctx.SendMessageMarkdown(msg.Chat.ID, "Введено неправильное значение", 0, nil)
		addTask(msg)
		return
	}

	if err := rus.Save(ctx); err != nil {
		ctx.Log().Errorf("Unable to save user status: %s", err)
		return
	}
	addTask(msg)
}
