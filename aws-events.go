package main

import (
	"os"

	ledisdbConfig "github.com/siddontang/ledisdb/config"

	"github.com/siddontang/ledisdb/ledis"
	log "github.com/sirupsen/logrus"
)

var application = &app{}

func init() {

	application.loadConfig()

	cfg := ledisdbConfig.NewConfigDefault()

	cfg.DataDir = application.Config.Ledis.Path
	cfg.Databases = 1

	l, _ := ledis.Open(cfg)
	db, _ := l.Select(0)

	application.DB = db

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

}

func main() {
	log.Info("Starting event alerter")

	// Get region list
	regionList, err := fetchRegionList()
	if err != err {
		log.Fatal("Unable to fetch region list")
	}

	//Check  each region
	for _, awsRegionName := range regionList {
		resp, err := getRegionInstanceStatus(awsRegionName)
		if err != nil {
			log.Fatal(err)
		}

		log.WithFields(log.Fields{
			"awsProfile": application.Config.AWS.Profile,
			"awsRegion":  awsRegionName,
		}).Info("Processing region")

		// Process instance to see if we need to open event
		processInstance(resp, awsRegionName)
	}

}
