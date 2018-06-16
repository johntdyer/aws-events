package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/johntdyer/aws-events/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Make sure the config is value
func validateAndParseLogLevel(level string) {
	l := strcase.ToCamel(level)
	switch l {
	case "Panic":
		log.SetLevel(log.PanicLevel)
	case "Error":
		log.SetLevel(log.ErrorLevel)
	case "Warn":
		log.SetLevel(log.WarnLevel)
	case "Info":
		log.SetLevel(log.InfoLevel)
	case "Debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.Fatalf("Log level '%s' is not a supported level\n", l)
		// not really possible but needed

	}
}

func (app *app) loadConfig() {

	v := viper.New()

	v.SetConfigName("config")
	v.AddConfigPath(".")

	// Set defaults
	v.SetEnvPrefix("AWS_EVENT") // will be uppercased automatically
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("LogLevel", "info")
	v.SetDefault("Jira.Port", 443)
	v.SetDefault("Jira.Protocol", "https")
	v.SetDefault("Jira.Project", "SPARK")
	v.SetDefault("Jira.IssueType", "Task")

	v.SetDefault("Ledis.Path", "./database/ledis")
	v.SetDefault("Ledis.Database", 0)

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("couldn't load config: %s", err)
		os.Exit(1)
	}
	var configuration config.Configuration
	if err := v.Unmarshal(&configuration); err != nil {
		fmt.Printf("couldn't read config: %s", err)
	}

	validateAndParseLogLevel(configuration.Application.LogLevel)

	app.Config = &configuration
}
