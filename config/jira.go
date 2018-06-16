package config

type jiraServerConfiguration struct {
	Port      int
	Host      string
	Protocol  string
	Path      string
	Username  string
	Password  string
	Project   string
	IssueType string
}
