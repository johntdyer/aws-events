package config

type Configuration struct {
	Jira        jiraServerConfiguration
	AWS         AWSConfig
	Ledis       ledisConfig
	Application applicationConfig
}
