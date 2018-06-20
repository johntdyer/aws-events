package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/coduno/netrc"
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

// Set my defaults
func setMyDefaults(v *viper.Viper) {
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
	v.SetDefault("Ledis.KeyTimeExpireInSeconds", 5184000)

	v.SetDefault("Jira.IssueMapping", map[string]string{
		"production":  "P1",
		"integration": "P2",
		"default":     "P3",
	})
}
func (app *app) loadConfig() {

	v := viper.New()

	setMyDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("couldn't load config: %s", err)
		os.Exit(1)
	}
	var configuration config.Configuration
	if err := v.Unmarshal(&configuration); err != nil {
		fmt.Printf("couldn't read config: %s", err)
	}

	validateAndParseLogLevel(configuration.Application.LogLevel)

	nc, err := netrc.Parse()
	if err != nil {
		log.Error(err)
	}
	app.Config = &configuration

	// Read credentials from netrc if they are not set in config
	if app.Config.Jira.Password == "" {
		app.Config.Jira.Password = nc[app.Config.Jira.Host].Password
		if nc[app.Config.Jira.Host].Login == "" {
			log.Fatal("Jira username required")
		}
	}

	if app.Config.Jira.Username == "" {
		app.Config.Jira.Username = nc[app.Config.Jira.Host].Login
		if nc[app.Config.Jira.Host].Login == "" {
			log.Fatal("Jira password required")
		}
	}

}
