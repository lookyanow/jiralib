package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
)

func main() {
	base := "https://my.jira.com"
	tp := jira.BasicAuthTransport{
		Username: "username",
		Password: "token",
	}

	jiraClient, err := jira.NewClient(tp.Client(), base)
	if err != nil {
		panic(err)
	}

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Assignee: &jira.User{
				Name: "myuser",
			},
			Reporter: &jira.User{
				Name: "youruser",
			},
			Description: "Test Issue",
			Type: jira.IssueType{
				Name: "Bug",
			},
			Project: jira.Project{
				Key: "PROJ1",
			},
			Summary: "Just a demo issue",
		},
	}
	issue, _, err := jiraClient.Issue.Create(&i)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s: %+v\n", issue.Key, issue.Fields.Summary)
}
