package main

import (
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// type instanceState string

// func (py instanceState) Split(str string) (string, string, error) {
// 	s := strings.Split(string(py), str)
// 	if len(s) < 2 {
// 		return "", "", errors.New("Minimum match not found")
// 	}
// 	return s[0], s[1], nil
// }

// Used to support multiple --region flags and represent as slice
type regionArrayFromFlags []string

func (i *regionArrayFromFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *regionArrayFromFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func newInstanceEvent(key string, awsRegion string) (bool, error) {

	res, err := application.DB.HGetAll([]byte(key)) //, "code")
	if err != nil {
		return false, err
	}

	// Key not not exist there for its a new event and must be handled

	if len(res) == 0 {
		log.WithFields(log.Fields{
			"instanceID": string(key),
			"awsProfile": application.Config.AWS.Profile,
			"awsRegion":  awsRegion,
		}).Info("Found new instance event")
		return true, nil
	}

	return false, nil

}

func trimEventDescription(event string, awsRegion string, instanceID string) string {

	re := regexp.MustCompile("(^The instance\\s)")

	// Do we have a match ?
	matched := re.MatchString(event)

	// Clean up event by removing state prefix
	//   Eg: 'The instance is running on degraded hardware' -> 'is running on degraded hardware'
	afterMatch := re.ReplaceAllString(event, "")

	if matched {
		log.WithFields(log.Fields{
			"instanceID":  instanceID,
			"awsProfile":  application.Config.AWS.Profile,
			"awsRegion":   awsRegion,
			"beforeMatch": event,
			"postMatch":   afterMatch,
		}).Debug("Trimed issue description prefix")
	}

	return afterMatch

}

func maintenanceComplete(event string, awsRegion string, instanceID string) (bool, string) {
	pattern := "(^\\[Completed]\\s)"
	re := regexp.MustCompile(pattern)

	// Do we have a match ?
	matched := re.MatchString(event)

	// Clean up event by removing state prefix
	//   Eg: '[Completed] The instance is running on degraded hardware'  --> 'The instance is running on degraded hardware'
	afterMatch := re.ReplaceAllString(event, "")

	if matched {
		log.WithFields(log.Fields{
			"instanceID":  instanceID,
			"awsProfile":  application.Config.AWS.Profile,
			"awsRegion":   awsRegion,
			"beforeMatch": event,
			"postMatch":   afterMatch,
		}).Debug("Remove [Completed] from description")
	}

	return matched, afterMatch
}
