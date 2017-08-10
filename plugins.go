package main

import (
	"fmt"
	"io/ioutil"
	"multibot/context"
	"os"
	"path/filepath"
	"plugin"

	"gopkg.in/telegram-bot-api.v4"

	log "github.com/sirupsen/logrus"
)

// BotPlugin struct for store one plugin
type BotPlugin struct {
	Name              string
	Description       string
	Commands          []string
	EachUpdateHandler func(tgbotapi.Update) error
	RunCommandHandler func(string, tgbotapi.Update) error
}

var (
	botPlugins          map[string]*BotPlugin
	botPluginsByCommand map[string]*BotPlugin
)

// LoadPlugins load all plugins from directory
func LoadPlugins() (err error) {
	var pluginFiles []os.FileInfo
	botPlugins = make(map[string]*BotPlugin)
	botPluginsByCommand = make(map[string]*BotPlugin)

	if pluginFiles, err = ioutil.ReadDir(options.PluginDir); err != nil {
		return
	}

	for _, pluginFile := range pluginFiles {
		if pluginFile.IsDir() {
			continue
		}
		log.Debugf("Try to load plugin: %s", pluginFile.Name())
		fullPath := filepath.Join(options.PluginDir, pluginFile.Name())

		var (
			p                 *plugin.Plugin
			initPlugin        plugin.Symbol
			getName           plugin.Symbol
			getDescription    plugin.Symbol
			getCommands       plugin.Symbol
			updateHandler     plugin.Symbol
			runCommandHandler plugin.Symbol
		)
		if p, err = plugin.Open(fullPath); err != nil {
			return
		}
		if initPlugin, err = p.Lookup("InitPlugin"); err != nil {
			return
		}
		if err = initPlugin.(func(*context.MultiBotContext) error)(botContext); err != nil {
			return
		}

		if getName, err = p.Lookup("GetName"); err != nil {
			return
		}
		if getDescription, err = p.Lookup("GetDescription"); err != nil {
			return
		}
		if getCommands, err = p.Lookup("GetCommands"); err != nil {
			return
		}
		if updateHandler, err = p.Lookup("UpdateHandler"); err != nil {
			return
		}
		if runCommandHandler, err = p.Lookup("RunCommand"); err != nil {
			return
		}

		botPlugin := &BotPlugin{
			Name:              getName.(func() string)(),
			Description:       getDescription.(func() string)(),
			Commands:          getCommands.(func() []string)(),
			EachUpdateHandler: updateHandler.(func(tgbotapi.Update) error),
			RunCommandHandler: runCommandHandler.(func(string, tgbotapi.Update) error),
		}
		if _, ok := botPlugins[botPlugin.Name]; ok {
			return fmt.Errorf("plugin %s already exists", botPlugin.Name)
		}
		botPlugins[botPlugin.Name] = botPlugin

		for _, cmd := range botPlugin.Commands {
			if pl, ok := botPluginsByCommand[cmd]; ok {
				return fmt.Errorf("we have command with name %s in plugin %s", cmd, pl.Name)
			}
			botPluginsByCommand[cmd] = botPlugin
			log.Debugf("Set command %s for plugin %s", cmd, botPlugin.Name)
		}

		log.Debugf("Loaded plugin: %s (%s)", botPlugin.Name, botPlugin.Description)
	}

	return
}
