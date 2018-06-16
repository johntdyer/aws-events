package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
)

func getInstanceTagValue(inst *ec2.Instance, name string, missingValue string) string {
	for _, tag := range inst.Tags {
		if (name != "") && (*tag.Key == name) {
			return *tag.Value
		}
	}
	return missingValue
}

// getInstanceStatus Get list of instances w/ event data
func getRegionInstanceStatus(awsRegion string) (*ec2.DescribeInstanceStatusOutput, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(awsRegion), Credentials: credentials.NewSharedCredentials("", application.Config.AWS.Profile)},
		Profile: application.Config.AWS.Profile,
	}))

	svc := ec2.New(sess)

	params := &ec2.DescribeInstanceStatusInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("event.code"),
				Values: []*string{
					aws.String("instance-reboot"),
					aws.String("instance-stop"),
					aws.String("instance-retirement"),
					aws.String("system-reboot"),
					aws.String("system-maintenance"),
				},
			},
		},
	}

	resp, err := svc.DescribeInstanceStatus(params)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// getInstanceData - Get actual instance data , such as tags, vpc, ect
func getInstanceData(awsRegion string, instanceID string) (*ec2.Instance, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(awsRegion), Credentials: credentials.NewSharedCredentials("", application.Config.AWS.Profile)},
		Profile: application.Config.AWS.Profile,
	}))

	svc := ec2.New(sess)

	params := &ec2.DescribeInstancesInput{

		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-id"),
				Values: []*string{aws.String(instanceID)},
			},
		},
	}
	resp, err := svc.DescribeInstances(params)
	if err != nil {
		log.WithFields(log.Fields{
			"instanceID": instanceID,
			"awsProfile": application.Config.AWS.Profile,
			"awsRegion":  awsRegion,
		}).Error(fmt.Printf("there was an error listing instances in %s", err.Error()))
	}

	// log.WithFields(log.Fields{
	// 	"instanceID": instanceID,
	// 	"awsProfile": application.Config.AWS.Profile,
	// 	"awsRegion":  awsRegion,
	// }).Debug("Found instance data")
	return resp.Reservations[0].Instances[0], nil
}

func fetchRegionList() ([]string, error) {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String("us-west-2"), Credentials: credentials.NewSharedCredentials("", application.Config.AWS.Profile)},
		Profile: application.Config.AWS.Profile,
	}))

	svc := ec2.New(awsSession)

	awsRegions, err := svc.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, err
	}
	regions := make([]string, 0, len(awsRegions.Regions))
	for _, region := range awsRegions.Regions {
		regions = append(regions, *region.RegionName)
	}

	return regions, nil
}

func processInstance(resp *ec2.DescribeInstanceStatusOutput, awsRegion string) {
	// Iterate over them
	for _, instance := range resp.InstanceStatuses {
		for _, event := range instance.Events {

			// If instance ID does not exist in K/V store then its assumed to be a new event and we must act on it
			// Find out if mainteance is complete and parse string

			newInstance, err := newInstanceEvent(*instance.InstanceId, awsRegion)
			if err != nil {
				log.Fatal(err)
			}

			if newInstance == true {
				isComplete, msg := maintenanceComplete(*event.Description, awsRegion, *instance.InstanceId)

				if isComplete {
					log.WithFields(log.Fields{
						"instanceID": *instance.InstanceId,
						"awsProfile": application.Config.AWS.Profile,
						"awsRegion":  awsRegion,
					}).Warn("Maintance completed")
					break
				}

				log.WithFields(log.Fields{
					"instanceID": *instance.InstanceId,
					"awsProfile": application.Config.AWS.Profile,
					"awsRegion":  awsRegion,
					"completed":  isComplete,
					"eventCode":  *event.Code,
					"notBefore":  *event.NotBefore,
				}).Info(msg)

				issue := buildInstanceTicket(*event, *instance, awsRegion)
				issueKey := createJiraIssue(issue)

				_, err := application.DB.HSet([]byte(*instance.InstanceId), []byte("code"), []byte(*event.Code))
				if err != nil {
					log.Fatal(err)
				}

				_, err = application.DB.HSet([]byte(*instance.InstanceId), []byte("description"), []byte(*event.Description))
				if err != nil {
					log.Fatal(err)
				}

				_, err = application.DB.HSet([]byte(*instance.InstanceId), []byte("issue"), []byte(issueKey))
				if err != nil {
					log.Fatal(err)
				}
			} else {
				log.WithFields(log.Fields{
					"instanceID": *instance.InstanceId,
					"awsProfile": application.Config.AWS.Profile,
					"awsRegion":  awsRegion,
				}).Warn("Instance event already processed")
			}
		}
	}
}
