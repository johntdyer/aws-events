## AWS Eventer


### Purpose

This tool is designed to be run in a cron and detected instance events.  When events are detected a JIRA issue is opened up

### Building


#### Features:

* Mapping of enviromnment to JIRA issue priority.  Example integration is P4 and production is P3.
* Configurable JIRA Server, Issue type, Issue Project
* Toml configuration, with overrides via environmnent variables
* Issue state is tracted in Ledis key value store to avoid duplicate notifications
* Issues are opened with formatting that includes all tags, environment, ect