package main

import (
	"flag"
	"os"

	ledisdbConfig "github.com/siddontang/ledisdb/config"

	"github.com/siddontang/ledisdb/ledis"
	log "github.com/sirupsen/logrus"
)

// var myFlags arrayFlags
var (
	application     = &app{}
	regionListFlags regionArrayFromFlags
)

// var application = &app{}

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

	flag.Var(&regionListFlags, "region", "Check a single region.  Optional, default is to check all regions.")
	flag.Parse()

	// Get region list
	// If user passes a list of regions we'll use that otherwise we'll check all
	var regionList []string
	var err error
	if len(regionListFlags) > 0 {
		regionList = regionListFlags
	} else {
		regionList, err = fetchRegionList()
		if err != err {
			log.Fatal("Unable to fetch region list")
		}
	}

	//Check each region
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
