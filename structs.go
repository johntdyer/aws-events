package main

import (
	appConfig "github.com/johntdyer/aws-events/config"
	"github.com/siddontang/ledisdb/ledis"
)

// App config strcut
type app struct {
	// DB *badger.DB
	DB     *ledis.DB
	Config *appConfig.Configuration
}

type issue struct {
	Description    string
	Summary        string
	Tags           []string
	InstanceID     string
	awsRegion      string
	awsEnvironment string
}
