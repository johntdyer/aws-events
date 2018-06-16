package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
)

func tagsToJiraIssue(instance *ec2.Instance) string {
	arr := []string{"||Tag || Value ||"}

	for _, tag := range instance.Tags {
		arr = append(arr, "| "+*tag.Key+" | "+*tag.Value+" |")
	}
	return strings.Join(arr, "\n")
}

func buildInstanceTicket(event ec2.InstanceStatusEvent, instance ec2.InstanceStatus, region string) *issue {
	instanceData, err := getInstanceData(region, *instance.InstanceId)
	if err != nil {
		log.WithFields(log.Fields{
			"instanceID": *instance.InstanceId,
			"awsProfile": application.Config.AWS.Profile,
			"awsRegion":  region,
		}).Error(err)
	}

	// Get instance name by tag
	instanceName := getInstanceTagValue(instanceData, "Name", "unknown")

	// Get instance environment from tag
	instanceEnvironment := getInstanceTagValue(instanceData, "Environment Type", "unknown")

	// Trim instance description
	trimedInstanceDescription := trimEventDescription(*event.Description, region, *instance.InstanceId)

	// Build issue summary
	summary := "Instance: " + instanceName + " - (" + *instanceData.InstanceId + ") " + trimedInstanceDescription

	// Build issue description
	desciption := "[Event Panel|https://" + region + ".console.aws.amazon.com/ec2/v2/home?region=" + region + "#Events]\n\n-------\n\n" +
		"|| || ||\n" +
		"| Instance Name | [" + instanceName + "|https://" + region + ".console.aws.amazon.com/ec2/v2/home?region=" + region + "#Instances:tag:Name=" + instanceName + "] |\n" +
		"| Instance | [" + *instance.InstanceId + "|https://" + region + ".console.aws.amazon.com/ec2/v2/home?region=" + region + "#Instances:instanceId=" + *instance.InstanceId + "]\n" +
		"| Region | [" + region + "|https://console.aws.amazon.com/ec2/v2/home?region=" + region + "#/home]\n" +
		"| State | " + *instance.InstanceState.Name + "|\n" +
		"| Instance Status | " + *instance.InstanceStatus.Status + "|\n" +
		"| Availability Zone | [" + *instance.AvailabilityZone + "|https://" + region + ".console.aws.amazon.com/ec2/v2/home?region=" + region + "#Instances:availabilityZone=" + *instance.AvailabilityZone + "]|\n" +
		"| Tags | {panel}\n" +
		tagsToJiraIssue(instanceData) +
		"|{panel}|\n" +
		"| Event Code | " + *event.Code + " | \n" +
		"| Event Description | " + *event.Description + " | \n"
	// "| Event Description | " + *event.NotBefore.Format("2006-01-02 15:04:05") + " | \n"

	issue := &issue{
		Description:    desciption,
		Summary:        summary,
		Tags:           []string{"aws", "aws_maintenance", fmt.Sprintf("aws_" + instanceEnvironment)},
		InstanceID:     *instanceData.InstanceId,
		awsRegion:      region,
		awsEnvironment: instanceEnvironment,
	}

	return issue
}

func getPriority(env string) *jira.Priority {

	switch env {
	case "production":
		return &jira.Priority{
			Name: "P2",
		}
	case "integration":
		return &jira.Priority{
			Name: "P3",
		}
	default:
		return &jira.Priority{
			Name: "P3",
		}
	}
}

func createJiraIssue(anIssue *issue) string {
	base := application.Config.Jira.Protocol + "://" +
		application.Config.Jira.Host +
		application.Config.Jira.Path

	tp := jira.CookieAuthTransport{
		Username: application.Config.Jira.Username,
		Password: application.Config.Jira.Password,
		AuthURL:  fmt.Sprintf("%s/rest/auth/1/session", base),
	}

	jiraClient, err := jira.NewClient(tp.Client(), base)
	if err != nil {
		panic(err)
	}

	i := jira.Issue{
		Fields: &jira.IssueFields{
			// Assignee: &jira.User{
			// 	Name: "johndye",
			// },
			Description: anIssue.Description,

			Type: jira.IssueType{
				Name: application.Config.Jira.IssueType,
			},
			Project: jira.Project{
				Key: application.Config.Jira.Project,
			},
			Priority: getPriority(anIssue.awsEnvironment),

			Labels:  anIssue.Tags,
			Summary: anIssue.Summary,
		},
	}
	issue, response, err := jiraClient.Issue.Create(&i)
	log.WithFields(log.Fields{
		"instanceID": anIssue.InstanceID,
		"awsProfile": application.Config.AWS.Profile,
		"awsRegion":  anIssue.awsRegion,
	}).Debug(fmt.Printf("%+v", &i))

	if err != nil {

		log.WithFields(log.Fields{
			"instanceID": anIssue.InstanceID,
			"awsProfile": application.Config.AWS.Profile,
			"awsRegion":  anIssue.awsRegion,
		}).Error(err)
	}

	bodyBytes, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 201 {
		log.WithFields(log.Fields{
			"instanceID":   anIssue.InstanceID,
			"awsProfile":   application.Config.AWS.Profile,
			"awsRegion":    anIssue.awsRegion,
			"responseCode": response.Status,
		}).Warn("Body: " + string(bodyBytes))

	} else {
		log.WithFields(log.Fields{
			"instanceID":   anIssue.InstanceID,
			"awsProfile":   application.Config.AWS.Profile,
			"awsRegion":    anIssue.awsRegion,
			"JiraKey":      issue.Key,
			"responseCode": response.StatusCode,
		}).Debug("Body: " + string(bodyBytes))
	}

	return issue.Key

}
