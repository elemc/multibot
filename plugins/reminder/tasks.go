package main

import (
	"fmt"
	"multibot/context"

	"github.com/go-pg/pg"
)

// UserTask struct for store user tasks
type UserTask struct {
	ChatID  int64  `sql:",pk"`
	Name    string `sql:",pk"`
	Type    string
	Year    int
	Month   int
	Day     int
	Hour    int
	Minute  int
	Second  int
	WeekDay int
	Text    string
}

// Save function store user task to database
func (ut *UserTask) Save(ctx *context.MultiBotContext) (err error) {
	db := ctx.GetDB()
	temp := &UserTask{
		ChatID: ut.ChatID,
		Name:   ut.Name,
	}

	if err = db.Select(temp); err != nil && err != pg.ErrNoRows {
		return
	} else if err == pg.ErrNoRows || temp == nil {
		return db.Insert(ut)
	}
	return db.Update(ut)
}

// Delete function remove user task from database
func (ut *UserTask) Delete(ctx *context.MultiBotContext) (err error) {
	return ctx.GetDB().Delete(ut)
}

// String function return task as string
func (ut *UserTask) String() (result string) {
	result = fmt.Sprintf("*%s* (%s)", ut.Name, ut.Type)

	switch ut.Type {
	case taskTypeOnce:
		result = fmt.Sprintf("%s %04d-%02d-%02d %02d:%02d:%02d", result, ut.Year, ut.Month, ut.Day, ut.Hour, ut.Minute, ut.Second)
	case taskTypeYearly:
		result = fmt.Sprintf("%s %s %02d %02d:%02d:%02d", result, months[ut.Month-1], ut.Day, ut.Hour, ut.Minute, ut.Second)
	case taskTypeMonthly:
		result = fmt.Sprintf("%s %02d %02d:%02d:%02d", result, ut.Day, ut.Hour, ut.Minute, ut.Second)
	case taskTypeWeekly:
		result = fmt.Sprintf("%s %s %02d:%02d:%02d", result, weekDays[ut.WeekDay-1], ut.Hour, ut.Minute, ut.Second)
	case taskTypeDaily:
		result = fmt.Sprintf("%s %02d:%02d:%02d", result, ut.Hour, ut.Minute, ut.Second)
	case taskTypeHourly:
		result = fmt.Sprintf("%s %02d:%02d", result, ut.Minute, ut.Second)
	}

	return
}

func getUserTasks(ctx *context.MultiBotContext, chatID int64) (tasks []UserTask, err error) {
	err = ctx.GetDB().Model(&tasks).Where("chat_id = ?", chatID).Select()
	return
}

func getTaskByString(ctx *context.MultiBotContext, chatID int64, s string) (ut *UserTask, err error) {
	var tasks []UserTask
	if tasks, err = getUserTasks(ctx, chatID); err != nil {
		return
	}

	for _, task := range tasks {
		if task.String() == s {
			ut = &task
			return
		}
	}
	return
}
