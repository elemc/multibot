// -*- Go -*-

package main

import (
	"multibot/context"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var options *context.Options

// LoadConfig function loads configuration file and set options
func LoadConfig() (err error) {
	log.Warnf("Load configuration file...")

	viper.SetConfigName(configName)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("/usr/local/etc")

	viper.SetDefault("main.plugin_dir", "plugins")

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	options = &context.Options{
		AppName:         configName,
		APIKey:          viper.GetString("main.api_key"),
		PgSQLDSN:        viper.GetString("pgsql.dsn"),
		LogLevel:        viper.GetString("log.level"),
		Debug:           viper.GetBool("main.debug"),
		PluginDir:       viper.GetString("main.plugin_dir"),
		PluginsSettings: make(map[string]map[string]interface{}),
	}
	return
}

func loadPLuginConfig(pluginName string) {
	log.Debugf("Loading plugin \"%s\" settings from configuration file...", pluginName)
	sub := viper.Sub(pluginName)
	if sub == nil {
		log.Debugf("No settings for plugin \"%s\"", pluginName)
		return
	}

	tempMap := make(map[string]interface{})
	for _, key := range sub.AllKeys() {
		tempMap[key] = sub.Get(key)
	}
	options.PluginsSettings[pluginName] = tempMap
	log.Debugf("Settings for plugin \"%s\" loaded successful", pluginName)
}
