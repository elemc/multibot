package main

import (
	"errors"
	"multibot/context"
	"time"

	"github.com/go-pg/pg"
)

const (
	taskAddStateSelectType = iota
	taskAddStateSelectYear
	taskAddStateSelectMonth
	taskAddStateSelectDay
	taskAddStateSelectHour
	taskAddStateSelectMinute
	taskAddStateSelectSecond

	taskTypeOnce    = "Одиночная"
	taskTypeYearly  = "Ежегодная"
	taskTypeMonthly = "Ежемесячная"
	taskTypeWeekly  = "Еженедельная"
	taskTypeDaily   = "Ежедневная"
	taskTypeHourly  = "Ежечасная"
)

var (
	taskTypes            = []string{taskTypeOnce, taskTypeYearly, taskTypeMonthly, taskTypeWeekly, taskTypeDaily, taskTypeHourly}
	errUserStateNotFound = errors.New("user state not found in database")
)

// ReminderUserState struct for store user command queue state
type ReminderUserState struct {
	ChatID    int64  `sql:",pk"`
	Command   string `sql:",pk"`
	Type      string
	Year      int
	Month     int
	Day       int
	Hour      int
	Minute    int
	Second    int
	State     int
	Timestamp time.Time
}

// initReminderUserState function initialize user state object
func initReminderUserState(chatID int64, cmd string) (rus *ReminderUserState) {
	rus = &ReminderUserState{
		ChatID:  chatID,
		Command: cmd,
	}
	return
}

// Save function store user state to database
func (rus *ReminderUserState) Save(ctx *context.MultiBotContext) (err error) {
	db := ctx.GetDB()
	temp := &ReminderUserState{
		ChatID:  rus.ChatID,
		Command: rus.Command,
	}

	rus.Timestamp = time.Now()
	if db.Select(temp); err != nil && err == pg.ErrNoRows {
		return db.Insert(rus)
	} else if err != nil {
		return
	}
	return db.Update(rus)
}

// getReminderUserState function return user state from database
func getReminderUserState(ctx *context.MultiBotContext, chatID int64, cmd string) (rus *ReminderUserState, err error) {
	rus = initReminderUserState(chatID, cmd)
	if err = ctx.GetDB().Select(rus); err != nil && err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return
	}
	if time.Since(rus.Timestamp) > defaultUserStateTimeout {
		rus = initReminderUserState(chatID, cmd)
	}

	return
}

// Del function remove user state from database after timeout or complete command
func (rus *ReminderUserState) Del(ctx *context.MultiBotContext) (err error) {
	err = ctx.GetDB().Delete(rus)
	return
}
