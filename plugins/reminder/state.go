package main

import (
	"errors"
	"fmt"
	"multibot/context"
	"time"

	"github.com/go-pg/pg"
)

const (
	taskAddStateSelectType = iota
	taskAddStateSelectYear
	taskAddStateSelectMonth
	taskAddStateSelectWeekDay
	taskAddStateSelectDay
	taskAddStateSelectHour
	taskAddStateSelectMinute
	taskAddStateSelectSecond
	taskAddStateSelectName
	taskAddStateSelectFinish
)

const (
	taskDelStateSelectTask = iota
	taskDelStateFinish
)

const (
	taskTypeOnce    = "Одиночная"
	taskTypeYearly  = "Ежегодная"
	taskTypeMonthly = "Ежемесячная"
	taskTypeWeekly  = "Еженедельная"
	taskTypeDaily   = "Ежедневная"
	taskTypeHourly  = "Ежечасная"
)

const (
	weekDayMonday    = "Понедельник"
	weekDayTuesday   = "Вторник"
	weekDayWednesday = "Среда"
	weekDayThursday  = "Четверг"
	weekDayFriday    = "Пятница"
	weekDaySaturday  = "Суббота"
	weekDaySunday    = "Воскресенье"
)

const (
	monthJan = "Январь"
	monthFeb = "Февраль"
	monthMar = "Март"
	monthApr = "Апрель"
	monthMay = "Май"
	monthJun = "Июнь"
	monthJul = "Июль"
	monthAug = "Август"
	monthSep = "Сентябрь"
	monthOct = "Октябрь"
	monthNov = "Ноябрь"
	monthDec = "Декабрь"
)

const (
	globalCancel = "Отмена"
)

var (
	taskTypes            = []string{taskTypeOnce, taskTypeYearly, taskTypeMonthly, taskTypeWeekly, taskTypeDaily, taskTypeHourly}
	weekDays             = []string{weekDayMonday, weekDayTuesday, weekDayWednesday, weekDayThursday, weekDayFriday, weekDaySaturday, weekDaySunday}
	months               = []string{monthJan, monthFeb, monthMar, monthApr, monthMay, monthJun, monthJul, monthAug, monthSep, monthOct, monthNov, monthDec}
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
	WeekDay   int
	Name      string
	Text      string
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
func (rus *ReminderUserState) Save() (err error) {
	db := ctx.GetDB()
	temp := &ReminderUserState{
		ChatID:  rus.ChatID,
		Command: rus.Command,
	}

	rus.Timestamp = time.Now()
	if err = db.Select(temp); err != nil && err != pg.ErrNoRows {
		return
	} else if err == pg.ErrNoRows || temp == nil {
		return db.Insert(rus)
	}
	return db.Update(rus)
}

// getReminderUserState function return user state from database
func getReminderUserState(chatID int64, cmd string) (rus *ReminderUserState, err error) {
	rus = initReminderUserState(chatID, cmd)
	if err = ctx.GetDB().Select(rus); err != nil && err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return
	}
	if time.Since(rus.Timestamp) > defaultUserStateTimeout {
		rus = nil
	}

	return
}

func getReminderUserStates(chatID int64) (rusA, rusD *ReminderUserState, err error) {
	if rusA, err = getReminderUserState(chatID, taskAddCommand); err != nil {
		return
	}
	if rusD, err = getReminderUserState(chatID, taskDelCommand); err != nil {
		return
	}
	return
}

// Delete function remove user state from database after timeout or complete command
func (rus *ReminderUserState) Delete() (err error) {
	err = ctx.GetDB().Delete(rus)
	return
}

// ToUserTask convert user state to user task object
func (rus *ReminderUserState) ToUserTask() (ut *UserTask) {
	ut = &UserTask{
		ChatID:  rus.ChatID,
		Name:    rus.Name,
		Type:    rus.Type,
		Year:    rus.Year,
		Month:   rus.Month,
		Day:     rus.Day,
		Hour:    rus.Hour,
		Minute:  rus.Minute,
		Second:  rus.Second,
		WeekDay: rus.WeekDay,
		Text:    rus.Text,
	}
	return
}

// GetLastDay function return last day of a month in year
func (rus *ReminderUserState) GetLastDay() int {
	var (
		bd    time.Time
		year  int
		month int
		err   error
	)
	if year = rus.Year; year == 0 {
		year = time.Now().Year()
	}
	if month = rus.Month; month == 0 {
		month = 1
	}
	if bd, err = time.Parse("2006.01.02", fmt.Sprintf("%d.%02d.%02d", year, month, 1)); err != nil {
		ctx.Log().Errorf("Unable to parse date %d-%02d-%02d: %s", year, month, 1, err)
		return 0
	}
	lastDay := bd.AddDate(0, 1, -1).Day()
	ctx.Log().Debugf("Calculate last month day for %s is %d", bd.String(), lastDay)
	return lastDay
}

func deleteUserStates(ctx *context.MultiBotContext, chatID int64) (err error) {
	if rusA, rusD, err := getReminderUserStates(chatID); err != nil {
		ctx.Log().Errorf("Unable to get user states: %s", err)
	} else {
		if rusA != nil {
			if err = rusA.Delete(); err != nil {
				return err
			}
		}
		if rusD != nil {
			if err = rusD.Delete(); err != nil {
				return err
			}
		}
	}
	return
}
