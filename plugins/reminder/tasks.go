package main

import (
	"fmt"
	"time"

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
func (ut *UserTask) Save() (err error) {
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
	if err = db.Update(ut); err != nil {
		return
	}
	setTimerAndRunJob(ut)
	return
}

// Delete function remove user task from database
func (ut *UserTask) Delete() (err error) {
	if err = ctx.GetDB().Delete(ut); err != nil {
		return
	}

	taskID := TaskID{ChatID: ut.ChatID, Name: ut.Name}
	timers.mutex.Lock()
	defer timers.mutex.Unlock()

	if t, ok := timers.timers[taskID]; ok {
		t.Stop()
	} else {
		ctx.Log().Warnf("Not found timer for chat ID %d and name %s", ut.ChatID, ut.Name)
	}

	return
}

// GetTimer function return time from user task
func (ut *UserTask) GetTimer() (t *time.Timer) {
	t = time.NewTimer(time.Until(ut.NextRunTime()))
	return
}

// NextRunTime function return next time for run task
func (ut *UserTask) NextRunTime() (t time.Time) {
	switch ut.Type {
	case taskTypeOnce:
		t = time.Date(ut.Year, time.Month(ut.Month), ut.Day, ut.Hour, ut.Minute, ut.Second, 0, time.Local)
	case taskTypeYearly:
		t = time.Date(time.Now().Year(), time.Month(ut.Month), ut.Day, ut.Hour, ut.Minute, ut.Second, 0, time.Local)
		if t.Before(time.Now()) {
			t = t.AddDate(1, 0, 0)
		}
	case taskTypeMonthly:
		t = time.Date(time.Now().Year(), time.Now().Month(), ut.Day, ut.Hour, ut.Minute, ut.Second, 0, time.Local)
		if t.Before(time.Now()) {
			t = t.AddDate(0, 1, 0)
		}
	case taskTypeWeekly:
		var day int
		wdn := int(time.Now().Weekday()) - 1
		if wdn < 0 {
			wdn = 6
		}
		if wdn == ut.WeekDay {
			day = time.Now().Day()
		} else if wdn < ut.WeekDay {
			day = time.Now().Day() + (ut.WeekDay - wdn)
		} else {
			day = time.Now().Day() + (7 - wdn - ut.WeekDay)
		}
		t = time.Date(time.Now().Year(), time.Now().Month(), day, ut.Hour, ut.Minute, ut.Second, 0, time.Local)
		if t.Before(time.Now()) {
			t = t.AddDate(0, 0, 7)
		}
	case taskTypeDaily:
		t = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), ut.Hour, ut.Minute, ut.Second, 0, time.Local)
		if t.Before(time.Now()) {
			t.AddDate(0, 0, 0)
		}
	case taskTypeHourly:
		t = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), ut.Minute, ut.Second, 0, time.Local)
		if t.Before(time.Now()) {
			t = t.Add(time.Hour * 1)
		}
	}

	return
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

func getUserTasks(chatID int64) (tasks []UserTask, err error) {
	err = ctx.GetDB().Model(&tasks).Where("chat_id = ?", chatID).Select()
	return
}

func getAllTasks() (tasks []UserTask, err error) {
	err = ctx.GetDB().Model(&tasks).Select()
	return
}

func getTaskByString(chatID int64, s string) (ut *UserTask, err error) {
	var tasks []UserTask
	if tasks, err = getUserTasks(chatID); err != nil {
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
