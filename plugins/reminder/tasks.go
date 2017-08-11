package main

import (
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

// Save function store user state to database
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
