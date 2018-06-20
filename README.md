# AWS Event Tool


## Purpose

This tool is designed to be run in a cron and detected instance events.  When events are detected a JIRA issue is opened up

Example:

![Ticket Example](contrib/example.png?v=4&s=200)

## Use

```bash
cp config-example.toml config.toml
# edit config file
./aws-events
```

The minimal IAM permissions required to run the app are below.  


```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "ec2:DescribeInstances",
                "ec2:DescribeRegions",
                "ec2:DescribeInstanceStatus"
            ],
            "Resource": "*"
        }
    ]
}
```

It is strongly recommended you create a user with only these permissions rather then using your personal keys


#### Flags

By default the application will check all regions, however you can pass one or more `--region` flags to define specific regions to check.

```
./aws-events --region us-east-1 --region us-east-2
```

##### Supported Config Options

| Config | Purpose | Environment Variable | Default |
| --------- |--------- |--------- |--------- |
| application.log_level | Set application log level, supported options are debug, warn, error, fatal | AWS_EVENT_LOG_LEVEL |  info |
| jira.protocol  | Protocol to use, http or https | AWS_EVENT_JIRA_PROTOCOL | https | 
| jira.port  | Jira server port | AWS_EVENT_JIRA_PORT | 443 | 
| jira.host  | Jira server hostname | AWS_EVENT_JIRA_HOST | "jira-eng-gpk2.example.com" | 
| jira.path  | Jira server bath path | AWS_EVENT_JIRA_PATH | /jira |
| jira.environmentPriorityMapping | Mapping between environment and issue priority | | production = "P1" <br/> integration = "P2"<br/> default    = "P3" <br/> | 
| aws.profileName | aws config profilee name |  AWS_EVENT_AWS_PROFILE_NAME | sparkdev |
| ledis.path | Path to database file for state |  AWS_EVENT_LEDIS_PATH | "./database/ledis | 
| ledis.database | Database to use, suggest never changing this | AWS_EVENT_LEDIS_DATABASE | 0 |
| ledis.key_expire_time | Time ( in seconds ) to expire keys in k/v data store, default is 60 days | AWS_EVEMNT_LEDIS_KEY_EXPIRE_TIME | 5184000 |

## Building

```bash
dep ensure
make
```


## Features:

* Mapping of enviromnment to JIRA issue priority.  Example integration is P4 and production is P3.
* Configurable JIRA Server, Issue type, Issue Project
* Toml configuration, with overrides via environmnent variables
* Issue state is tracted in Ledis key value store to avoid duplicate notifications
* Issues are opened with formatting that includes all tags, environment, ect