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
		AppName:   configName,
		APIKey:    viper.GetString("main.api_key"),
		PgSQLDSN:  viper.GetString("pgsql.dsn"),
		LogLevel:  viper.GetString("log.level"),
		Debug:     viper.GetBool("main.debug"),
		PluginDir: viper.GetString("main.plugin_dir"),
	}
	return
}
