package main

import "multibot/context"

const (
	commandsAddCommand  = "ssh_commands_add"
	commandsDelCommand  = "ssh_commands_del"
	commandsListCommand = "ssh_commands_list"

	commandsAddKeyboard  = "Добавить SSH команду"
	commandsDelKeyboard  = "Удалить SSH команду"
	commandsListKeyboard = "Список SSH команд"
)

var (
	ctx *context.MultiBotContext
)

func InitPlugin(mbc *context.MultiBotContext) (err error) {
	ctx = mbc

	if err = ctx.DBCreateTable(&UserCommand{}); err != nil {
		ctx.Log().Errorf("Unable to create table for ssh user commands: %s", err)
		return
	}

	return
}

func GetName() string {
	return "ssh"
}

func GetDescription() string {
	return "Plugin for run command over SSH on remote hosts periodicaly and notification about it"
}

func GetCommands() []string {
	return []string{
		commandsAddCommand,
		commandsDelCommand,
		commandsListCommand,
	}
}
