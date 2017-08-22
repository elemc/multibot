package main

import (
	"sync"
	"time"
)

// TaskID struct for keys in map timers
type TaskID struct {
	ChatID int64
	Name   string
}

// Timers struct for store all task timers
type Timers struct {
	timers map[TaskID]*time.Timer
	mutex  sync.RWMutex
}

var timers Timers

func loadTasks() (err error) {
	var (
		tasks []UserTask
	)
	timers.timers = make(map[TaskID]*time.Timer)

	if tasks, err = getAllTasks(); err != nil {
		return
	}

	for _, task := range tasks {
		setTimerAndRunJob(task)
	}
	return
}

func setTimerAndRunJob(task UserTask) {
	taskID := TaskID{ChatID: task.ChatID, Name: task.Name}
	timer := task.GetTimer()

	timers.mutex.Lock()
	defer timers.mutex.Unlock()
	timers.timers[taskID] = timer
	go job(timer, task)
	ctx.Log().Debugf("Load task %s for %d and run on %s", task.Name, task.ChatID, task.NextRunTime().String())
}

func job(t *time.Timer, usertask UserTask) {
	if ct, ok := <-t.C; ok {
		ctx.Log().Warnf("Job run %s for user %d: in time %s", usertask.Name, usertask.ChatID, ct.String())
		ctx.SendMessageMarkdown(usertask.ChatID, usertask.Text, 0, nil)
		if usertask.Type != taskTypeOnce {
			setTimerAndRunJob(usertask)
		} else {
			if err := usertask.Delete(); err != nil {
				ctx.Log().Errorf("Unable to delete user task '%d:%s'", usertask.ChatID, usertask.Name)
				return
			}
		}
	}
}
